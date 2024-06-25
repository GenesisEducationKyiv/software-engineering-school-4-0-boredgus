package dispatch_service

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"subscription-api/config"
	db "subscription-api/internal/db"
	"subscription-api/internal/entities"
	"subscription-api/internal/mailing"
	"subscription-api/internal/services"
)

type Store interface {
	db.Database
	IsError(error, db.Error) bool
}

type UserRepo interface {
	CreateUser(ctx context.Context, email string) error
}

type SubRepo interface {
	CreateSubscription(ctx context.Context, args db.SubscriptionData) error
}

type DispatchRepo interface {
	GetDispatchByID(ctx context.Context, dispatchId string) (entities.CurrencyDispatch, error)
	GetSubscribersOfDispatch(ctx context.Context, dispatchId string) ([]string, error)
	GetAllDispatches(ctx context.Context) ([]services.DispatchData, error)
}

type Mailman interface {
	Send(email mailing.Email) error
}

type CurrencyServiceClient interface {
	Convert(ctx context.Context, params services.ConvertCurrencyParams) (map[string]float64, error)
}

type dispatchService struct {
	log          config.Logger
	userRepo     UserRepo
	subRepo      SubRepo
	dispatchRepo DispatchRepo
	mailman      Mailman
	csClient     CurrencyServiceClient
}

func NewDispatchService(
	logger config.Logger,
	mailman Mailman,
	currencyService CurrencyServiceClient,
	userRepo UserRepo,
	subRepo SubRepo,
	dispatchRepo DispatchRepo,
) *dispatchService {
	return &dispatchService{
		userRepo:     userRepo,
		subRepo:      subRepo,
		dispatchRepo: dispatchRepo,
		mailman:      mailman,
		csClient:     currencyService,
		log:          logger,
	}
}

func (s *dispatchService) GetAllDispatches(ctx context.Context) ([]services.DispatchData, error) {
	return s.dispatchRepo.GetAllDispatches(ctx)
}

func (s *dispatchService) SubscribeForDispatch(ctx context.Context, email, dispatchId string) error {
	_, err := s.dispatchRepo.GetDispatchByID(ctx, dispatchId)
	if err != nil {
		return err
	}

	if err = s.userRepo.CreateUser(ctx, email); err != nil && !errors.Is(err, services.UniqueViolationErr) {
		return err
	}

	// TODO: send welcome email if creation of subscription was successful
	return s.subRepo.CreateSubscription(ctx, db.SubscriptionData{Email: email, Dispatch: dispatchId})
}

var TemplateParseErr = errors.New("template error")

func (s *dispatchService) parseHTMLTemplate(templateName string, data any) ([]byte, error) {
	templateFile := mailing.PathToTemplate(templateName + ".html")
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		s.log.Errorf("failed to parse html template %s: %v", templateName, err)

		return nil, errors.Join(TemplateParseErr, err)
	}
	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {
		s.log.Errorf("failed to execute html template %s: %v", templateName, err)

		return nil, errors.Join(TemplateParseErr, err)
	}

	return buffer.Bytes(), nil
}

type ExchangeRateTemplateParams struct {
	BaseCurrency string
	Rates        map[string]float64
}

func (s *dispatchService) SendDispatch(ctx context.Context, dispatchId string) error {
	dispatch, err := s.dispatchRepo.GetDispatchByID(ctx, dispatchId)
	if err != nil {
		return err
	}
	if dispatch.CountOfSubscribers == 0 {
		return nil
	}

	subscribers, err := s.dispatchRepo.GetSubscribersOfDispatch(ctx, dispatchId)
	if err != nil {
		return err
	}

	curencyRates, err := s.csClient.Convert(ctx, services.ConvertCurrencyParams{
		Base:   dispatch.Details.BaseCurrency,
		Target: dispatch.Details.TargetCurrencies,
	})
	if err != nil {
		return err
	}

	htmlContent, err := s.parseHTMLTemplate(dispatch.TemplateName, ExchangeRateTemplateParams{
		BaseCurrency: dispatch.Details.BaseCurrency,
		Rates:        curencyRates,
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
