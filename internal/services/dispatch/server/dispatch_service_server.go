package dispatch_server

import (
	"context"
	"errors"
	"subscription-api/config"
	"subscription-api/internal/services"
	grpc "subscription-api/internal/services/dispatch/server/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type dispatchServiceServer struct {
	service services.DispatchService
	logger  config.Logger
	grpc.UnimplementedDispatchServiceServer
}

func NewDispatchServiceServer(s services.DispatchService, l config.Logger) *dispatchServiceServer {
	return &dispatchServiceServer{service: s, logger: l}
}

func (s *dispatchServiceServer) log(method string, req any) {
	s.logger.Infof("DispatchService.%v(%+v)", method, req)
}

func (s *dispatchServiceServer) SubscribeForDispatch(ctx context.Context, req *grpc.SubscribeForDispatchRequest) (*grpc.SubscribeForDispatchResponse, error) {
	s.log("SubscribeForDispatch", req.String())
	err := s.service.SubscribeForDispatch(ctx, req.Email, req.DispatchId)
	if errors.Is(err, services.UniqueViolationErr) {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}
	if errors.Is(err, services.NotFoundErr) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpc.SubscribeForDispatchResponse{}, nil
}

func (s *dispatchServiceServer) SendDispatch(ctx context.Context, req *grpc.SendDispatchRequest) (*grpc.SendDispatchResponse, error) {
	s.log("SendDispatch", req.String())
	err := s.service.SendDispatch(ctx, req.DispatchId)
	if errors.Is(err, services.NotFoundErr) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpc.SendDispatchResponse{}, nil
}

func (s *dispatchServiceServer) GetAllDispatches(ctx context.Context, req *grpc.GetAllDispatchesRequest) (*grpc.GetAllDispatchesResponse, error) {
	s.log("GetAllDispatches", req.String())
	allDispatches, err := s.service.GetAllDispatches(ctx)
	if errors.Is(err, services.NotFoundErr) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, services.InvalidArgumentErr) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if len(allDispatches) == 0 {
		return &grpc.GetAllDispatchesResponse{Dispatches: []*grpc.DispatchData{}}, nil
	}

	dispatches := make([]*grpc.DispatchData, 0, len(allDispatches))
	for _, dsptch := range allDispatches {
		dispatches = append(dispatches, ToProtoDispatch(dsptch))
	}

	return &grpc.GetAllDispatchesResponse{Dispatches: dispatches}, nil
}

func ToProtoDispatch(data services.DispatchData) *grpc.DispatchData {
	return &grpc.DispatchData{
		Id:                 data.Id,
		Label:              data.Label,
		SendAt:             data.SendAt,
		CountOfSubscribers: int64(data.CountOfSubscribers),
	}
}
