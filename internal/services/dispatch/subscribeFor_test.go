package dispatch_service

import (
	"context"
	"reflect"
	"subscription-api/internal/services"
	dispatch_grpc "subscription-api/internal/services/dispatch/grpc"
	"testing"
)

func Test_DispatchServiceServer_SubscribeForDispatch(t *testing.T) {
	type fields struct {
		s                                  services.DispatchService
		UnimplementedDispatchServiceServer dispatch_grpc.UnimplementedDispatchServiceServer
	}
	type args struct {
		ctx context.Context
		req *dispatch_grpc.SubscribeForDispatchRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dispatch_grpc.SubscribeForDispatchRequest
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
			got, err := s.SubscribeForDispatch(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("dispatchServiceServer.SubscribeFor() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dispatchServiceServer.SubscribeFor() = %v, want %v", got, tt.want)
			}
		})
	}
}
