package grpc

import (
	"context"
	"errors"
	"subscription-api/config"
	"subscription-api/internal/services"
	ds "subscription-api/internal/services/dispatch"
	pb_ds "subscription-api/pkg/grpc/dispatch_service"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type dispatchServiceServer struct {
	s ds.DispatchService
	pb_ds.UnimplementedDispatchServiceServer
}

func NewDispatchServiceServer(s ds.DispatchService) pb_ds.DispatchServiceServer {
	return &dispatchServiceServer{s: s}
}

func (s *dispatchServiceServer) log(method string, req any) {
	config.Log().Infof("DispatchService.%v(%+v)", method, req)
}

func (s *dispatchServiceServer) SubscribeFor(ctx context.Context, req *pb_ds.SubscribeForRequest) (*pb_ds.SubscribeForResponse, error) {
	s.log("SubscribeFor", req.String())
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

	return &pb_ds.SubscribeForResponse{}, nil
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

func (s *dispatchServiceServer) GetDispatch(ctx context.Context, req *pb_ds.GetDispatchRequest) (*pb_ds.GetDispatchResponse, error) {
	s.log("GetDispatch", req.String())
	d, err := s.s.GetDispatch(ctx, req.DispatchId)
	if errors.Is(err, services.NotFoundErr) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, services.InvalidArgumentErr) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb_ds.GetDispatchResponse{
		Dispatch: &pb_ds.DispatchData{
			Id:                 d.Id,
			BaseCurrency:       d.Details.BaseCurrency,
			TargetCurrencies:   d.Details.TargetCurrencies,
			SendAt:             d.SendAt.Format(time.TimeOnly),
			CountOfSubscribers: int64(d.CountOfSubscribers),
		},
	}, nil
}
