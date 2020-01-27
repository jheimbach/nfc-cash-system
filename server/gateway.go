package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func GatewayHandler(ctx context.Context, grpcEndpoint, certFile string) (*runtime.ServeMux, error) {
	creds, err := credentials.NewClientTLSFromFile(certFile, "nfc-cash-system.local")
	if err != nil {
		return nil, fmt.Errorf("could not create credentials from %q: %v", certFile, err)
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	mux := runtime.NewServeMux()
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
	mux, err := GatewayHandler(ctx, grpcEndpoint, certFile)
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
