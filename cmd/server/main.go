package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jheimbach/nfc-cash-system/pkg/server"
	"github.com/jheimbach/nfc-cash-system/pkg/server/internals/database"
	"github.com/jheimbach/nfc-cash-system/pkg/server/internals/test"
)

func main() {
	address := flag.String("address", "", "Host address for server")
	port := flag.String("port", "50051", "Host port for server")
	certFile := flag.String("grpc-cert", "../../tls/cert.pem", "TLS certificate for grpc server")
	keyFile := flag.String("grpc-key", "../../tls/cert-key.pem", "TLS key for grpc server")
	accessTknKey := flag.String("access-token-key", "7QC/y4Dkke2izCGyArkfH074ETD9Hyf6PxIV/D7L2Nw=", "TLS key for grpc server")
	refreshTknKey := flag.String("refresh-token-key", "tA2ZFqRCgYBEX4Y9/Q4Au9U0qrbW2oBcqJ8uRPavj9g=", "TLS key for grpc server")

	flag.Parse()

	certPath := test.EnvWithDefault("TLS_CERT", *certFile)
	keyPath := test.EnvWithDefault("TLS_KEY", *keyFile)

	// start Database
	if err := database.CheckEnvVars(); err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	log.Println("connecting to database...")
	db, err := database.OpenDatabase(database.DefaultDSN)
	if err != nil {
		log.Fatalf("could not connect to database, %v", err)
	}
	defer db.Close()
	log.Println("connected to database")

	log.Println("start grpc server...")
	// start grpc server
	grpcSrv, err := server.NewGrpcServer(db, certPath, keyPath, *accessTknKey, *refreshTknKey)
	if err != nil {
		log.Fatalf("could not create grpc server: %v", err)
	}

	go func() {
		// stop server gracefully on sigterm
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		s := <-signals
		// closing grpc server
		log.Printf("recieved signal %s, closing grpc server...\n", s)

		grpcSrv.GracefulStop()
	}()

	log.Println("grpc server started")
	if err := grpcSrv.Start(net.JoinHostPort(*address, *port)); err != nil {
		log.Fatalf("could not start grpc server: %v", err)
	}
}
