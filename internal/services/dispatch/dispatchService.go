package ds

import (
	"context"
	"errors"
	"subscription-api/config"
	db "subscription-api/internal/db"
	e "subscription-api/internal/entities"
	"subscription-api/internal/mailing"
	"subscription-api/internal/services"
	pb_cs "subscription-api/pkg/grpc/currency_service"

	"google.golang.org/grpc"
)

type Store interface {
	WithTx(ctx context.Context, f func(db.DB) error) error
}

type UserRepo interface {
	CreateUser(ctx context.Context, db db.DB, email string) error
}

type SubRepo interface {
	CreateSubscription(ctx context.Context, db db.DB, args db.SubscriptionData) error
}

type DispatchRepo interface {
	GetDispatchByID(ctx context.Context, db db.DB, dispatchId string) (e.CurrencyDispatch, error)
	GetSubscribersOfDispatch(ctx context.Context, db db.DB, dispatchId string) ([]string, error)
	GetAllDispatches(ctx context.Context, db db.DB) ([]e.CurrencyDispatch, error)
}

type HTMLTemplateParser interface {
	Parse(templateName string, data any) ([]byte, error)
}

type Mailman interface {
	Send(email mailing.Email) error
}

type CurrencyServiceClient interface {
	Convert(ctx context.Context, in *pb_cs.ConvertRequest, opts ...grpc.CallOption) (*pb_cs.ConvertResponse, error)
}

type dispatchService struct {
	log          config.Logger
	store        Store
	userRepo     UserRepo
	subRepo      SubRepo
	dispatchRepo DispatchRepo
	htmlParser   HTMLTemplateParser
	mailman      Mailman
	csClient     CurrencyServiceClient
}

type DispatchServiceParams struct {
	Store           Store
	Logger          config.Logger
	Mailman         Mailman
	CurrencyService CurrencyServiceClient
}

func NewDispatchService(params DispatchServiceParams) *dispatchService {
	return &dispatchService{
		store:        params.Store,
		userRepo:     db.NewUserRepo(),
		subRepo:      db.NewSubRepo(),
		dispatchRepo: db.NewDispatchRepo(),
		htmlParser:   mailing.NewHTMLTemplateParser(params.Logger),
		mailman:      params.Mailman,
		csClient:     params.CurrencyService,
		log:          params.Logger,
	}
}

func (s *dispatchService) GetAllDispatches(ctx context.Context) ([]e.CurrencyDispatch, error) {
	var dispatches []e.CurrencyDispatch
	err := s.store.WithTx(ctx, func(db db.DB) error {
		d, err := s.dispatchRepo.GetAllDispatches(ctx, db)
		if err != nil {
			return err
		}
		dispatches = make([]e.CurrencyDispatch, 0, len(d))
		for _, dsptch := range d {
			dispatches = append(dispatches, dsptch)
		}

		return nil
	})

	return dispatches, err
}

func (s *dispatchService) SubscribeForDispatch(ctx context.Context, email, dispatchId string) error {
	if err := s.store.WithTx(ctx, func(d db.DB) error {
		_, err := s.dispatchRepo.GetDispatchByID(ctx, d, dispatchId)
		if err != nil {
			return err
		}

		if err = s.userRepo.CreateUser(ctx, d, email); err != nil && !errors.Is(err, services.UniqueViolationErr) {
			return err
		}

		return s.subRepo.CreateSubscription(ctx, d, db.SubscriptionData{Email: email, Dispatch: dispatchId})
	}); err != nil {
		return err
	}
	// TODO: send welcome email
	return nil
}

type ExchangeRates struct {
	BaseCurrency string
	Rates        map[string]float64
}

func (s *dispatchService) SendDispatch(ctx context.Context, dispatchId string) error {
	var dispatch e.CurrencyDispatch
	var subscribers []string

	if err := s.store.WithTx(ctx, func(d db.DB) error {
		dsptch, err := s.dispatchRepo.GetDispatchByID(ctx, d, dispatchId)
		if err != nil {
			return err
		}
		if dsptch.CountOfSubscribers == 0 {
			return nil
		}
		dispatch = dsptch
		subscribers, err = s.dispatchRepo.GetSubscribersOfDispatch(ctx, d, dispatchId)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	if len(subscribers) == 0 {
		return nil
	}

	resp, err := s.csClient.Convert(ctx, &pb_cs.ConvertRequest{
		BaseCurrency:     dispatch.Details.BaseCurrency,
		TargetCurrencies: dispatch.Details.TargetCurrencies,
	})
	if err != nil {
		return err
	}

	htmlContent, err := s.htmlParser.Parse(dispatch.TemplateName, ExchangeRates{
		BaseCurrency: resp.BaseCurrency,
		Rates:        resp.Rates,
	})
	if err != nil {
		return err
	}

	return s.mailman.Send(mailing.Email{
		To:       subscribers,
		Subject:  dispatch.Label,
		HTMLBody: string(htmlContent),
	})
}
