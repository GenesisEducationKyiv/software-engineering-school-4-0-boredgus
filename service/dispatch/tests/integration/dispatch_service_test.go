package tests

import (
	"context"
	"database/sql"
	"slices"
	"testing"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/mailing"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/repo"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service"
	service_err "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service/err"
	pkg_sql "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/sql"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/tests"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/tests/stubs"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/tests/testdata"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/shared/db/postgres"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type (
	CurrencyServiceServer interface {
		Stop()
	}

	Store interface {
		repo.DB
		service.Store
	}

	DispatchServiceSuite struct {
		suite.Suite
		ctx             context.Context
		dispatchService service.DispatchService

		pgContainer  *tests.PostgresContainer
		dispatchRepo service.DispatchRepo
		dbConnection *sql.DB

		logger  config.Logger
		mailman *stubs.MailmanStub
	}
)

func (s *DispatchServiceSuite) SetupSuite() {
	s.logger = config.InitLogger(config.TestMode)

	s.mailman = stubs.NewMailmanStub()
	s.ctx = context.Background()

	pgContainer, err := tests.CreatePostgresContainer(s.ctx)
	s.NoErrorf(err, "failed to create postgres container")
	s.pgContainer = pgContainer

	dbConnection, err := postgres.NewPostrgreSQL(
		s.pgContainer.ConnectionString,
		pkg_sql.PostgeSQLMigrationsUp(nil),
	)
	s.NoError(err)
	s.dbConnection = dbConnection
	storage := repo.NewStore(dbConnection, postgres.IsPqError)
	s.dispatchRepo = repo.NewDispatchRepo(storage)

	s.dispatchService = service.NewDispatchService(
		s.logger,
		s.mailman,
		stubs.NewCurrencyServiceClient(),
		repo.NewUserRepo(storage),
		repo.NewSubRepo(storage),
		s.dispatchRepo,
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

func (s *DispatchServiceSuite) Test_GetAllDispatches() {
	ctx := context.Background()

	dispatches, err := s.dispatchService.GetAllDispatches(ctx)

	s.NoError(err)
	s.Equal(1, len(dispatches))
	s.Equal(dispatches[0].Id, testdata.USD_UAH_DISPATCH_ID)
}

func (s *DispatchServiceSuite) Test_SendDispatch() {
	ctx := context.Background()
	s.NoError(s.pgContainer.ExecuteSQLFiles(ctx, "add_couple_of_subscribers_for_usd_uah_dispatch"))

	s.mailman.On("Send", mock.Anything).Return(nil)

	s.NoError(s.dispatchService.SendDispatch(ctx, testdata.USD_UAH_DISPATCH_ID))

	s.Equal(1, len(s.mailman.Calls))
	actualEmailReceivers := s.mailman.Calls[0].Arguments.Get(0).(mailing.Email).To
	expectedEmailReceivers := testdata.SubscribersOfUSDUAHDispatch
	s.Equal(actualEmailReceivers, expectedEmailReceivers)
}

func (s *DispatchServiceSuite) Test_SubscribeForDispatch_Success() {
	emailToSubscribe := "email_1@gmail.com"
	dispatchID := testdata.USD_UAH_DISPATCH_ID
	ctx := context.Background()

	s.NoError(s.dispatchService.SubscribeForDispatch(ctx, emailToSubscribe, dispatchID))

	subscribers, err := s.dispatchRepo.GetSubscribersOfDispatch(ctx, dispatchID)
	s.NoError(err)
	s.True(slices.Contains(subscribers, emailToSubscribe))
}

func (s *DispatchServiceSuite) Test_SubscribeForDispatch_UserAlreadySubscribedForThisDispatch() {
	email := "email_2@gmail.com"
	dispatchId := testdata.USD_UAH_DISPATCH_ID
	ctx := context.Background()

	s.NoError(s.dispatchService.SubscribeForDispatch(ctx, email, dispatchId))
	s.ErrorIs(
		s.dispatchService.SubscribeForDispatch(ctx, email, dispatchId),
		service_err.UniqueViolationErr,
	)
}

func TestIntegration_DispatchService(t *testing.T) {
	suite.Run(t, new(DispatchServiceSuite))
}
