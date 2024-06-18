package dispatch_service

import (
	"context"
	"reflect"
	"subscription-api/internal/services"
	grpc "subscription-api/internal/services/dispatch/grpc"
	"testing"
)

func Test_dispatchServiceServer_GetDispatch(t *testing.T) {
	type fields struct {
		s                                  services.DispatchService
		UnimplementedDispatchServiceServer grpc.UnimplementedDispatchServiceServer
	}
	type args struct {
		ctx context.Context
		req *grpc.GetAllDispatchesRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *grpc.GetAllDispatchesResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dispatchServiceServer{
				s:                                  tt.fields.s,
				UnimplementedDispatchServiceServer: tt.fields.UnimplementedDispatchServiceServer,
			}
			got, err := s.GetAllDispatches(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("dispatchServiceServer.GetDispatch() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dispatchServiceServer.GetDispatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
