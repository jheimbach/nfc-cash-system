package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/JHeimbach/nfc-cash-system/internals/database"
	"github.com/JHeimbach/nfc-cash-system/server"
	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/models/mysql"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func main() {
	grpcHost := flag.String("grpc-host", "", "Host address for grpc server")
	grpcPort := flag.String("grpc-port", "50051", "Port for grpc server")
	restHost := flag.String("rest-host", "", "Host address for rest server")
	restPort := flag.String("rest-port", "8080", "Port for rest server")
	dsn := flag.String("dsn", "${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST})/${DB_NAME}?parseTime=true&multiStatements=true", "MySQL data source name")
	flag.Parse()

	populatedDsn := os.ExpandEnv(*dsn)
	db, err := database.OpenDatabase(populatedDsn)
	if err != nil {
		log.Fatalf("could not open database, %v", err)
	}
	defer db.Close()

	go func(db *sql.DB) {
		if err := startGrpcServer(*grpcHost, *grpcPort, db); err != nil {
			log.Fatal(err)
		}
	}(db)

	if err := startRestGatewayServer(*grpcHost, *grpcPort, *restHost, *restPort); err != nil {
		log.Fatal(err)
	}
}

func startGrpcServer(host, port string, database *sql.DB) error {
	lis, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	server.RegisterAccountServer(s, mysql.NewAccountModel(database))
	server.RegisterGroupServer(s, mysql.NewGroupModel(database))
	server.RegisterTransactionServer(s, mysql.NewTransactionModel(database))

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func startRestGatewayServer(grpcHost, grpcPort, restHost, restPort string) error {
	grpcEndpoint := net.JoinHostPort(grpcHost, grpcPort)
	restEndpoint := net.JoinHostPort(restHost, restPort)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := api.RegisterAccountServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)

	if err != nil {
		return err
	}

	return http.ListenAndServe(restEndpoint, mux)
}
