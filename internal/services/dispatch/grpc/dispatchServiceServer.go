package grpc

import (
	"context"
	"errors"
	"subscription-api/config"
	"subscription-api/internal/services"
	ss "subscription-api/internal/services"
	pb_ds "subscription-api/pkg/grpc/dispatch_service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type dispatchServiceServer struct {
	s ss.DispatchService
	l config.Logger
	pb_ds.UnimplementedDispatchServiceServer
}

func NewDispatchServiceServer(s ss.DispatchService, l config.Logger) *dispatchServiceServer {
	return &dispatchServiceServer{s: s, l: l}
}

func (s *dispatchServiceServer) log(method string, req any) {
	s.l.Infof("DispatchService.%v(%+v)", method, req)
}

func (s *dispatchServiceServer) SubscribeForDispatch(ctx context.Context, req *pb_ds.SubscribeForDispatchRequest) (*pb_ds.SubscribeForDispatchResponse, error) {
	s.log("SubscribeForDispatch", req.String())
	err := s.s.SubscribeForDispatch(ctx, req.Email, req.DispatchId)
	if errors.Is(err, services.UniqueViolationErr) {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}
	if errors.Is(err, services.NotFoundErr) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb_ds.SubscribeForDispatchResponse{}, nil
}

func (s *dispatchServiceServer) SendDispatch(ctx context.Context, req *pb_ds.SendDispatchRequest) (*pb_ds.SendDispatchResponse, error) {
	s.log("SendDispatch", req.String())
	err := s.s.SendDispatch(ctx, req.DispatchId)
	if errors.Is(err, services.NotFoundErr) {
		return nil, status.Error(codes.Canceled, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb_ds.SendDispatchResponse{}, nil
}

func (s *dispatchServiceServer) GetAllDispatches(ctx context.Context, req *pb_ds.GetAllDispatchesRequest) (*pb_ds.GetAllDispatchesResponse, error) {
	s.log("GetAllDispatches", req.String())
	d, err := s.s.GetAllDispatches(ctx)
	if errors.Is(err, services.NotFoundErr) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, services.InvalidArgumentErr) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	dispatches := make([]*pb_ds.DispatchData, 0, len(d))
	for _, dsptch := range d {
		dispatches = append(dispatches, dsptch.ToProto())
	}

	return &pb_ds.GetAllDispatchesResponse{Dispatches: dispatches}, nil
}
