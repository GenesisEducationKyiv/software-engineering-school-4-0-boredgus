package dispatch_service_test

import (
	"context"
	"log"
	"net"
	"subscription-api/config"
	store "subscription-api/internal/db"
	"subscription-api/internal/services"
	"subscription-api/internal/sql"
	testhelpers "subscription-api/internal/tests/helpers"
	"subscription-api/internal/tests/stubs"
	"subscription-api/pkg/db"
	"subscription-api/pkg/utils"
	"testing"

	currency_service "subscription-api/internal/services/currency"
	currency_grpc "subscription-api/internal/services/currency/grpc"
	dispatch_service "subscription-api/internal/services/dispatch"

	grpc_client "subscription-api/pkg/grpc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DispatchServiceSuite struct {
	suite.Suite
	cs interface {
		Stop()
	}
	service     services.DispatchService
	l           config.Logger
	dbContainer *testhelpers.PostgresContainer
	ctx         context.Context
	mailman     *stubs.MailmanStub
}

const CurrencyServiceURL = "localhost:4040"

func (s *DispatchServiceSuite) startCurrencyServiceServer() {
	lis, err := net.Listen("tcp", CurrencyServiceURL)
	s.NoErrorf(err, "failed to listen url")

	server := grpc.NewServer()
	currency_grpc.RegisterCurrencyServiceServer(server,
		currency_service.NewCurrencyServiceServer(
			currency_service.NewCurrencyService(stubs.NewCurrencyAPIClient()),
			s.l,
		))
	s.cs = server

	go func() {
		if err := server.Serve(lis); err != nil {
			s.Fail(err.Error())
		}
	}()
}

func (s *DispatchServiceSuite) SetupSuite() {
	s.l = config.InitLogger(config.TestMode)
	s.startCurrencyServiceServer()

	s.ctx = context.Background()
	cont, err := testhelpers.CreatePostgresContainer(s.ctx, testhelpers.DBParams{
		Database: "subs",
		Username: "postgres",
		Password: "pass",
	})
	if err != nil {
		s.Failf("failed to create postgres container: %v", err.Error())
	}
	s.dbContainer = cont

	currencyServiceConn, err := grpc.NewClient(
		CurrencyServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.NoErrorf(err, "failed to create connection to currency service")

	mailman := stubs.NewMailmanStub()
	s.mailman = mailman

	s.service = dispatch_service.NewDispatchService(&dispatch_service.DispatchServiceParams{
		Logger: s.l,
		Store: store.NewStore(
			utils.Must(db.NewPostrgreSQL(
				s.dbContainer.ConnectionString,
				sql.PostgeSQLMigrationsUp(nil),
			)),
			db.IsPqError,
		),
		Mailman:         mailman,
		CurrencyService: grpc_client.NewCurrencyServiceClient(currencyServiceConn),
	})
}

func (s *DispatchServiceSuite) TearDownSuite() {
	s.cs.Stop()
	if err := s.dbContainer.Terminate(s.ctx); err != nil {
		log.Fatalf("error terminating database container: %s", err)
	}
}

func (s *DispatchServiceSuite) SetupTest() {}

func (s *DispatchServiceSuite) TestIntegration_GetAllDispatches() {
	ctx := context.Background()

	dispatches1, err1 := s.service.GetAllDispatches(ctx)
	dispatches2, err2 := s.service.GetAllDispatches(ctx)

	s.NoError(err1)
	s.NoError(err2)
	s.Equal(1, len(dispatches1))
	s.Equal(1, len(dispatches2))
	s.Equal(dispatches1, dispatches2)
	s.Equal(dispatches1[0].Id, services.USD_UAH_DISPATCH_ID)
}

func (s *DispatchServiceSuite) TestIntegration_SendDispatch() {
	dispatchId := services.USD_UAH_DISPATCH_ID
	emails := []string{"send.dispatch1@gmail.com", "send.dispatch2@gmail.com"}
	ctx := context.Background()

	for _, email := range emails {
		s.NoError(s.service.SubscribeForDispatch(ctx, email, dispatchId))
	}

	dispatchesAfter, err := s.service.GetAllDispatches(ctx)
	for _, d := range dispatchesAfter {
		if d.Id == dispatchId {
			assert.GreaterOrEqual(s.T(), d.CountOfSubscribers, len(emails))
		}
	}
	s.mailman.On("Send", mock.Anything).Once()

	s.NoError(err)
	s.NoError(s.service.SendDispatch(ctx, services.USD_UAH_DISPATCH_ID))
	s.Equal(1, len(s.mailman.Calls))
}

func (s *DispatchServiceSuite) TestIntegration_SubscribeForDispatch() {
	email1, email2 := "email1@gmail.com", "email2@gmail.com"
	dispatchId := services.USD_UAH_DISPATCH_ID
	countBefore, countAfter := 0, 0

	dispatchesBefore, err1 := s.service.GetAllDispatches(context.Background())
	err2 := s.service.SubscribeForDispatch(context.Background(), email1, dispatchId)
	err3 := s.service.SubscribeForDispatch(context.Background(), email1, dispatchId)
	err4 := s.service.SubscribeForDispatch(context.Background(), email1, "invalid-dispatch-uuid")
	err5 := s.service.SubscribeForDispatch(context.Background(), email1, "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	err6 := s.service.SubscribeForDispatch(context.Background(), email2, dispatchId)
	dispatchesAfter, err6 := s.service.GetAllDispatches(context.Background())
	for _, d := range dispatchesBefore {
		if d.Id == dispatchId {
			countBefore = d.CountOfSubscribers
		}
	}
	for _, d := range dispatchesAfter {
		if d.Id == dispatchId {
			countAfter = d.CountOfSubscribers
		}
	}

	s.NoError(err1)
	s.NoError(err2)
	s.ErrorIs(err3, services.UniqueViolationErr)
	s.ErrorIs(err4, services.InvalidArgumentErr)
	s.ErrorIs(err5, services.NotFoundErr)
	s.NoError(err6)
	s.Equal(2, countAfter-countBefore)
}

func TestIntegration_DispatchService(t *testing.T) {
	suite.Run(t, new(DispatchServiceSuite))
}
