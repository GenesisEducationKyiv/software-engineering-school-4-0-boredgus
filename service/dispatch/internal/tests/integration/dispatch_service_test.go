package tests

import (
	"context"
	"database/sql"
	"slices"
	"testing"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/repo"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/tests"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/tests/stubs"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/tests/testdata"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type (
	DispatchService interface {
		SubscribeForDispatch(ctx context.Context, email, dispatchId string) error
		UnsubscribeFromDispatch(ctx context.Context, email, dispatchId string) error
	}

	DispatchServiceSuite struct {
		suite.Suite
		ctx             context.Context
		dispatchService DispatchService

		pgContainer  *tests.PostgresContainer
		dispatchRepo service.DispatchRepo
		dbConnection *sql.DB

		logger config.Logger
		broker *stubs.BrokerStub
	}
)

func (s *DispatchServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	pgContainer, err := tests.CreatePostgresContainer(s.ctx)
	s.NoErrorf(err, "failed to create postgres container")
	s.pgContainer = pgContainer

	dbConnection, err := db.NewPostrgreSQL(
		s.pgContainer.ConnectionString,
		db.PostgeSQLMigrationsUp,
	)
	s.NoError(err)
	s.dbConnection = dbConnection
	storage := repo.NewStore(dbConnection, db.IsPqError)
	s.dispatchRepo = repo.NewDispatchRepo(storage)
	s.broker = stubs.NewBrokerStub()

	s.dispatchService = service.NewDispatchService(
		repo.NewUserRepo(storage),
		repo.NewSubRepo(storage),
		s.dispatchRepo,
		s.broker,
	)
}

func (s *DispatchServiceSuite) TearDownSuite() {
	err := s.dbConnection.Close()
	if err != nil {
		s.Fail("failed to close connection to db", err)
	}

	err = s.pgContainer.Terminate(s.ctx)
	if err != nil {
		s.Fail("failed to terminate database container", err)
	}
}

func (s *DispatchServiceSuite) Test_SubscribeForDispatch_SuccessfullyCreatedSubscription() {
	emailToSubscribe := "email_for_created_subscription@gmail.com"
	dispatchID := testdata.USD_UAH_DISPATCH_ID
	ctx := context.Background()

	s.broker.On("CreateSubscription", mock.Anything)
	s.NoError(s.dispatchService.SubscribeForDispatch(ctx, emailToSubscribe, dispatchID))

	subscribers, err := s.dispatchRepo.GetSubscribersOfDispatch(ctx, dispatchID)
	s.NoError(err)
	s.True(slices.Contains(subscribers, emailToSubscribe))
}

func (s *DispatchServiceSuite) Test_SubscribeForDispatch_SuccessfullyRenewedSubscription() {
	data := testdata.CancelledSubscriptionData
	ctx := context.Background()
	s.NoError(s.pgContainer.ExecuteSQLFiles(ctx, data.Filename))
	s.broker.On("RenewSubscription", mock.Anything)

	s.NoError(s.dispatchService.SubscribeForDispatch(ctx, data.Email, data.DispatchID))

	subscribers, err := s.dispatchRepo.GetSubscribersOfDispatch(ctx, data.DispatchID)
	s.NoError(err)
	s.True(slices.Contains(subscribers, data.Email))
}

func (s *DispatchServiceSuite) Test_SubscribeForDispatch_UserAlreadySubscribedForThisDispatch() {
	email := "duplicated_email@gmail.com"
	dispatchId := testdata.USD_UAH_DISPATCH_ID
	ctx := context.Background()

	s.broker.On("CreateSubscription", mock.Anything)
	s.NoError(s.dispatchService.SubscribeForDispatch(ctx, email, dispatchId))
	s.ErrorIs(
		s.dispatchService.SubscribeForDispatch(ctx, email, dispatchId),
		service.AlreadyExistsErr,
	)
}

func (s *DispatchServiceSuite) Test_UnsubscribeFromDispatch_SuccessfullyCancelledSubscription() {
	data := testdata.NewSubscriptionData
	ctx := context.Background()

	s.broker.On("CancelSubscription", mock.Anything)
	s.NoError(s.pgContainer.ExecuteSQLFiles(ctx, data.Filename))

	s.NoError(s.dispatchService.UnsubscribeFromDispatch(ctx, data.Email, data.DispatchID))
}

func (s *DispatchServiceSuite) Test_UnsubscribeFromDispatch_SubscriptionIsAlreadyCancelled() {
	data := testdata.CancelledSubscriptionData2
	ctx := context.Background()
	s.NoError(s.pgContainer.ExecuteSQLFiles(ctx, data.Filename))

	s.NoError(s.dispatchService.UnsubscribeFromDispatch(ctx, data.Email, data.DispatchID))
}

func TestIntegration_DispatchService(t *testing.T) {
	suite.Run(t, new(DispatchServiceSuite))
}
