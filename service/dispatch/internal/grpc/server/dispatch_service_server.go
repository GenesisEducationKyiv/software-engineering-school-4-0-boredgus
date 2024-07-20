package server

import (
	"context"
	"errors"

	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/grpc/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type (
	DispatchService interface {
		SubscribeForDispatch(ctx context.Context, email, dispatch string) error
		UnsubscribeFromDispatch(ctx context.Context, email, dispatch string) error
	}

	dispatchServiceServer struct {
		service DispatchService
		logger  config.Logger
		grpc_gen.UnimplementedDispatchServiceServer
	}
)

func NewDispatchServiceServer(s DispatchService, l config.Logger) *dispatchServiceServer {
	return &dispatchServiceServer{service: s, logger: l}
}

func (s *dispatchServiceServer) log(method string, req any) {
	s.logger.Infof("DispatchService.%v(%+v)", method, req)
}

func (s *dispatchServiceServer) SubscribeForDispatch(ctx context.Context, req *grpc_gen.SubscribeForDispatchRequest) (*emptypb.Empty, error) {
	s.log("SubscribeForDispatch", req.String())
	err := s.service.SubscribeForDispatch(ctx, req.Email, req.DispatchId)
	if errors.Is(err, service.UniqueViolationErr) {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}
	if errors.Is(err, service.NotFoundErr) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *dispatchServiceServer) UnsubscribeFromDispatch(ctx context.Context, req *grpc_gen.UnsubscribeFromDispatchRequest) (*emptypb.Empty, error) {
	s.log("UnsubscribeFromDispatch", req.String())

	err := s.service.UnsubscribeFromDispatch(ctx, req.Email, req.DispatchId)
	if errors.Is(err, service.NotFoundErr) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
