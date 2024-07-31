package server

import (
	"context"
	"errors"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/entities"
	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/grpc/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	DispatchService interface {
		SubscribeForDispatch(ctx context.Context, email, dispatch string) (*entities.Subscription, error)
		UnsubscribeFromDispatch(ctx context.Context, email, dispatch string) (*entities.Subscription, error)
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

	subscription, err := s.service.SubscribeForDispatch(ctx, req.Email, req.DispatchId)
	if errors.Is(err, service.ErrUniqueViolation) {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}
	if errors.Is(err, service.ErrNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpc_gen.SubscribeForDispatchResponse{
		Subscription: subscriptionToProto(subscription),
	}, nil
}

func (s *dispatchServiceServer) UnsubscribeFromDispatch(ctx context.Context, req *grpc_gen.UnsubscribeFromDispatchRequest) (*grpc_gen.UnsubscribeFromDispatchResponse, error) {
	s.log("UnsubscribeFromDispatch", req.String())

	subscription, err := s.service.UnsubscribeFromDispatch(ctx, req.Email, req.DispatchId)
	if errors.Is(err, service.ErrNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpc_gen.UnsubscribeFromDispatchResponse{
		Subscription: subscriptionToProto(subscription),
	}, nil
}

var statusToProto = map[entities.SubscriptionStatus]grpc_gen.SubscriptionStatus{
	entities.SubscriptionStatusActive:    grpc_gen.SubscriptionStatus_CREATED,
	entities.SubscriptionStatusCancelled: grpc_gen.SubscriptionStatus_CANCELLED,
}

func subscriptionToProto(s *entities.Subscription) *grpc_gen.Subscription {
	return &grpc_gen.Subscription{
		DispatchID:  s.DispatchID,
		Email:       s.Email,
		BaseCcy:     s.BaseCcy,
		TargetCcies: s.TargetCcies,
		Status:      statusToProto[s.Status],
		SendAt:      timestamppb.New(s.SendAt.UTC()),
	}
}
