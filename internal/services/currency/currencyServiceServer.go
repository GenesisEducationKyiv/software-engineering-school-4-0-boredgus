package currency_service

import (
	"context"
	"errors"
	"subscription-api/config"
	"subscription-api/internal/services"
	currency_grpc "subscription-api/internal/services/currency/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type currencyServiceServer struct {
	currency_grpc.UnimplementedCurrencyServiceServer
	s services.CurrencyService
	l config.Logger
}

func NewCurrencyServiceServer(s services.CurrencyService, l config.Logger) *currencyServiceServer {
	return &currencyServiceServer{s: s, l: l}
}

func (s *currencyServiceServer) log(method string, req any) {
	s.l.Infof("CurrencyService.%v(%+v)", method, req)
}

func (s *currencyServiceServer) Convert(ctx context.Context, req *currency_grpc.ConvertRequest) (*currency_grpc.ConvertResponse, error) {
	s.log("Convert", req)
	rates, err := s.s.Convert(ctx, services.ConvertCurrencyParams{
		Base:   req.BaseCurrency,
		Target: req.TargetCurrencies,
	})
	if errors.Is(err, services.InvalidArgumentErr) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if errors.Is(err, services.FailedPreconditionErr) {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	res := make(map[string]float64)
	for currency, rate := range rates {
		res[string(currency)] = rate
	}

	return &currency_grpc.ConvertResponse{
			BaseCurrency: req.BaseCurrency,
			Rates:        res,
		},
		nil
}
