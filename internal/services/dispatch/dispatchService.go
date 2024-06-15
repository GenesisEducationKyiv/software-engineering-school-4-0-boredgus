package ds

import (
	"context"
	"errors"
	"subscription-api/config"
	db "subscription-api/internal/db"
	e "subscription-api/internal/entities"
	"subscription-api/internal/mailing"
	"subscription-api/internal/services"
)

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

type dispatchService struct {
	store        db.Store
	userRepo     UserRepo
	subRepo      SubRepo
	dispatchRepo DispatchRepo
	htmlParser   HTMLTemplateParser
	mailman      Mailman
	log          config.Logger
}

func NewDispatchService(s db.Store, l config.Logger, smtpParams mailing.SMTPParams) *dispatchService {
	return &dispatchService{
		store:        s,
		userRepo:     db.NewUserRepo(),
		subRepo:      db.NewSubRepo(),
		dispatchRepo: db.NewDispatchRepo(),
		htmlParser:   mailing.NewHTMLTemplateParser(l),
		mailman:      mailing.NewMailman(smtpParams),
		log:          l,
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

		subscribers, err = s.dispatchRepo.GetSubscribersOfDispatch(ctx, d, dispatchId)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	s.log.Infof("fetched dispatch: %+v     subscribers: %v\n\n", dispatch, subscribers)

	htmlContent, err := s.htmlParser.Parse(dispatch.TemplateName, dispatch)
	if err != nil {
		return err
	}

	return s.mailman.Send(mailing.Email{
		To:       subscribers,
		Subject:  dispatch.Label,
		HTMLBody: string(htmlContent),
	})
}
