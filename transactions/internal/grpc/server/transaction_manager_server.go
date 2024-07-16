package server

import (
	"context"
	"errors"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/config"
	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/grpc/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/transactions/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type (
	TransactionManager interface {
		SubscribeForDispatch(ctx context.Context, email, dispatchID string) error
		UnsubscribeFromDispatch(ctx context.Context, email, dispatchID string) error
	}
	transactionManagerServer struct {
		grpc_gen.UnimplementedTransactionManagerServer
		manager TransactionManager
		logger  config.Logger
	}
)

func NewTransactionManagerServer(manager TransactionManager, logger config.Logger) *transactionManagerServer {
	return &transactionManagerServer{
		manager: manager,
		logger:  logger,
	}
}

func (s *transactionManagerServer) log(method string, req any) {
	s.logger.Infof("DispatchService.%v(%+v)", method, req)
}

func (s *transactionManagerServer) SubscribeForDispatch(ctx context.Context, req *grpc_gen.SubscribeForDispatchRequest) (*emptypb.Empty, error) {
	s.log("SubscribeForDispatch", req.String())

	err := s.manager.SubscribeForDispatch(ctx, req.Email, req.DispatchId)
	if errors.Is(err, service.ErrNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, service.ErrAlreadyExists) {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *transactionManagerServer) UnsubscribeFromDispatch(ctx context.Context, req *grpc_gen.UnsubscribeFromDispatchRequest) (*emptypb.Empty, error) {
	s.log("UnsubscribeFromDispatch", req.String())

	err := s.manager.UnsubscribeFromDispatch(ctx, req.Email, req.DispatchId)
	if errors.Is(err, service.ErrNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
