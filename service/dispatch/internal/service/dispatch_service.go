package service

import (
	"context"

	"errors"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service/deps"
	service_errors "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service/err"
)

type dispatchService struct {
	log          config.Logger
	userRepo     deps.UserRepo
	subRepo      deps.SubRepo
	dispatchRepo deps.DispatchRepo
	broker       deps.Broker
}

func NewDispatchService(
	logger config.Logger,
	userRepo deps.UserRepo,
	subRepo deps.SubRepo,
	dispatchRepo deps.DispatchRepo,
	broker deps.Broker,
) *dispatchService {
	return &dispatchService{
		userRepo:     userRepo,
		subRepo:      subRepo,
		dispatchRepo: dispatchRepo,
		log:          logger,
		broker:       broker,
	}
}

func (s *dispatchService) GetAllDispatches(ctx context.Context) ([]deps.DispatchData, error) {
	return s.dispatchRepo.GetAllDispatches(ctx)
}

func (s *dispatchService) SubscribeForDispatch(ctx context.Context, email, dispatchId string) error {
	dispatchData, err := s.dispatchRepo.GetDispatchByID(ctx, dispatchId)
	if err != nil {
		return err
	}

	if err = s.userRepo.CreateUser(ctx, email); err != nil && !errors.Is(err, service_errors.UniqueViolationErr) {
		return err
	}

	if err = s.subRepo.CreateSubscription(ctx, deps.SubscriptionData{Email: email, Dispatch: dispatchId}); err != nil {
		return err
	}

	s.broker.CreateSubscription(deps.Subscription{
		DispatchID:  dispatchId,
		Sources:     map[string]string{"email": email},
		BaseCcy:     dispatchData.Details.BaseCurrency,
		TargetCcies: dispatchData.Details.TargetCurrencies,
		SendAt:      dispatchData.SendAt,
	})

	return nil
}
