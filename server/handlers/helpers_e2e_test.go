package handlers

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/auth"
	"github.com/JHeimbach/nfc-cash-system/server/internals/test"
	"github.com/JHeimbach/nfc-cash-system/server/internals/test/mock"
	"github.com/JHeimbach/nfc-cash-system/server/repositories/mysql"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	GrpcEndpoint = "localhost:5550"
	RestEndPoint = "localhost:8880"
	dbName       = "nfc-cash-system_test"
	migrationDir = "../migrations"
	testDataDir  = "../testdata"
)

// login returns jwt access and refresh token
func login(t *testing.T) (string, string) {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, RestUrlWithPath("v1/user/login"), nil)
	if err != nil {
		t.Skipf("could not create auth request: %v", err)
	}
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("briste0@angelfire.com:lMvZARjM3pwe")))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Skipf("could not do auth request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var jErr *api.Status
		err := json.NewDecoder(res.Body).Decode(&jErr)
		if err != nil {
			t.Skipf("could not parse error: %v", err)
		}
		t.Skipf("could not authenticate: %s", jErr.Message)
	}

	var tkn api.AuthenticateResponse
	err = json.NewDecoder(res.Body).Decode(&tkn)
	if err != nil {
		if err, ok := err.(*json.UnmarshalTypeError); ok {
			if err.Field != "expires_in" {
				t.Skipf("could not parse auth response: %v", err)
			}
		}
	}

	return tkn.AccessToken, tkn.RefreshToken
}

func dataFor(t string) string {
	return path.Join(testDataDir, strings.Join([]string{t, "sql"}, "."))
}

func getTestDb(t *testing.T) (*sql.DB, func(...string), func()) {
	migrationDir := test.EnvWithDefault("DB_MIGRATIONS_DIR", migrationDir)
	dbName := test.EnvWithDefault("TEST_DB_NAME", dbName)
	return test.GetDb(t, dbName, dataFor("teardown"), migrationDir)
}

func startServers(t *testing.T) (teardown func()) {
	t.Helper()

	db, setup, dbTeardown := getTestDb(t)
	setup(dataFor("end-to-end"))

	certs, rmCertFiles := mock.CertFiles(t)

	grpcSrv, err := newGrpcServer(db, certs[0], certs[1])
	if err != nil {
		dbTeardown()
		log.Fatalf("could not create grpc server: %v", err)
	}
	go func() {
		lis, err := net.Listen("tcp", GrpcEndpoint)
		if err != nil {
			t.Skipf("could not listen to endpoint %s: %v", GrpcEndpoint, err)
		}

		if err := grpcSrv.Serve(lis); err != nil {
			t.Skipf("could not serve grpc server: %v", err)
		}
	}()

	mux, err := newGateWayHandler(certs[0])
	srv := &http.Server{
		Addr:    RestEndPoint,
		Handler: mux,
	}
	if err != nil {
		dbTeardown()
		grpcSrv.GracefulStop()
		log.Fatalf("could not create rest server: %v", err)
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Skip(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	return func() {
		srv.Shutdown(context.Background())
		grpcSrv.GracefulStop()
		dbTeardown()
		rmCertFiles()
	}
}

func newGrpcServer(database *sql.DB, cert, key string) (*grpc.Server, error) {
	creds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		return nil, err
	}
	tokenGen, _ := auth.NewJWTAuthenticator("e5pQGs6UOjZxZLjOg0Z5jNQ1JzARut4NZ3JGf7e0", "R04y8DEVfLoGm91oDRjIgDIsH1n1Y8ak0JFloBIm")

	s := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(auth.InitInterceptor(tokenGen)))

	userModel := mysql.NewUserModel(database)
	RegisterUserServer(s, userModel, tokenGen)

	groupRepository := mysql.NewGroupRepository(database)
	RegisterGroupServer(s, groupRepository)

	accountRepository := mysql.NewAccountRepository(database, groupRepository)
	transactionRepository := mysql.NewTransactionRepository(database, accountRepository)
	RegisterAccountServer(s, accountRepository, transactionRepository)
	RegisterTransactionServer(s, transactionRepository)

	return s, nil
}

func newGateWayHandler(cert string) (*runtime.ServeMux, error) {
	creds, err := credentials.NewClientTLSFromFile(cert, "")
	if err != nil {
		return nil, err
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	ctx := context.Background()
	mux := runtime.NewServeMux()
	api.RegisterUserServiceHandlerFromEndpoint(ctx, mux, GrpcEndpoint, opts)
	api.RegisterGroupsServiceHandlerFromEndpoint(ctx, mux, GrpcEndpoint, opts)
	api.RegisterAccountServiceHandlerFromEndpoint(ctx, mux, GrpcEndpoint, opts)
	api.RegisterTransactionsServiceHandlerFromEndpoint(ctx, mux, GrpcEndpoint, opts)

	return mux, nil
}

func RestUrlWithPath(path string) string {
	return fmt.Sprintf("http://%s/%s", RestEndPoint, path)
}

func checkError(response *http.Response, code int, errMsg string) error {
	if response.StatusCode != code {
		return fmt.Errorf("got statuscode %d, expected %d", response.StatusCode, code)
	}

	status, err := parseErrorMsg(response)
	if err != nil {
		return err
	}

	if status.Message != errMsg {
		return fmt.Errorf("got err msg: %q, wanted: %q", status.Message, errMsg)
	}

	return nil
}

func parseErrorMsg(response *http.Response) (*api.Status, error) {
	var jsonErr *api.Status
	err := json.NewDecoder(response.Body).Decode(&jsonErr)
	if err != nil {
		return nil, fmt.Errorf("could not parse status: %w", err)
	}
	return jsonErr, nil
}

func checkUnwantedErr(response *http.Response) error {
	if response.StatusCode != http.StatusOK {
		status, err := parseErrorMsg(response)
		if err != nil {
			return err
		}
		return fmt.Errorf("unexpected response error: %v", status)
	}
	return nil
}
