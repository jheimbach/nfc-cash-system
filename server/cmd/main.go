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

	"github.com/JHeimbach/nfc-cash-system/server"
	"github.com/JHeimbach/nfc-cash-system/server/internals/database"
)

func main() {
	grpcEndpoint := flag.String("grpc-host", "localhost:50051", "Host address for grpc server")
	restEndpoint := flag.String("rest-host", "localhost:8080", "Host address for rest server")
	certFile := flag.String("grpc-cert", "./tls/cert.pem", "TLS certificate for grpc server")
	keyFile := flag.String("grpc-key", "./tls/cert-key.pem", "TLS key for grpc server")

	flag.Parse()

	// start Database
	if err := database.CheckEnvVars(); err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	db, err := database.OpenDatabase(database.DefaultDSN)
	if err != nil {
		log.Fatalf("could not connect to database, %v", err)
	}
	defer db.Close()

	// start grpc server
	grpcSrv, err := server.NewGrpcServer(db, *certFile, *keyFile)
	if err != nil {
		log.Fatalf("could not create grpc server: %v", err)
	}
	go func() {
		if err := grpcSrv.Start(*grpcEndpoint); err != nil {
			log.Fatalf("could not start grpc server: %v", err)
		}
	}()

	// start rest server
	restCtx, restCancel := context.WithCancel(context.Background())
	defer restCancel()

	srv, err := server.NewGatewayServer(restCtx, *restEndpoint, *grpcEndpoint, *certFile)
	if err != nil {
		log.Fatalf("could not create rest gateway server: %v", err)
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

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

	// closing grpc server
	log.Printf("...closing grpc server\n")
	grpcSrv.GracefulStop()
}
