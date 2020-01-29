package server

import (
	"database/sql"
	"fmt"
	"net"

	"github.com/jheimbach/nfc-cash-system/pkg/server/auth"
	"github.com/jheimbach/nfc-cash-system/pkg/server/handlers"
	"github.com/jheimbach/nfc-cash-system/pkg/server/repositories/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Grpc struct {
	*grpc.Server
}

func NewGrpcServer(database *sql.DB, cert, certKey string, accessTknKey, refreshTknKey string) (*Grpc, error) {
	creds, err := credentials.NewServerTLSFromFile(cert, certKey)
	if err != nil {
		return nil, err
	}

	tokenGen, err := auth.NewJWTAuthenticator(accessTknKey, refreshTknKey)
	if err != nil {
		return nil, err
	}
	s := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(auth.InitInterceptor(tokenGen)))

	handlers.RegisterHealthServer(s)

	userModel := mysql.NewUserModel(database)
	handlers.RegisterUserServer(s, userModel, tokenGen)

	groupRepository := mysql.NewGroupRepository(database)
	handlers.RegisterGroupServer(s, groupRepository)

	accountRepository := mysql.NewAccountRepository(database, groupRepository)
	transactionRepository := mysql.NewTransactionRepository(database, accountRepository)
	handlers.RegisterAccountServer(s, accountRepository, transactionRepository)
	handlers.RegisterTransactionServer(s, transactionRepository)

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
