package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jheimbach/nfc-cash-system/pkg/server"
	"github.com/jheimbach/nfc-cash-system/pkg/server/internals/database"
	"github.com/spf13/viper"
)

func main() {
	err := initConfig()
	if err != nil {
		log.Fatalf("could not load config: %v\n", err)
	}
	log.Println("connecting to database...")
	db, err := database.OpenDatabase(
		database.CreateDsn(
			viper.GetString("database.user"),
			viper.GetString("database.password"),
			viper.GetString("database.host"),
			viper.GetString("database.name"),
		),
	)
	if err != nil {
		log.Fatalf("could not connect to database, %v", err)
	}
	defer db.Close()
	log.Println("connected to database")

	log.Println("start grpc server...")
	// start grpc server
	grpcSrv, err := server.NewGrpcServer(
		db,
		viper.GetString("tls_cert"),
		viper.GetString("tls_key"),
		viper.GetString("access_token_key"),
		viper.GetString("refresh_token_key"),
	)
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

	endpointAddr := net.JoinHostPort(viper.GetString("host"), viper.GetString("port"))
	log.Printf("starting grpc server at: %s", endpointAddr)
	if err := grpcSrv.Start(endpointAddr); err != nil {
		log.Fatalf("could not start grpc server at %s: %v", endpointAddr, err)
	}
}
