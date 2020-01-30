package e2e

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/jheimbach/nfc-cash-system/api"
	"github.com/jheimbach/nfc-cash-system/pkg/server/internals/test"
	"github.com/jheimbach/nfc-cash-system/pkg/server/internals/test/mock"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	grpcServerName = "grpc-endpoint"
	grpcPort       = "50051"
	restPort       = "8080"
	testDataDir    = "./testdata"
	certPath       = "/run/tls/cert.pem"
	keyPath        = "/run/tls/cert-key.pem"
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
	log.Println("db container started")
	db, err := test.OpenAndMigrateDatabase(laddr, "../../../migrations")
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not migrate db container: %w", err)
	}
	log.Println("db migrated")

	files, certTd, err := mock.CertFiles()
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not create cert files %w", err)
	}

	log.Println("starting test container")
	_, err = createAndStartGrpcServer(networkName, files)
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not create grpc container: %w", err)
	}
	log.Println("grpc container started")

	restPort, err := createAndStartRestServer(networkName, files)
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not create rest container: %w", err)
	}
	log.Println("rest container started")

	return db, net.JoinHostPort("", restPort), certTd, nil
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
	if !_dbInitialized {
		_dbInitialized = true
		err = test.SetupDB(_conn, dataFor("end-to-end"))
	}
	return err, func() error {
		_dbInitialized = false
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

func createAndStartGrpcServer(networkName string, certFiles []string) (string, error) {
	dir, _ := os.Getwd()
	ctx := context.Background()
	port, err := nat.NewPort("tcp", grpcPort)
	if err != nil {
		return "", fmt.Errorf("could not create port: %w", err)
	}
	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			FromDockerfile: tc.FromDockerfile{
				Context:    filepath.Join(dir, "../../../"),
				Dockerfile: "server.Dockerfile",
			},
			Networks:       []string{networkName},
			NetworkAliases: map[string][]string{networkName: {grpcServerName}},
			Env: map[string]string{
				"SERVER_DATABASE.USER":     test.MysqlUser,
				"SERVER_DATABASE.PASSWORD": test.MysqlPassword,
				"SERVER_DATABASE.HOST":     "mysql-server",
				"SERVER_DATABASE.NAME":     test.MysqlDatabase,
				"SERVER_PORT":              port.Port(),
				"SERVER_TLS_CERT":          certPath,
				"SERVER_TLS_KEY":           keyPath,
			},
			ExposedPorts: []string{port.Port()},
			WaitingFor:   wait.ForLog("starting grpc server"),
			BindMounts: map[string]string{
				certFiles[0]: certPath,
				certFiles[1]: keyPath,
			},
		},
		Started: false,
	})

	if err != nil {
		return "", fmt.Errorf("failed to create container %w", err)
	}
	err = container.Start(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to start container %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, port)
	if err != nil {
		return "", fmt.Errorf("could not get endpoint address %w", err)
	}

	return mappedPort.Port(), nil
}

func createAndStartRestServer(networkName string, certFiles []string) (string, error) {
	dir, _ := os.Getwd()
	ctx := context.Background()

	port, err := nat.NewPort("tcp", restPort)
	if err != nil {
		return "", fmt.Errorf("could not create port: %w", err)
	}

	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			FromDockerfile: tc.FromDockerfile{
				Context:    filepath.Join(dir, "../../../"),
				Dockerfile: "gateway.Dockerfile",
			},
			Networks: []string{networkName},
			Env: map[string]string{
				"GATEWAY_REST_HOST": "",
				"GATEWAY_REST_PORT": port.Port(),
				"GATEWAY_TLS_CERT":  certPath,
				"GATEWAY_GRPC_HOST": grpcServerName,
				"GATEWAY_GRPC_PORT": grpcPort,
			},
			ExposedPorts: []string{port.Port()},
			WaitingFor:   wait.ForLog("rest server started"),
			BindMounts: map[string]string{
				certFiles[0]: certPath,
			},
		},
		Started: true,
	})

	if err != nil {
		return "", fmt.Errorf("failed to create container %w", err)
	}

	restEndpoint, err := container.MappedPort(ctx, port)
	if err != nil {
		return "", fmt.Errorf("could not get endpoint address %w", err)
	}

	return restEndpoint.Port(), nil
}
