package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var TestCertFile = []byte(`-----BEGIN CERTIFICATE-----
MIICEzCCAXygAwIBAgIQMIMChMLGrR+QvmQvpwAU6zANBgkqhkiG9w0BAQsFADAS
MRAwDgYDVQQKEwdBY21lIENvMCAXDTcwMDEwMTAwMDAwMFoYDzIwODQwMTI5MTYw
MDAwWjASMRAwDgYDVQQKEwdBY21lIENvMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCB
iQKBgQDuLnQAI3mDgey3VBzWnB2L39JUU4txjeVE6myuDqkM/uGlfjb9SjY1bIw4
iA5sBBZzHi3z0h1YV8QPuxEbi4nW91IJm2gsvvZhIrCHS3l6afab4pZBl2+XsDul
rKBxKKtD1rGxlG4LjncdabFn9gvLZad2bSysqz/qTAUStTvqJQIDAQABo2gwZjAO
BgNVHQ8BAf8EBAMCAqQwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDwYDVR0TAQH/BAUw
AwEB/zAuBgNVHREEJzAlggtleGFtcGxlLmNvbYcEfwAAAYcQAAAAAAAAAAAAAAAA
AAAAATANBgkqhkiG9w0BAQsFAAOBgQCEcetwO59EWk7WiJsG4x8SY+UIAA+flUI9
tyC4lNhbcF2Idq9greZwbYCqTTTr2XiRNSMLCOjKyI7ukPoPjo16ocHj+P3vZGfs
h1fIw3cSS2OolhloGw/XM6RWPWtPAlGykKLciQrBru5NAPvCMsb/I1DAceTiotQM
fblo6RBxUQ==
-----END CERTIFICATE-----`)

func TestErrCouldNotRegisterService(t *testing.T) {
	err := errors.New("something went wrong")
	want := fmt.Errorf("could not register userservice: %v", err)
	got := errCouldNotRegisterService("user", err)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v; wanted %v", got, want)
	}
}

func Test_startRestGatewayServer(t *testing.T) {

	tFile, err := ioutil.TempFile("", "*.pem")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer tFile.Close()

	err = ioutil.WriteFile(tFile.Name(), TestCertFile, 0666)
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}

	mux, err := newGatewayServer(context.Background(), "localhost:90050", tFile.Name())
	if err != nil {
		t.Fatalf("could not create rest gateway servemux: %v", err)
	}
	srv := httptest.NewServer(mux)
	defer srv.Close()

	// because we don't start the grpc server, we expect 503 as the "success" code
	tests := []struct {
		name   string
		path   string
		status int
	}{
		{
			name:   "existing route",
			path:   "/v1/accounts",
			status: http.StatusServiceUnavailable,
		},
		{
			name:   "non existing route",
			path:   "/v1/account",
			status: http.StatusNotFound,
		},
		{
			name:   "method not allowed",
			path:   "/v1/user/logout",
			status: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", srv.URL, tt.path), nil)
			if err != nil {
				t.Errorf("got unexpected err: %v", err)
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("got unexpected err: %v", err)
			}

			// because we don't start the grpc server, we expect 503
			if res.StatusCode != tt.status {
				t.Errorf("got %s, expected %s", res.Status, http.StatusText(tt.status))
			}
		})
	}
}
