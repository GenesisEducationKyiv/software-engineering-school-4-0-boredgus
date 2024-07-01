package server

import (
	"context"
	"errors"

	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/grpc/gen"
	service_errors "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service/err"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/repo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	DispatchService interface {
		SubscribeForDispatch(ctx context.Context, email, dispatch string) error
		SendDispatch(ctx context.Context, dispatch string) error
		GetAllDispatches(ctx context.Context) ([]repo.DispatchData, error)
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

func (s *dispatchServiceServer) SubscribeForDispatch(ctx context.Context, req *grpc_gen.SubscribeForDispatchRequest) (*grpc_gen.SubscribeForDispatchResponse, error) {
	s.log("SubscribeForDispatch", req.String())
	err := s.service.SubscribeForDispatch(ctx, req.Email, req.DispatchId)
	if errors.Is(err, service_errors.UniqueViolationErr) {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}
	if errors.Is(err, service_errors.NotFoundErr) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpc_gen.SubscribeForDispatchResponse{}, nil
}

func (s *dispatchServiceServer) SendDispatch(ctx context.Context, req *grpc_gen.SendDispatchRequest) (*grpc_gen.SendDispatchResponse, error) {
	s.log("SendDispatch", req.String())
	err := s.service.SendDispatch(ctx, req.DispatchId)
	if errors.Is(err, service_errors.NotFoundErr) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpc_gen.SendDispatchResponse{}, nil
}

func (s *dispatchServiceServer) GetAllDispatches(ctx context.Context, req *grpc_gen.GetAllDispatchesRequest) (*grpc_gen.GetAllDispatchesResponse, error) {
	s.log("GetAllDispatches", req.String())
	allDispatches, err := s.service.GetAllDispatches(ctx)
	if errors.Is(err, service_errors.NotFoundErr) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, service_errors.InvalidArgumentErr) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if len(allDispatches) == 0 {
		return &grpc_gen.GetAllDispatchesResponse{Dispatches: []*grpc_gen.DispatchData{}}, nil
	}

	dispatches := make([]*grpc_gen.DispatchData, 0, len(allDispatches))
	for _, dsptch := range allDispatches {
		dispatches = append(dispatches, ToProtoDispatch(dsptch))
	}

	return &grpc_gen.GetAllDispatchesResponse{Dispatches: dispatches}, nil
}

func ToProtoDispatch(data repo.DispatchData) *grpc_gen.DispatchData {
	return &grpc_gen.DispatchData{
		Id:                 data.Id,
		Label:              data.Label,
		SendAt:             data.SendAt,
		CountOfSubscribers: int64(data.CountOfSubscribers),
	}
}
