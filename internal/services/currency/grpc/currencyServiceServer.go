package grpc

import (
	"context"
	"errors"
	"subscription-api/config"
	"subscription-api/internal/entities"
	"subscription-api/internal/services"
	cs "subscription-api/internal/services/currency"
	pb_cs "subscription-api/pkg/grpc/currency_service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type currencyServiceServer struct {
	pb_cs.UnimplementedCurrencyServiceServer
	s cs.CurrencyService
	l config.Logger
}

func NewCurrencyServiceServer(s cs.CurrencyService, l config.Logger) pb_cs.CurrencyServiceServer {
	return &currencyServiceServer{s: s, l: l}
}

func (s *currencyServiceServer) log(method string, req any) {
	s.l.Infof("CurrencyService.%v(%+v)", method, req)
}

func (s *currencyServiceServer) Convert(ctx context.Context, req *pb_cs.ConvertRequest) (*pb_cs.ConvertResponse, error) {
	s.log("Convert", req)
	rates, err := s.s.Convert(ctx, cs.ConvertCurrencyParams{
		Base:   entities.Currency(req.BaseCurrency),
		Target: entities.FromString(req.TargetCurrencies),
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

	return &pb_cs.ConvertResponse{BaseCurrency: req.BaseCurrency, Rates: res}, nil
}
