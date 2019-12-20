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
	"os/signal"
	"syscall"
	"time"

	"github.com/JHeimbach/nfc-cash-system/internals/database"
	"github.com/JHeimbach/nfc-cash-system/server"
	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/auth"
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

	grpcEndpoint := net.JoinHostPort(*grpcHost, *grpcPort)
	restEndpoint := net.JoinHostPort(*restHost, *restPort)

	grpcSrv, err := newGrpcServer(db, *certFile, *keyFile)
	if err != nil {
		log.Fatalf("could not create grpc server: %v", err)
	}
	go func() {
		if err := startGrpcServer(grpcSrv, grpcEndpoint); err != nil {
			log.Fatalf("could not start grpc server: %v", err)
		}
	}()

	restCtx, restCancel := context.WithCancel(context.Background())
	defer restCancel()

	mux, err := newGatewayServer(restCtx, grpcEndpoint, *certFile)
	if err != nil {
		log.Fatalf("could not create rest gateway server: %v", err)
	}

	srv := http.Server{
		Addr:    restEndpoint,
		Handler: mux,
	}

	go func() {
		if err := http.ListenAndServe(restEndpoint, mux); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	s := <-signals
	log.Printf("recieved signal %s, closing rest server\n", s)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("could not shutdown server gracefully: %v", err)
	}
	log.Printf("...closing grpc server\n")
	grpcSrv.GracefulStop()
}

func startGrpcServer(s *grpc.Server, grpcEndpoint string) error {
	lis, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		return fmt.Errorf("could not listen to grpcEndpoint %s: %v", grpcEndpoint, err)
	}
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("could not serve grpc server: %v", err)
	}

	return nil
}

func newGrpcServer(database *sql.DB, cert, certKey string) (*grpc.Server, error) {
	creds, err := credentials.NewServerTLSFromFile(cert, certKey)
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(auth.UnaryInterceptor))

	userModel := mysql.NewUserModel(database)
	server.RegisterUserServer(s, userModel)

	groupModel := mysql.NewGroupModel(database)
	server.RegisterGroupServer(s, groupModel)

	accountModel := mysql.NewAccountModel(database, groupModel)
	server.RegisterAccountServer(s, accountModel)

	transactionModel := mysql.NewTransactionModel(database, accountModel)
	server.RegisterTransactionServer(s, transactionModel)

	return s, nil
}

func newGatewayServer(ctx context.Context, grpcEndpoint, certFile string) (*runtime.ServeMux, error) {
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		return nil, fmt.Errorf("could not create credentials from %q: %v", certFile, err)
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	mux := runtime.NewServeMux()
	err = api.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return nil, errCouldNotRegisterService("user", err)
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

func errCouldNotRegisterService(service string, err error) error {
	return fmt.Errorf("could not register %sservice: %v", service, err)
}
