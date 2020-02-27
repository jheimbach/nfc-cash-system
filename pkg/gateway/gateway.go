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
		Handler: secureHeaders(mux),
	}
	return srv, nil
}

func errCouldNotRegisterService(service string, err error) error {
	return fmt.Errorf("could not register %sservice: %v", service, err)
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1;mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
