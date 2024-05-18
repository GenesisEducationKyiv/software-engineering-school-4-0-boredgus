package grpc

import (
	"context"
	"errors"
	"subscription-api/internal/entities"
	cs "subscription-api/internal/services/currency"
	pb_cs "subscription-api/pkg/grpc/currency_service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type currencyServiceServer struct {
	pb_cs.UnimplementedCurrencyServiceServer
	s cs.CurrencyService
}

func NewCurrencyServiceServer(s cs.CurrencyService) pb_cs.CurrencyServiceServer {
	return &currencyServiceServer{s: s}
}

func (s *currencyServiceServer) Convert(ctx context.Context, req *pb_cs.ConvertRequest) (*pb_cs.ConvertResponse, error) {
	rates, err := s.s.Convert(ctx, cs.ConvertParams{
		From: entities.Currency(req.BaseCurrency),
		To:   entities.FromString(req.TargetCurrencies),
	})
	if errors.Is(err, cs.InvalidArgumentErr) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if errors.Is(err, cs.FailedPreconditionErr) {
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
