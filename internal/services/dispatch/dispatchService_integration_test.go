package ds

import (
	"context"
	"log"
	"net"
	"subscription-api/config"
	store "subscription-api/internal/db"
	"subscription-api/internal/services"
	"subscription-api/internal/sql"
	"subscription-api/internal/testhelpers"
	"subscription-api/internal/testhelpers/stubs"
	"subscription-api/pkg/db"
	pb_cs "subscription-api/pkg/grpc/currency_service"
	"subscription-api/pkg/utils"
	"testing"

	cs "subscription-api/internal/services/currency"
	g "subscription-api/internal/services/currency/grpc"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DispatchServiceSuite struct {
	suite.Suite
	cs interface {
		Stop()
	}
	service     *dispatchService
	l           config.Logger
	dbContainer *testhelpers.PostgresContainer
	ctx         context.Context
}

const CurrencyServiceURL = "localhost:4040"

func (s *DispatchServiceSuite) startCurrencyServiceServer() {
	lis, err := net.Listen("tcp", CurrencyServiceURL)
	s.NoErrorf(err, "failed to listen url")

	server := grpc.NewServer()
	pb_cs.RegisterCurrencyServiceServer(server,
		g.NewCurrencyServiceServer(
			cs.NewCurrencyService(stubs.NewCurrencyAPIClient()),
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

	s.l = config.InitLogger(config.TestMode)

	currencyServiceConn, err := grpc.NewClient(
		CurrencyServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.NoErrorf(err, "failed to create connection to currency service")

	s.service = NewDispatchService(&DispatchServiceParams{
		Logger: s.l,
		Store: store.NewStore(
			utils.Must(db.NewPostrgreSQL(
				s.dbContainer.ConnectionString,
				sql.PostgeSQLMigrationsUp(nil),
			)),
			db.IsPqError,
		),
		Mailman:         stubs.NewMailmanStub(),
		CurrencyService: pb_cs.NewCurrencyServiceClient(currencyServiceConn),
	})
	// s.NoError(s.dbContainer.Snapshot(s.ctx, postgres.WithSnapshotName("initial")))
}

func (s *DispatchServiceSuite) TearDownSuite() {
	s.cs.Stop()
	if err := s.dbContainer.Terminate(s.ctx); err != nil {
		log.Fatalf("error terminating database container: %s", err)
	}
}

func (s *DispatchServiceSuite) SetupTest() {
	// s.NoError(s.dbContainer.Restore(s.ctx, postgres.WithSnapshotName("initial")))
}

func (s *DispatchServiceSuite) Test_GetAllDispatches() {
	dispatches, err := s.service.GetAllDispatches(context.Background())
	s.NoError(err)
	s.Equal(1, len(dispatches))
	d := dispatches[0]
	s.Equal(d.Id, services.USD_UAH_DISPATCH_ID)
}

func (s *DispatchServiceSuite) Test_SubscribeForDispatch() {
	email1, email2 := "email@gmail.com", "qwerty@gmail.com"

	err := s.service.SubscribeForDispatch(context.Background(), email1, services.USD_UAH_DISPATCH_ID)
	s.NoError(err)
	err = s.service.SubscribeForDispatch(context.Background(), email1, services.USD_UAH_DISPATCH_ID)
	s.ErrorIs(err, services.UniqueViolationErr)
	err = s.service.SubscribeForDispatch(context.Background(), email2, services.USD_UAH_DISPATCH_ID)
	s.NoError(err)
	dispatches, err := s.service.GetAllDispatches(context.Background())
	s.NoError(err)
	for _, d := range dispatches {
		if d.Id == services.USD_UAH_DISPATCH_ID {
			s.Equal(2, d.CountOfSubscribers)
		}
	}
}

func (s *DispatchServiceSuite) Test_SendDispatch() {
	s.NoError(s.service.SendDispatch(context.Background(), services.USD_UAH_DISPATCH_ID))
}

func TestDispatchServiceSuite(t *testing.T) {
	suite.Run(t, new(DispatchServiceSuite))
}
