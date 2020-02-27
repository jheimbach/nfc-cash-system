package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jheimbach/nfc-cash-system/pkg/gateway"
	"github.com/spf13/viper"
)

func main() {
	/*
		restEndpoint := flag.String("rest-host", ":8080", "Host address for rest server")
		grpcEndpoint := flag.String("grpc-host", ":50051", "Host address for grpc server")
		certFile := flag.String("grpc-cert", "./tls/cert.pem", "TLS certificate for grpc server")
		flag.Parse()
	*/

	err := initConfig()
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(os.Stdout, "[REST-GATEWAY] ", log.Flags())

	// start rest server
	restCtx, restCancel := context.WithCancel(context.Background())
	defer restCancel()

	logger.Println("start rest server...")
	srv, err := gateway.NewGatewayServer(
		restCtx,
		net.JoinHostPort(viper.GetString("rest_host"), viper.GetString("rest_port")),
		net.JoinHostPort(viper.GetString("grpc_host"), viper.GetString("grpc_port")),
		viper.GetString("tls_cert"))

	if err != nil {
		logger.Fatalf("could not create rest gateway server: %v", err)
	}
	srv.ErrorLog = logger

	go handleShutdown(srv)

	logger.Println("rest server started")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal(err)
	}
}

func handleShutdown(srv *http.Server) {
	// stop server gracefully on sigterm
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	s := <-signals
	log.Printf("recieved signal %s, shutting rest server down...\n", s)
	// closing rest server
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("could not shutdown server gracefully: %v", err)
	} else {
		log.Println("rest server shutdown, goodbye")
	}
}
