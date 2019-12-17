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
	"google.golang.org/grpc/credentials"
)

func main() {
	grpcHost := flag.String("grpc-host", "localhost", "Host address for grpc server")
	grpcPort := flag.String("grpc-port", "50051", "Port for grpc server")
	restHost := flag.String("rest-host", "localhost", "Host address for rest server")
	restPort := flag.String("rest-port", "8080", "Port for rest server")
	dsn := flag.String("dsn", "${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST})/${DB_NAME}?parseTime=true&multiStatements=true", "MySQL data source name")
	certFile := flag.String("grpc-cert", "./tls/cert.pem", "TLS certificate for grpc server")
	keyFile := flag.String("grpc-key", "./tls/cert-key.pem", "TLS key for grpc server")

	flag.Parse()

	populatedDsn := os.ExpandEnv(*dsn)
	db, err := database.OpenDatabase(populatedDsn)
	if err != nil {
		log.Fatalf("could not open database, %v", err)
	}
	defer db.Close()

	go func(db *sql.DB) {
		if err := startGrpcServer(*grpcHost, *grpcPort, db, *certFile, *keyFile); err != nil {
			log.Fatal(err)
		}
	}(db)

	if err := startRestGatewayServer(*grpcHost, *grpcPort, *restHost, *restPort, *certFile); err != nil {
		log.Fatal(err)
	}
}

func startGrpcServer(host, port string, database *sql.DB, cert, certKey string) error {
	lis, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile(cert, certKey)
	if err != nil {
		return err
	}

	s := grpc.NewServer(grpc.Creds(creds))

	userModel := mysql.NewUserModel(database)
	groupModel := mysql.NewGroupModel(database)
	accountModel := mysql.NewAccountModel(database, groupModel)
	transactionModel := mysql.NewTransactionModel(database, accountModel)
	server.RegisterAuthServer(s, userModel)
	server.RegisterGroupServer(s, groupModel)
	server.RegisterAccountServer(s, accountModel)
	server.RegisterTransactionServer(s, transactionModel)

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func startRestGatewayServer(grpcHost, grpcPort, restHost, restPort, certFile string) error {
	grpcEndpoint := net.JoinHostPort(grpcHost, grpcPort)
	restEndpoint := net.JoinHostPort(restHost, restPort)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		return err
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	mux := runtime.NewServeMux()
	err = api.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return err
	}

	err = api.RegisterGroupsServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return err
	}

	err = api.RegisterAccountServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return err
	}

	err = api.RegisterTransactionsServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(restEndpoint, mux)
}
