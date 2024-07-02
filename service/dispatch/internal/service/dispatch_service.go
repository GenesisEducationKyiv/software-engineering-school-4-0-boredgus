package service

import (
	"bytes"
	"context"

	"errors"
	"html/template"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/emails"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service/deps"
	service_errors "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service/err"
)

type dispatchService struct {
	log          config.Logger
	userRepo     deps.UserRepo
	subRepo      deps.SubRepo
	dispatchRepo deps.DispatchRepo
	mailman      deps.Mailman
	csClient     deps.CurrencyServiceClient
}

func NewDispatchService(
	logger config.Logger,
	mailman deps.Mailman,
	currencyService deps.CurrencyServiceClient,
	userRepo deps.UserRepo,
	subRepo deps.SubRepo,
	dispatchRepo deps.DispatchRepo,
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

func (s *dispatchService) GetAllDispatches(ctx context.Context) ([]deps.DispatchData, error) {
	return s.dispatchRepo.GetAllDispatches(ctx)
}

func (s *dispatchService) SubscribeForDispatch(ctx context.Context, email, dispatchId string) error {
	_, err := s.dispatchRepo.GetDispatchByID(ctx, dispatchId)
	if err != nil {
		return err
	}

	if err = s.userRepo.CreateUser(ctx, email); err != nil && !errors.Is(err, service_errors.UniqueViolationErr) {
		return err
	}

	// TODO: send welcome email if creation of subscription was successful
	return s.subRepo.CreateSubscription(ctx, deps.SubscriptionData{Email: email, Dispatch: dispatchId})
}

var TemplateParseErr = errors.New("template error")

func (s *dispatchService) parseHTMLTemplate(templateName string, data any) ([]byte, error) {
	templateFile := emails.PathToTemplate(templateName + ".html")
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

	curencyRates, err := s.csClient.Convert(ctx, dispatch.Details.BaseCurrency, dispatch.Details.TargetCurrencies)
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

	return s.mailman.Send(deps.Email{
		To:       subscribers,
		Subject:  dispatch.Label,
		HTMLBody: string(htmlContent),
	})
}
