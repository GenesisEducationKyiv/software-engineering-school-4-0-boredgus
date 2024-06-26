// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.26.1
// source: internal/services/dispatch/grpc/dispatchService.proto

package __

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SubscribeForDispatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email      string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	DispatchId string `protobuf:"bytes,2,opt,name=dispatch_id,json=dispatchId,proto3" json:"dispatch_id,omitempty"`
}

func (x *SubscribeForDispatchRequest) Reset() {
	*x = SubscribeForDispatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeForDispatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeForDispatchRequest) ProtoMessage() {}

func (x *SubscribeForDispatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeForDispatchRequest.ProtoReflect.Descriptor instead.
func (*SubscribeForDispatchRequest) Descriptor() ([]byte, []int) {
	return file_internal_services_dispatch_grpc_dispatchService_proto_rawDescGZIP(), []int{0}
}

func (x *SubscribeForDispatchRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *SubscribeForDispatchRequest) GetDispatchId() string {
	if x != nil {
		return x.DispatchId
	}
	return ""
}

type SubscribeForDispatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SubscribeForDispatchResponse) Reset() {
	*x = SubscribeForDispatchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeForDispatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeForDispatchResponse) ProtoMessage() {}

func (x *SubscribeForDispatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeForDispatchResponse.ProtoReflect.Descriptor instead.
func (*SubscribeForDispatchResponse) Descriptor() ([]byte, []int) {
	return file_internal_services_dispatch_grpc_dispatchService_proto_rawDescGZIP(), []int{1}
}

type SendDispatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DispatchId string `protobuf:"bytes,1,opt,name=dispatch_id,json=dispatchId,proto3" json:"dispatch_id,omitempty"`
}

func (x *SendDispatchRequest) Reset() {
	*x = SendDispatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendDispatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendDispatchRequest) ProtoMessage() {}

func (x *SendDispatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendDispatchRequest.ProtoReflect.Descriptor instead.
func (*SendDispatchRequest) Descriptor() ([]byte, []int) {
	return file_internal_services_dispatch_grpc_dispatchService_proto_rawDescGZIP(), []int{2}
}

func (x *SendDispatchRequest) GetDispatchId() string {
	if x != nil {
		return x.DispatchId
	}
	return ""
}

type SendDispatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendDispatchResponse) Reset() {
	*x = SendDispatchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendDispatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendDispatchResponse) ProtoMessage() {}

func (x *SendDispatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendDispatchResponse.ProtoReflect.Descriptor instead.
func (*SendDispatchResponse) Descriptor() ([]byte, []int) {
	return file_internal_services_dispatch_grpc_dispatchService_proto_rawDescGZIP(), []int{3}
}

type GetAllDispatchesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetAllDispatchesRequest) Reset() {
	*x = GetAllDispatchesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllDispatchesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllDispatchesRequest) ProtoMessage() {}

func (x *GetAllDispatchesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllDispatchesRequest.ProtoReflect.Descriptor instead.
func (*GetAllDispatchesRequest) Descriptor() ([]byte, []int) {
	return file_internal_services_dispatch_grpc_dispatchService_proto_rawDescGZIP(), []int{4}
}

type DispatchData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// string baseCurrency = 2;
	// repeated string targetCurrencies = 3;
	SendAt             string `protobuf:"bytes,4,opt,name=send_at,json=sendAt,proto3" json:"send_at,omitempty"`
	CountOfSubscribers int64  `protobuf:"varint,5,opt,name=countOfSubscribers,proto3" json:"countOfSubscribers,omitempty"`
	Label              string `protobuf:"bytes,6,opt,name=label,proto3" json:"label,omitempty"`
}

func (x *DispatchData) Reset() {
	*x = DispatchData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DispatchData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DispatchData) ProtoMessage() {}

func (x *DispatchData) ProtoReflect() protoreflect.Message {
	mi := &file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DispatchData.ProtoReflect.Descriptor instead.
func (*DispatchData) Descriptor() ([]byte, []int) {
	return file_internal_services_dispatch_grpc_dispatchService_proto_rawDescGZIP(), []int{5}
}

func (x *DispatchData) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *DispatchData) GetSendAt() string {
	if x != nil {
		return x.SendAt
	}
	return ""
}

func (x *DispatchData) GetCountOfSubscribers() int64 {
	if x != nil {
		return x.CountOfSubscribers
	}
	return 0
}

func (x *DispatchData) GetLabel() string {
	if x != nil {
		return x.Label
	}
	return ""
}

type GetAllDispatchesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Dispatches []*DispatchData `protobuf:"bytes,1,rep,name=dispatches,proto3" json:"dispatches,omitempty"`
}

func (x *GetAllDispatchesResponse) Reset() {
	*x = GetAllDispatchesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllDispatchesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllDispatchesResponse) ProtoMessage() {}

func (x *GetAllDispatchesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllDispatchesResponse.ProtoReflect.Descriptor instead.
func (*GetAllDispatchesResponse) Descriptor() ([]byte, []int) {
	return file_internal_services_dispatch_grpc_dispatchService_proto_rawDescGZIP(), []int{6}
}

func (x *GetAllDispatchesResponse) GetDispatches() []*DispatchData {
	if x != nil {
		return x.Dispatches
	}
	return nil
}

var File_internal_services_dispatch_grpc_dispatchService_proto protoreflect.FileDescriptor

var file_internal_services_dispatch_grpc_dispatchService_proto_rawDesc = []byte{
	0x0a, 0x35, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x2f, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x2f, 0x67, 0x72, 0x70,
	0x63, 0x2f, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x6d, 0x61, 0x69, 0x6e, 0x22, 0x54, 0x0a,
	0x1b, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x46, 0x6f, 0x72, 0x44, 0x69, 0x73,
	0x70, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05,
	0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61,
	0x69, 0x6c, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63,
	0x68, 0x49, 0x64, 0x22, 0x1e, 0x0a, 0x1c, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65,
	0x46, 0x6f, 0x72, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x36, 0x0a, 0x13, 0x53, 0x65, 0x6e, 0x64, 0x44, 0x69, 0x73, 0x70, 0x61,
	0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x69,
	0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x49, 0x64, 0x22, 0x16, 0x0a, 0x14, 0x53,
	0x65, 0x6e, 0x64, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x19, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x44, 0x69, 0x73,
	0x70, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x7d,
	0x0a, 0x0c, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x44, 0x61, 0x74, 0x61, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x17,
	0x0a, 0x07, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x65, 0x6e, 0x64, 0x41, 0x74, 0x12, 0x2e, 0x0a, 0x12, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x4f, 0x66, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x73, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x12, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x4f, 0x66, 0x53, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x22, 0x4e, 0x0a,
	0x18, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x65,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x0a, 0x64, 0x69, 0x73,
	0x70, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x0a, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x32, 0x90, 0x02,
	0x0a, 0x0f, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x5f, 0x0a, 0x14, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x46, 0x6f,
	0x72, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x12, 0x21, 0x2e, 0x6d, 0x61, 0x69, 0x6e,
	0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x46, 0x6f, 0x72, 0x44, 0x69, 0x73,
	0x70, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x6d,
	0x61, 0x69, 0x6e, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x46, 0x6f, 0x72,
	0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x47, 0x0a, 0x0c, 0x53, 0x65, 0x6e, 0x64, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74,
	0x63, 0x68, 0x12, 0x19, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x44, 0x69,
	0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e,
	0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63,
	0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x53, 0x0a, 0x10, 0x47,
	0x65, 0x74, 0x41, 0x6c, 0x6c, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x12,
	0x1d, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x44, 0x69, 0x73,
	0x70, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e,
	0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x44, 0x69, 0x73, 0x70,
	0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_services_dispatch_grpc_dispatchService_proto_rawDescOnce sync.Once
	file_internal_services_dispatch_grpc_dispatchService_proto_rawDescData = file_internal_services_dispatch_grpc_dispatchService_proto_rawDesc
)

func file_internal_services_dispatch_grpc_dispatchService_proto_rawDescGZIP() []byte {
	file_internal_services_dispatch_grpc_dispatchService_proto_rawDescOnce.Do(func() {
		file_internal_services_dispatch_grpc_dispatchService_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_services_dispatch_grpc_dispatchService_proto_rawDescData)
	})
	return file_internal_services_dispatch_grpc_dispatchService_proto_rawDescData
}

var file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_internal_services_dispatch_grpc_dispatchService_proto_goTypes = []interface{}{
	(*SubscribeForDispatchRequest)(nil),  // 0: main.SubscribeForDispatchRequest
	(*SubscribeForDispatchResponse)(nil), // 1: main.SubscribeForDispatchResponse
	(*SendDispatchRequest)(nil),          // 2: main.SendDispatchRequest
	(*SendDispatchResponse)(nil),         // 3: main.SendDispatchResponse
	(*GetAllDispatchesRequest)(nil),      // 4: main.GetAllDispatchesRequest
	(*DispatchData)(nil),                 // 5: main.DispatchData
	(*GetAllDispatchesResponse)(nil),     // 6: main.GetAllDispatchesResponse
}
var file_internal_services_dispatch_grpc_dispatchService_proto_depIdxs = []int32{
	5, // 0: main.GetAllDispatchesResponse.dispatches:type_name -> main.DispatchData
	0, // 1: main.DispatchService.SubscribeForDispatch:input_type -> main.SubscribeForDispatchRequest
	2, // 2: main.DispatchService.SendDispatch:input_type -> main.SendDispatchRequest
	4, // 3: main.DispatchService.GetAllDispatches:input_type -> main.GetAllDispatchesRequest
	1, // 4: main.DispatchService.SubscribeForDispatch:output_type -> main.SubscribeForDispatchResponse
	3, // 5: main.DispatchService.SendDispatch:output_type -> main.SendDispatchResponse
	6, // 6: main.DispatchService.GetAllDispatches:output_type -> main.GetAllDispatchesResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_internal_services_dispatch_grpc_dispatchService_proto_init() }
func file_internal_services_dispatch_grpc_dispatchService_proto_init() {
	if File_internal_services_dispatch_grpc_dispatchService_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscribeForDispatchRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscribeForDispatchResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendDispatchRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendDispatchResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllDispatchesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DispatchData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllDispatchesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_services_dispatch_grpc_dispatchService_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_services_dispatch_grpc_dispatchService_proto_goTypes,
		DependencyIndexes: file_internal_services_dispatch_grpc_dispatchService_proto_depIdxs,
		MessageInfos:      file_internal_services_dispatch_grpc_dispatchService_proto_msgTypes,
	}.Build()
	File_internal_services_dispatch_grpc_dispatchService_proto = out.File
	file_internal_services_dispatch_grpc_dispatchService_proto_rawDesc = nil
	file_internal_services_dispatch_grpc_dispatchService_proto_goTypes = nil
	file_internal_services_dispatch_grpc_dispatchService_proto_depIdxs = nil
}
