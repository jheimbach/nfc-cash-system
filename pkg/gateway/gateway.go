package gateway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jheimbach/nfc-cash-system/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func handler(ctx context.Context, grpcEndpoint, certFile string) (*runtime.ServeMux, error) {
	creds, err := credentials.NewClientTLSFromFile(certFile, "nfc-cash-system.local")
	if err != nil {
		return nil, fmt.Errorf("could not create credentials from %q: %v", certFile, err)
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(func(s string) (string, bool) {
		if s == "x-refresh-token" || s == "X-Refresh-Token" {
			return "x-refresh-token", true
		}
		return "", false
	}))
	err = api.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return nil, errCouldNotRegisterService("user", err)
	}

	err = api.RegisterHealthServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return nil, errCouldNotRegisterService("health", err)
	}

	err = api.RegisterGroupsServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return nil, errCouldNotRegisterService("group", err)
	}

	err = api.RegisterAccountServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return nil, errCouldNotRegisterService("account", err)
	}

	err = api.RegisterTransactionsServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return nil, errCouldNotRegisterService("transaction", err)

	}

	return mux, nil
}

func NewGatewayServer(ctx context.Context, restEndpoint, grpcEndpoint, certFile string) (*http.Server, error) {
	mux, err := handler(ctx, grpcEndpoint, certFile)
	if err != nil {
		return nil, err
	}

	srv := &http.Server{
		Addr:    restEndpoint,
		Handler: mux,
	}
	return srv, nil
}

func errCouldNotRegisterService(service string, err error) error {
	return fmt.Errorf("could not register %sservice: %v", service, err)
}
