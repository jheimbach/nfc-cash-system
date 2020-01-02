package server

import (
	"database/sql"
	"fmt"
	"net"

	"github.com/JHeimbach/nfc-cash-system/server/auth"
	"github.com/JHeimbach/nfc-cash-system/server/handlers"
	"github.com/JHeimbach/nfc-cash-system/server/models/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Grpc struct {
	*grpc.Server
}

func NewGrpcServer(database *sql.DB, cert, certKey string) (*Grpc, error) {
	creds, err := credentials.NewServerTLSFromFile(cert, certKey)
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(auth.UnaryInterceptor))

	userModel := mysql.NewUserModel(database)
	handlers.RegisterUserServer(s, userModel)

	groupModel := mysql.NewGroupModel(database)
	handlers.RegisterGroupServer(s, groupModel)

	accountModel := mysql.NewAccountModel(database, groupModel)
	handlers.RegisterAccountServer(s, accountModel)

	transactionModel := mysql.NewTransactionModel(database, accountModel)
	handlers.RegisterTransactionServer(s, transactionModel)

	return &Grpc{Server: s}, nil
}

func (s *Grpc) Start(endpoint string) error {
	lis, err := net.Listen("tcp", endpoint)
	if err != nil {
		return fmt.Errorf("could not listen to endpoint %s: %v", endpoint, err)
	}

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("could not serve grpc server: %v", err)
	}

	return nil
}
