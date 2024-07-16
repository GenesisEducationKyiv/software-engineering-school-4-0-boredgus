package server

import (
	"context"
	"reflect"
	"testing"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/config"
	grpc_gen "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/grpc/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Test_customerServiceServer_CreateCustomer(t *testing.T) {
	type fields struct {
		UnimplementedCustomerServiceServer grpc_gen.UnimplementedCustomerServiceServer
		service                            CustomerService
		logger                             config.Logger
	}
	type args struct {
		ctx context.Context
		req *grpc_gen.CreateCustomerRequest
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		expectedResponse *emptypb.Empty
		expectedErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &customerServiceServer{
				UnimplementedCustomerServiceServer: tt.fields.UnimplementedCustomerServiceServer,
				service:                            tt.fields.service,
				logger:                             tt.fields.logger,
			}
			got, err := s.CreateCustomer(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.expectedErr {
				t.Errorf("customerServiceServer.CreateCustomer() error = %v, wantErr %v", err, tt.expectedErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expectedResponse) {
				t.Errorf("customerServiceServer.CreateCustomer() = %v, want %v", got, tt.expectedResponse)
			}
		})
	}
}
