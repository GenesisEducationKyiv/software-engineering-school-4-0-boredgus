package tests

import (
	"context"
	"log"
	"net"
	"slices"
	"subscription-api/config"
	store "subscription-api/internal/db"
	"subscription-api/internal/mailing"
	"subscription-api/internal/services"
	"subscription-api/internal/sql"
	"subscription-api/internal/tests"
	"subscription-api/internal/tests/stubs"
	"subscription-api/pkg/db"
	"subscription-api/pkg/utils"
	"testing"

	currency_service "subscription-api/internal/services/currency"
	currency_grpc "subscription-api/internal/services/currency/grpc"
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

	DispatchServiceSuite struct {
		suite.Suite
		currencyServiceServer CurrencyServiceServer
		service               services.DispatchService
		logger                config.Logger
		pgContainer           *tests.PostgresContainer
		ctx                   context.Context
		mailman               *stubs.MailmanStub
	}
)

func (s *DispatchServiceSuite) startCurrencyServiceServer() {
	lis, err := net.Listen("tcp", CurrencyServiceURL)
	s.NoErrorf(err, "failed to listen url")

	currencyAPI := stubs.NewCurrencyAPIClient()
	currencyService := currency_service.NewCurrencyService(currencyAPI)

	server := grpc.NewServer()
	currency_grpc.RegisterCurrencyServiceServer(
		server,
		currency_service.NewCurrencyServiceServer(currencyService, s.logger),
	)
	s.currencyServiceServer = server

	go func() {
		if err := server.Serve(lis); err != nil {
			s.Fail(err.Error())
		}
	}()
}

func (s *DispatchServiceSuite) getSubscribersOfDispatch(ctx context.Context, dispatchID string) ([]string, error) {
	postgresqlDB, err := db.NewPostrgreSQL(s.pgContainer.ConnectionString)
	s.NoError(err)
	defer postgresqlDB.Close()

	dispatchRepo := store.NewDispatchRepo()

	return dispatchRepo.GetSubscribersOfDispatch(
		ctx,
		store.NewStore(postgresqlDB, db.IsPqError),
		dispatchID)
}

func (s *DispatchServiceSuite) SetupSuite() {
	s.logger = config.InitLogger(config.TestMode)

	s.startCurrencyServiceServer()

	s.ctx = context.Background()

	pgContainer, err := tests.CreatePostgresContainer(s.ctx)
	s.NoErrorf(err, "failed to create postgres container")
	s.pgContainer = pgContainer

	currencyServiceConn, err := grpc.NewClient(
		CurrencyServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.NoErrorf(err, "failed to create connection to currency service")

	s.mailman = stubs.NewMailmanStub()

	s.service = dispatch_service.NewDispatchService(&dispatch_service.DispatchServiceParams{
		Logger: s.logger,
		Store: store.NewStore(
			utils.Must(db.NewPostrgreSQL(
				s.pgContainer.ConnectionString,
				sql.PostgeSQLMigrationsUp(nil),
			)),
			db.IsPqError,
		),
		Mailman:         s.mailman,
		CurrencyService: grpc_client.NewCurrencyServiceClient(currencyServiceConn),
	})
}

func (s *DispatchServiceSuite) TearDownSuite() {
	s.currencyServiceServer.Stop()

	if err := s.pgContainer.Terminate(s.ctx); err != nil {
		log.Fatalf("error terminating database container: %s", err)
	}
}

func (s *DispatchServiceSuite) SetupTest() {}

func (s *DispatchServiceSuite) Test_GetAllDispatches() {
	ctx := context.Background()

	dispatches, err := s.service.GetAllDispatches(ctx)

	s.NoError(err)
	s.Equal(1, len(dispatches))
	s.Equal(dispatches[0].Id, services.USD_UAH_DISPATCH_ID)
}

func (s *DispatchServiceSuite) Test_SendDispatch() {
	ctx := context.Background()
	s.NoError(s.pgContainer.ExecuteSQLFiles(ctx, "add_couple_of_subscribers_for_usd_uah_dispatch"))

	expectedEmailReceivers := []string{"test_email_1@gmail.com", "test_email_2@gmail.com"}
	s.mailman.On("Send", mock.Anything).Return(nil)

	s.NoError(s.service.SendDispatch(ctx, services.USD_UAH_DISPATCH_ID))
	s.Equal(1, len(s.mailman.Calls))
	actualEmailReceivers := s.mailman.Calls[0].Arguments.Get(0).(mailing.Email).To
	s.Equal(actualEmailReceivers, expectedEmailReceivers)
}

func (s *DispatchServiceSuite) Test_SubscribeForDispatch_Success() {
	email := "email_1@gmail.com"
	dispatchId := services.USD_UAH_DISPATCH_ID
	ctx := context.Background()

	s.NoError(s.service.SubscribeForDispatch(ctx, email, dispatchId))
	subscribers, err := s.getSubscribersOfDispatch(ctx, dispatchId)

	s.NoError(err)
	s.True(slices.Contains(subscribers, email))
}

func (s *DispatchServiceSuite) Test_SubscribeForDispatch_UserAlreadySubscribedForThisDispatch() {
	email := "email_2@gmail.com"
	dispatchId := services.USD_UAH_DISPATCH_ID
	ctx := context.Background()

	s.NoError(s.service.SubscribeForDispatch(ctx, email, dispatchId))
	s.ErrorIs(
		s.service.SubscribeForDispatch(ctx, email, dispatchId),
		services.UniqueViolationErr,
	)
}

func TestIntegration_DispatchService(t *testing.T) {
	suite.Run(t, new(DispatchServiceSuite))
}
