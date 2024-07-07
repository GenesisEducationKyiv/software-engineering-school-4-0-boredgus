// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.26.1
// source: messages/subscription/subscription_messages.proto

package messages

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SubscriptionStatus int32

const (
	SubscriptionStatus_CREATED SubscriptionStatus = 0
)

// Enum value maps for SubscriptionStatus.
var (
	SubscriptionStatus_name = map[int32]string{
		0: "CREATED",
	}
	SubscriptionStatus_value = map[string]int32{
		"CREATED": 0,
	}
)

func (x SubscriptionStatus) Enum() *SubscriptionStatus {
	p := new(SubscriptionStatus)
	*p = x
	return p
}

func (x SubscriptionStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SubscriptionStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_messages_subscription_subscription_messages_proto_enumTypes[0].Descriptor()
}

func (SubscriptionStatus) Type() protoreflect.EnumType {
	return &file_messages_subscription_subscription_messages_proto_enumTypes[0]
}

func (x SubscriptionStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SubscriptionStatus.Descriptor instead.
func (SubscriptionStatus) EnumDescriptor() ([]byte, []int) {
	return file_messages_subscription_subscription_messages_proto_rawDescGZIP(), []int{0}
}

type Subscription struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DispatchID  string                 `protobuf:"bytes,1,opt,name=dispatchID,proto3" json:"dispatchID,omitempty"`
	BaseCcy     string                 `protobuf:"bytes,2,opt,name=baseCcy,proto3" json:"baseCcy,omitempty"`
	TargetCcies []string               `protobuf:"bytes,3,rep,name=targetCcies,proto3" json:"targetCcies,omitempty"`
	Status      SubscriptionStatus     `protobuf:"varint,4,opt,name=status,proto3,enum=messages.subscription.SubscriptionStatus" json:"status,omitempty"`
	Email       string                 `protobuf:"bytes,5,opt,name=email,proto3" json:"email,omitempty"`
	SendAt      *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=sendAt,proto3" json:"sendAt,omitempty"`
}

func (x *Subscription) Reset() {
	*x = Subscription{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_subscription_subscription_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Subscription) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Subscription) ProtoMessage() {}

func (x *Subscription) ProtoReflect() protoreflect.Message {
	mi := &file_messages_subscription_subscription_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Subscription.ProtoReflect.Descriptor instead.
func (*Subscription) Descriptor() ([]byte, []int) {
	return file_messages_subscription_subscription_messages_proto_rawDescGZIP(), []int{0}
}

func (x *Subscription) GetDispatchID() string {
	if x != nil {
		return x.DispatchID
	}
	return ""
}

func (x *Subscription) GetBaseCcy() string {
	if x != nil {
		return x.BaseCcy
	}
	return ""
}

func (x *Subscription) GetTargetCcies() []string {
	if x != nil {
		return x.TargetCcies
	}
	return nil
}

func (x *Subscription) GetStatus() SubscriptionStatus {
	if x != nil {
		return x.Status
	}
	return SubscriptionStatus_CREATED
}

func (x *Subscription) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *Subscription) GetSendAt() *timestamppb.Timestamp {
	if x != nil {
		return x.SendAt
	}
	return nil
}

type SubscriptionMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventType EventType              `protobuf:"varint,1,opt,name=eventType,proto3,enum=messages.EventType" json:"eventType,omitempty"`
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Payload   *Subscription          `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *SubscriptionMessage) Reset() {
	*x = SubscriptionMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_subscription_subscription_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscriptionMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscriptionMessage) ProtoMessage() {}

func (x *SubscriptionMessage) ProtoReflect() protoreflect.Message {
	mi := &file_messages_subscription_subscription_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscriptionMessage.ProtoReflect.Descriptor instead.
func (*SubscriptionMessage) Descriptor() ([]byte, []int) {
	return file_messages_subscription_subscription_messages_proto_rawDescGZIP(), []int{1}
}

func (x *SubscriptionMessage) GetEventType() EventType {
	if x != nil {
		return x.EventType
	}
	return EventType_SUBSCRIPTION_CREATED
}

func (x *SubscriptionMessage) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *SubscriptionMessage) GetPayload() *Subscription {
	if x != nil {
		return x.Payload
	}
	return nil
}

var File_messages_subscription_subscription_messages_proto protoreflect.FileDescriptor

var file_messages_subscription_subscription_messages_proto_rawDesc = []byte{
	0x0a, 0x31, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2f, 0x73, 0x75, 0x62, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x15, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x73, 0x75,
	0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x73, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf7, 0x01, 0x0a, 0x0c, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63,
	0x68, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x69, 0x73, 0x70, 0x61,
	0x74, 0x63, 0x68, 0x49, 0x44, 0x12, 0x18, 0x0a, 0x07, 0x62, 0x61, 0x73, 0x65, 0x43, 0x63, 0x79,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x62, 0x61, 0x73, 0x65, 0x43, 0x63, 0x79, 0x12,
	0x20, 0x0a, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x43, 0x63, 0x69, 0x65, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x43, 0x63, 0x69, 0x65,
	0x73, 0x12, 0x41, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x29, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x73, 0x75, 0x62,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x32, 0x0a, 0x06, 0x73, 0x65,
	0x6e, 0x64, 0x41, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x41, 0x74, 0x22, 0xc1,
	0x01, 0x0a, 0x13, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x31, 0x0a, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x73, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x12, 0x3d, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e,
	0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x75, 0x62,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x2a, 0x21, 0x0a, 0x12, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x52, 0x45, 0x41,
	0x54, 0x45, 0x44, 0x10, 0x00, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x3b, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_messages_subscription_subscription_messages_proto_rawDescOnce sync.Once
	file_messages_subscription_subscription_messages_proto_rawDescData = file_messages_subscription_subscription_messages_proto_rawDesc
)

func file_messages_subscription_subscription_messages_proto_rawDescGZIP() []byte {
	file_messages_subscription_subscription_messages_proto_rawDescOnce.Do(func() {
		file_messages_subscription_subscription_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_messages_subscription_subscription_messages_proto_rawDescData)
	})
	return file_messages_subscription_subscription_messages_proto_rawDescData
}

var file_messages_subscription_subscription_messages_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_messages_subscription_subscription_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_messages_subscription_subscription_messages_proto_goTypes = []interface{}{
	(SubscriptionStatus)(0),       // 0: messages.subscription.SubscriptionStatus
	(*Subscription)(nil),          // 1: messages.subscription.Subscription
	(*SubscriptionMessage)(nil),   // 2: messages.subscription.SubscriptionMessage
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
	(EventType)(0),                // 4: messages.EventType
}
var file_messages_subscription_subscription_messages_proto_depIdxs = []int32{
	0, // 0: messages.subscription.Subscription.status:type_name -> messages.subscription.SubscriptionStatus
	3, // 1: messages.subscription.Subscription.sendAt:type_name -> google.protobuf.Timestamp
	4, // 2: messages.subscription.SubscriptionMessage.eventType:type_name -> messages.EventType
	3, // 3: messages.subscription.SubscriptionMessage.timestamp:type_name -> google.protobuf.Timestamp
	1, // 4: messages.subscription.SubscriptionMessage.payload:type_name -> messages.subscription.Subscription
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_messages_subscription_subscription_messages_proto_init() }
func file_messages_subscription_subscription_messages_proto_init() {
	if File_messages_subscription_subscription_messages_proto != nil {
		return
	}
	file_messages_messages_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_messages_subscription_subscription_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Subscription); i {
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
		file_messages_subscription_subscription_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscriptionMessage); i {
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
			RawDescriptor: file_messages_subscription_subscription_messages_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_messages_subscription_subscription_messages_proto_goTypes,
		DependencyIndexes: file_messages_subscription_subscription_messages_proto_depIdxs,
		EnumInfos:         file_messages_subscription_subscription_messages_proto_enumTypes,
		MessageInfos:      file_messages_subscription_subscription_messages_proto_msgTypes,
	}.Build()
	File_messages_subscription_subscription_messages_proto = out.File
	file_messages_subscription_subscription_messages_proto_rawDesc = nil
	file_messages_subscription_subscription_messages_proto_goTypes = nil
	file_messages_subscription_subscription_messages_proto_depIdxs = nil
}