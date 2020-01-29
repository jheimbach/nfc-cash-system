package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jheimbach/nfc-cash-system/pkg/gateway"
	"github.com/jheimbach/nfc-cash-system/pkg/server/internals/test"
)

func main() {
	restEndpoint := flag.String("rest-host", ":8080", "Host address for rest server")
	grpcEndpoint := flag.String("grpc-host", ":50051", "Host address for grpc server")
	certFile := flag.String("grpc-cert", "./tls/cert.pem", "TLS certificate for grpc server")
	flag.Parse()

	certPath := test.EnvWithDefault("TLS_CERT", *certFile)
	grpcHost := test.EnvWithDefault("GRPC_HOST", *grpcEndpoint)

	// start rest server
	restCtx, restCancel := context.WithCancel(context.Background())
	defer restCancel()

	log.Println("start rest server...")
	srv, err := gateway.NewGatewayServer(restCtx, *restEndpoint, grpcHost, certPath)
	if err != nil {
		log.Fatalf("could not create rest gateway server: %v", err)
	}
	srv.ErrorLog = log.New(os.Stdout, "REST-GATEWAY", log.Flags())

	go func() {
		// stop server gracefully on sigterm
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		s := <-signals
		log.Printf("recieved signal %s, closing rest server...\n", s)
		// closing rest server
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("could not shutdown server gracefully: %v", err)
		}
	}()
	log.Println("rest server started")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
