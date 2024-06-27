package tests

import (
	"context"
	"database/sql"
	"net"
	"slices"
	"subscription-api/config"
	store "subscription-api/internal/db"
	"subscription-api/internal/mailing"
	"subscription-api/internal/services"
	sql_pkg "subscription-api/internal/sql"
	"subscription-api/internal/tests"
	"subscription-api/internal/tests/stubs"
	"subscription-api/internal/tests/testdata"
	"subscription-api/pkg/db"
	"testing"

	currency_service "subscription-api/internal/services/currency"
	currency_server "subscription-api/internal/services/currency/server"
	currency_grpc "subscription-api/internal/services/currency/server/grpc"
	dispatch_service "subscription-api/internal/services/dispatch"

	grpc_client "subscription-api/pkg/grpc"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	CurrencyServiceURL string = "localhost:4040"
)

type (
	CurrencyServiceServer interface {
		Stop()
	}

	Store interface {
		store.DB
		dispatch_service.Store
	}

	DispatchServiceSuite struct {
		suite.Suite
		ctx             context.Context
		dispatchService services.DispatchService

		currencyServiceServer CurrencyServiceServer

		pgContainer  *tests.PostgresContainer
		dispatchRepo dispatch_service.DispatchRepo
		dbConnection *sql.DB

		logger  config.Logger
		mailman *stubs.MailmanStub
	}
)

func (s *DispatchServiceSuite) startCurrencyServiceServer(url string) {
	lis, err := net.Listen("tcp", url)
	s.NoErrorf(err, "failed to listen url")

	currencyAPI := stubs.NewCurrencyAPIClient()
	currencyService := currency_service.NewCurrencyService(currencyAPI)

	server := grpc.NewServer()
	currency_grpc.RegisterCurrencyServiceServer(
		server,
		currency_server.NewCurrencyServiceServer(currencyService, s.logger),
	)
	s.currencyServiceServer = server

	go func() {
		if err := server.Serve(lis); err != nil {
			s.Fail(err.Error())
		}
	}()
}

func (s *DispatchServiceSuite) SetupSuite() {
	s.logger = config.InitLogger(config.TestMode)

	s.startCurrencyServiceServer(CurrencyServiceURL)
	connToCurrencyService, err := grpc.NewClient(
		CurrencyServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.NoErrorf(err, "failed to create connection to currency service")

	s.mailman = stubs.NewMailmanStub()
	s.ctx = context.Background()

	pgContainer, err := tests.CreatePostgresContainer(s.ctx)
	s.NoErrorf(err, "failed to create postgres container")
	s.pgContainer = pgContainer

	dbConnection, err := db.NewPostrgreSQL(
		s.pgContainer.ConnectionString,
		sql_pkg.PostgeSQLMigrationsUp(nil),
	)
	s.NoError(err)
	s.dbConnection = dbConnection
	storage := store.NewStore(dbConnection, db.IsPqError)
	s.dispatchRepo = store.NewDispatchRepo(storage)

	s.dispatchService = dispatch_service.NewDispatchService(
		s.logger,
		s.mailman,
		grpc_client.NewCurrencyServiceClient(connToCurrencyService),
		store.NewUserRepo(storage),
		store.NewSubRepo(storage),
		s.dispatchRepo,
	)
}

func (s *DispatchServiceSuite) TearDownSuite() {
	s.currencyServiceServer.Stop()

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
	s.Equal(dispatches[0].Id, services.USD_UAH_DISPATCH_ID)
}

func (s *DispatchServiceSuite) Test_SendDispatch() {
	ctx := context.Background()
	s.NoError(s.pgContainer.ExecuteSQLFiles(ctx, "add_couple_of_subscribers_for_usd_uah_dispatch"))

	s.mailman.On("Send", mock.Anything).Return(nil)

	s.NoError(s.dispatchService.SendDispatch(ctx, services.USD_UAH_DISPATCH_ID))
	s.Equal(1, len(s.mailman.Calls))
	actualEmailReceivers := s.mailman.Calls[0].Arguments.Get(0).(mailing.Email).To
	expectedEmailReceivers := testdata.SubscribersOfUSDUAHDispatch
	s.Equal(actualEmailReceivers, expectedEmailReceivers)
}

func (s *DispatchServiceSuite) Test_SubscribeForDispatch_Success() {
	emailToSubscribe := "email_1@gmail.com"
	dispatchID := services.USD_UAH_DISPATCH_ID
	ctx := context.Background()

	s.NoError(s.dispatchService.SubscribeForDispatch(ctx, emailToSubscribe, dispatchID))

	subscribers, err := s.dispatchRepo.GetSubscribersOfDispatch(ctx, dispatchID)
	s.NoError(err)
	s.True(slices.Contains(subscribers, emailToSubscribe))
}

func (s *DispatchServiceSuite) Test_SubscribeForDispatch_UserAlreadySubscribedForThisDispatch() {
	email := "email_2@gmail.com"
	dispatchId := services.USD_UAH_DISPATCH_ID
	ctx := context.Background()

	s.NoError(s.dispatchService.SubscribeForDispatch(ctx, email, dispatchId))
	s.ErrorIs(
		s.dispatchService.SubscribeForDispatch(ctx, email, dispatchId),
		services.UniqueViolationErr,
	)
}

func TestIntegration_DispatchService(t *testing.T) {
	suite.Run(t, new(DispatchServiceSuite))
}
