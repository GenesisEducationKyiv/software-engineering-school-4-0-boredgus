package grpc

import (
	"context"
	"reflect"
	ds "subscription-api/internal/services/dispatch"
	pb_ds "subscription-api/pkg/grpc/dispatch_service"
	"testing"
)

func Test_dispatchServiceServer_SubscribeFor(t *testing.T) {
	type fields struct {
		s                                  ds.DispatchService
		UnimplementedDispatchServiceServer pb_ds.UnimplementedDispatchServiceServer
	}
	type args struct {
		ctx context.Context
		req *pb_ds.SubscribeForRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb_ds.SubscribeForResponse
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
			got, err := s.SubscribeFor(tt.args.ctx, tt.args.req)
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
