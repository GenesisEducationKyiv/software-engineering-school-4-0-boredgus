package grpc

import (
	"context"
	"reflect"
	"subscription-api/internal/services"
	pb_ds "subscription-api/pkg/grpc/dispatch_service"
	"testing"
)

func Test_dispatchServiceServer_GetDispatch(t *testing.T) {
	type fields struct {
		s                                  services.DispatchService
		UnimplementedDispatchServiceServer pb_ds.UnimplementedDispatchServiceServer
	}
	type args struct {
		ctx context.Context
		req *pb_ds.GetAllDispatchesRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb_ds.GetAllDispatchesResponse
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
