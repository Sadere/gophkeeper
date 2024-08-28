// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: proto/keeper/v1/notification.proto

package keeperv1

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SubscribeV1Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *SubscribeV1Request) Reset() {
	*x = SubscribeV1Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_keeper_v1_notification_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeV1Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeV1Request) ProtoMessage() {}

func (x *SubscribeV1Request) ProtoReflect() protoreflect.Message {
	mi := &file_proto_keeper_v1_notification_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeV1Request.ProtoReflect.Descriptor instead.
func (*SubscribeV1Request) Descriptor() ([]byte, []int) {
	return file_proto_keeper_v1_notification_proto_rawDescGZIP(), []int{0}
}

func (x *SubscribeV1Request) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type SubscribeV1Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Updated bool   `protobuf:"varint,2,opt,name=updated,proto3" json:"updated,omitempty"`
}

func (x *SubscribeV1Response) Reset() {
	*x = SubscribeV1Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_keeper_v1_notification_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeV1Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeV1Response) ProtoMessage() {}

func (x *SubscribeV1Response) ProtoReflect() protoreflect.Message {
	mi := &file_proto_keeper_v1_notification_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeV1Response.ProtoReflect.Descriptor instead.
func (*SubscribeV1Response) Descriptor() ([]byte, []int) {
	return file_proto_keeper_v1_notification_proto_rawDescGZIP(), []int{1}
}

func (x *SubscribeV1Response) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *SubscribeV1Response) GetUpdated() bool {
	if x != nil {
		return x.Updated
	}
	return false
}

var File_proto_keeper_v1_notification_proto protoreflect.FileDescriptor

var file_proto_keeper_v1_notification_proto_rawDesc = []byte{
	0x0a, 0x22, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2f, 0x76,
	0x31, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6b, 0x65, 0x65, 0x70,
	0x65, 0x72, 0x2e, 0x76, 0x31, 0x22, 0x24, 0x0a, 0x12, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69,
	0x62, 0x65, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x22, 0x3f, 0x0a, 0x13, 0x53,
	0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x07, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x32, 0x71, 0x0a, 0x13,
	0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x5a, 0x0a, 0x0b, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65,
	0x56, 0x31, 0x12, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6b, 0x65, 0x65, 0x70, 0x65,
	0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x56, 0x31,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72,
	0x69, 0x62, 0x65, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42,
	0xa0, 0x01, 0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6b, 0x65,
	0x65, 0x70, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x42, 0x11, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x18, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x3b, 0x6b, 0x65,
	0x65, 0x70, 0x65, 0x72, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x50, 0x4b, 0x58, 0xaa, 0x02, 0x0f, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x56, 0x31, 0xca, 0x02,
	0x0f, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x5c, 0x56, 0x31,
	0xe2, 0x02, 0x1b, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x5c,
	0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x11, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x3a, 0x3a, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x3a, 0x3a,
	0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_keeper_v1_notification_proto_rawDescOnce sync.Once
	file_proto_keeper_v1_notification_proto_rawDescData = file_proto_keeper_v1_notification_proto_rawDesc
)

func file_proto_keeper_v1_notification_proto_rawDescGZIP() []byte {
	file_proto_keeper_v1_notification_proto_rawDescOnce.Do(func() {
		file_proto_keeper_v1_notification_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_keeper_v1_notification_proto_rawDescData)
	})
	return file_proto_keeper_v1_notification_proto_rawDescData
}

var file_proto_keeper_v1_notification_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_keeper_v1_notification_proto_goTypes = []any{
	(*SubscribeV1Request)(nil),  // 0: proto.keeper.v1.SubscribeV1Request
	(*SubscribeV1Response)(nil), // 1: proto.keeper.v1.SubscribeV1Response
}
var file_proto_keeper_v1_notification_proto_depIdxs = []int32{
	0, // 0: proto.keeper.v1.NotificationService.SubscribeV1:input_type -> proto.keeper.v1.SubscribeV1Request
	1, // 1: proto.keeper.v1.NotificationService.SubscribeV1:output_type -> proto.keeper.v1.SubscribeV1Response
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_keeper_v1_notification_proto_init() }
func file_proto_keeper_v1_notification_proto_init() {
	if File_proto_keeper_v1_notification_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_keeper_v1_notification_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*SubscribeV1Request); i {
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
		file_proto_keeper_v1_notification_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*SubscribeV1Response); i {
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
			RawDescriptor: file_proto_keeper_v1_notification_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_keeper_v1_notification_proto_goTypes,
		DependencyIndexes: file_proto_keeper_v1_notification_proto_depIdxs,
		MessageInfos:      file_proto_keeper_v1_notification_proto_msgTypes,
	}.Build()
	File_proto_keeper_v1_notification_proto = out.File
	file_proto_keeper_v1_notification_proto_rawDesc = nil
	file_proto_keeper_v1_notification_proto_goTypes = nil
	file_proto_keeper_v1_notification_proto_depIdxs = nil
}
