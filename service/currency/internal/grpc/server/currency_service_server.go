package grpc_server

import (
	"context"
	"errors"

	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/grpc/gen"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/pkg/logger"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/currency/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CurrencyService interface {
	Convert(ctx context.Context, params service.ConvertCurrencyParams) (map[string]float64, error)
}

type currencyServiceServer struct {
	grpc_gen.UnimplementedCurrencyServiceServer
	service CurrencyService
	logger  logger.Logger
}

func NewCurrencyServiceServer(s CurrencyService, l logger.Logger) *currencyServiceServer {
	return &currencyServiceServer{service: s, logger: l}
}

func (s *currencyServiceServer) log(method string, req any) {
	s.logger.Infof("CurrencyService.%v(%+v)", method, req)
}

func (s *currencyServiceServer) Convert(ctx context.Context, req *grpc_gen.ConvertRequest) (*grpc_gen.ConvertResponse, error) {
	s.log("Convert", req)
	rates, err := s.service.Convert(ctx, service.ConvertCurrencyParams{
		Base:   req.BaseCurrency,
		Target: req.TargetCurrencies,
	})
	if errors.Is(err, service.InvalidArgumentErr) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if errors.Is(err, service.FailedPreconditionErr) {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	res := make(map[string]float64)
	for currency, rate := range rates {
		res[string(currency)] = rate
	}

	return &grpc_gen.ConvertResponse{
			BaseCurrency: req.BaseCurrency,
			Rates:        res,
		},
		nil
}
