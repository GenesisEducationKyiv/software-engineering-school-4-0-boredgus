package server

import (
	"context"
	"errors"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/config"
	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/grpc/gen"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CustomerService interface {
	CreateCustomer(ctx context.Context, email string) error
}

type customerServiceServer struct {
	grpc_gen.UnimplementedCustomerServiceServer
	service CustomerService
	logger  config.Logger
}

func NewCustomerServiceServer(s CustomerService, l config.Logger) *customerServiceServer {
	return &customerServiceServer{service: s, logger: l}
}

func (s *customerServiceServer) log(method string, req any) {
	s.logger.Infof("CustomerService.%v(%+v)", method, req)
}

func (s *customerServiceServer) CreateCustomer(ctx context.Context, req *grpc_gen.CreateCustomerRequest) (*emptypb.Empty, error) {
	s.log("CreateCustomer", req)
	err := s.service.CreateCustomer(ctx, req.Email)
	if errors.Is(err, service.AlreadyExistsErr) {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
