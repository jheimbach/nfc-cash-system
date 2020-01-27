package main

import (
	"context"
	"fmt"
	"log"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile("../tls/cert.pem", "")
	if err != nil {
		fmt.Println(err)
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	conn, err := grpc.Dial("localhost:50051", opts...)
	if err != nil {
		log.Fatal(err)
	}

	client := api.NewUserServiceClient(conn)
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{"Authorization": "Basic YnJpc3RlMEBhbmdlbGZpcmUuY29tOmxNdlpBUmpNM3B3ZQ=="}))
	auth, err  := client.AuthenticateUser(ctx, &empty.Empty{})
	if err != nil {
		fmt.Printf("%q", err)
	}
	fmt.Printf("%q", auth)
}


