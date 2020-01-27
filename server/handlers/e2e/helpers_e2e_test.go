package e2e

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/internals/test"
	"github.com/JHeimbach/nfc-cash-system/server/internals/test/mock"
	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	GrpcEndpoint = "localhost:5550"
	testDataDir  = "./testdata"
)

// login returns jwt access and refresh token
func login() (string, string, error) {
	req, err := http.NewRequest(http.MethodGet, RestUrlWithPath("v1/user/login"), nil)
	if err != nil {
		return "", "", fmt.Errorf("could not create auth request: %w", err)
	}
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("briste0@angelfire.com:lMvZARjM3pwe")))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("could not do auth request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var jErr *api.Status
		err := json.NewDecoder(res.Body).Decode(&jErr)
		if err != nil {
			return "", "", fmt.Errorf("could not parse error: %v", err)
		}
		return "", "", fmt.Errorf("could not authenticate: %s", jErr.Message)
	}

	var tkn api.AuthenticateResponse
	err = json.NewDecoder(res.Body).Decode(&tkn)
	if err != nil {
		if err, ok := err.(*json.UnmarshalTypeError); ok {
			if err.Field != "expires_in" {
				return "", "", fmt.Errorf("could not parse auth response: %v", err)
			}
		}
	}

	return tkn.AccessToken, tkn.RefreshToken, nil
}

func dataFor(t string) string {
	return path.Join(testDataDir, strings.Join([]string{t, "sql"}, "."))
}

func startServers() (*sql.DB, string, func(), error) {
	ctx := context.Background()

	const networkName = "integration-test-rest"
	network := tc.NetworkRequest{
		Driver: "bridge",
		Name:   networkName,
	}

	provider, err := tc.NewDockerProvider()
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not create docker provider %w", err)
	}

	if _, err := provider.GetNetwork(context.Background(), network); err != nil {
		if _, err := provider.CreateNetwork(context.Background(), network); err != nil {
			return nil, "", nil, fmt.Errorf("could not create network %q, %w", networkName, err)
		}
	}

	log.Println("starting db container")
	laddr, _, err := test.StartDbContainer(networkName)
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not start db container: %w", err)
	}
	db, err := test.OpenAndMigrateDatabase(laddr)
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not migrate db container: %w", err)
	}

	restPort, err := nat.NewPort("tcp", "8080")
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not create rest port: %w", err)
	}
	grpcPort, err := nat.NewPort("tcp", "50051")
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not create grpc port: %w", err)
	}
	files, certTd, err := mock.CertFiles()
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not create cert files %w", err)
	}

	dir, _ := os.Getwd()
	log.Println("starting test container")
	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			FromDockerfile: tc.FromDockerfile{
				Context:    filepath.Join(dir, "../.."),
			},
			Networks: []string{networkName},
			Labels: map[string]string{
				tc.TestcontainerLabel:  "true",
				"nfc-cash-test-server": "true",
			},
			Env: map[string]string{
				"DB_USER":     test.MysqlUser,
				"DB_PASSWORD": test.MysqlPassword,
				"DB_HOST":     "mysql-server",
				"DB_NAME":     test.MysqlDatabase,
			},
			ExposedPorts: []string{restPort.Port(), grpcPort.Port()},
			WaitingFor:   wait.ForLog("rest server started"),
			BindMounts: map[string]string{
				files[0]: "/run/tls/cert.pem",
				files[1]: "/run/tls/cert-key.pem",
			},
		},
		Started: true,
	})
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not create container %w", err)
	}

	restEndpoint, err := container.MappedPort(ctx, restPort)
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not get container endpoint address %w", err)
	}

	return db, fmt.Sprintf("0.0.0.0:%s", restEndpoint.Port()), certTd, nil
}

func RestUrlWithPath(path string) string {
	return fmt.Sprintf("http://%s/%s", _restEndpoint, path)
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

func setupDb() (error, func() error) {
	var err error
	if !_dbinitialized {
		_dbinitialized = true
		err = test.SetupDB(_conn, dataFor("end-to-end"))
	}
	return err, func() error {
		_dbinitialized = false
		return test.TeardownDB(_conn, dataFor("teardown"))()
	}
}

func prepareTest(t *testing.T) func() error {
	t.Helper()
	test.IsIntegrationTest(t)
	err, td := setupDb()
	if err != nil {
		t.Log(err)
	}
	return td
}
