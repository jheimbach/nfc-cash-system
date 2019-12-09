package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/JHeimbach/nfc-cash-system/server"
	"github.com/JHeimbach/nfc-cash-system/server/models/mysql"
	"google.golang.org/grpc"
)

func main() {
	host := flag.String("grpc-host", "", "Host for grpc server")
	port := flag.String("grpc-port", "50051", "Port for grpc server")
	flag.Parse()

	if err := startGrpcServer(*host, *port); err != nil {
		log.Fatal(err)
	}
}

func startGrpcServer(host, port string) error {
	lis, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	server.RegisterAccountServer(s, mysql.NewAccountModel(&sql.DB{}))

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
