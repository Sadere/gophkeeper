// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: proto/keeper/v1/keeper.proto

package keeperv1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
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

type LoginV1Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Login    string `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *LoginV1Request) Reset() {
	*x = LoginV1Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_keeper_v1_keeper_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginV1Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginV1Request) ProtoMessage() {}

func (x *LoginV1Request) ProtoReflect() protoreflect.Message {
	mi := &file_proto_keeper_v1_keeper_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginV1Request.ProtoReflect.Descriptor instead.
func (*LoginV1Request) Descriptor() ([]byte, []int) {
	return file_proto_keeper_v1_keeper_proto_rawDescGZIP(), []int{0}
}

func (x *LoginV1Request) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *LoginV1Request) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type LoginV1Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccessToken string `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
}

func (x *LoginV1Response) Reset() {
	*x = LoginV1Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_keeper_v1_keeper_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginV1Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginV1Response) ProtoMessage() {}

func (x *LoginV1Response) ProtoReflect() protoreflect.Message {
	mi := &file_proto_keeper_v1_keeper_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginV1Response.ProtoReflect.Descriptor instead.
func (*LoginV1Response) Descriptor() ([]byte, []int) {
	return file_proto_keeper_v1_keeper_proto_rawDescGZIP(), []int{1}
}

func (x *LoginV1Response) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

type RegisterV1Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Login    string `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *RegisterV1Request) Reset() {
	*x = RegisterV1Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_keeper_v1_keeper_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterV1Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterV1Request) ProtoMessage() {}

func (x *RegisterV1Request) ProtoReflect() protoreflect.Message {
	mi := &file_proto_keeper_v1_keeper_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterV1Request.ProtoReflect.Descriptor instead.
func (*RegisterV1Request) Descriptor() ([]byte, []int) {
	return file_proto_keeper_v1_keeper_proto_rawDescGZIP(), []int{2}
}

func (x *RegisterV1Request) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *RegisterV1Request) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type RegisterV1Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccessToken string `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
}

func (x *RegisterV1Response) Reset() {
	*x = RegisterV1Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_keeper_v1_keeper_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterV1Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterV1Response) ProtoMessage() {}

func (x *RegisterV1Response) ProtoReflect() protoreflect.Message {
	mi := &file_proto_keeper_v1_keeper_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterV1Response.ProtoReflect.Descriptor instead.
func (*RegisterV1Response) Descriptor() ([]byte, []int) {
	return file_proto_keeper_v1_keeper_proto_rawDescGZIP(), []int{3}
}

func (x *RegisterV1Response) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

var File_proto_keeper_v1_keeper_proto protoreflect.FileDescriptor

var file_proto_keeper_v1_keeper_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2f, 0x76,
	0x31, 0x2f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x1a,
	0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x58, 0x0a, 0x0e,
	0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f,
	0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xba,
	0x48, 0x06, 0x72, 0x04, 0x10, 0x05, 0x18, 0x64, 0x52, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12,
	0x25, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x09, 0xba, 0x48, 0x06, 0x72, 0x04, 0x10, 0x06, 0x18, 0x14, 0x52, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x34, 0x0a, 0x0f, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x56,
	0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x5b, 0x0a, 0x11,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1f, 0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x09, 0xba, 0x48, 0x06, 0x72, 0x04, 0x10, 0x05, 0x18, 0x64, 0x52, 0x05, 0x6c, 0x6f, 0x67,
	0x69, 0x6e, 0x12, 0x25, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xba, 0x48, 0x06, 0x72, 0x04, 0x10, 0x06, 0x18, 0x14, 0x52,
	0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x37, 0x0a, 0x12, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x32, 0xb2, 0x01, 0x0a, 0x0b, 0x41, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x4c, 0x0a, 0x07, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x56, 0x31, 0x12, 0x1f, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e,
	0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x76, 0x31,
	0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x55, 0x0a, 0x0a, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x56, 0x31, 0x12, 0x22,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x76, 0x31,
	0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6b, 0x65, 0x65, 0x70, 0x65,
	0x72, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x56, 0x31, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x9a, 0x01, 0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x42,
	0x0b, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x18,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x3b,
	0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x50, 0x4b, 0x58, 0xaa, 0x02,
	0x0f, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x56, 0x31,
	0xca, 0x02, 0x0f, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x5c,
	0x56, 0x31, 0xe2, 0x02, 0x1b, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x4b, 0x65, 0x65, 0x70, 0x65,
	0x72, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x11, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x3a, 0x3a, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72,
	0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_keeper_v1_keeper_proto_rawDescOnce sync.Once
	file_proto_keeper_v1_keeper_proto_rawDescData = file_proto_keeper_v1_keeper_proto_rawDesc
)

func file_proto_keeper_v1_keeper_proto_rawDescGZIP() []byte {
	file_proto_keeper_v1_keeper_proto_rawDescOnce.Do(func() {
		file_proto_keeper_v1_keeper_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_keeper_v1_keeper_proto_rawDescData)
	})
	return file_proto_keeper_v1_keeper_proto_rawDescData
}

var file_proto_keeper_v1_keeper_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_keeper_v1_keeper_proto_goTypes = []any{
	(*LoginV1Request)(nil),     // 0: proto.keeper.v1.LoginV1Request
	(*LoginV1Response)(nil),    // 1: proto.keeper.v1.LoginV1Response
	(*RegisterV1Request)(nil),  // 2: proto.keeper.v1.RegisterV1Request
	(*RegisterV1Response)(nil), // 3: proto.keeper.v1.RegisterV1Response
}
var file_proto_keeper_v1_keeper_proto_depIdxs = []int32{
	0, // 0: proto.keeper.v1.AuthService.LoginV1:input_type -> proto.keeper.v1.LoginV1Request
	2, // 1: proto.keeper.v1.AuthService.RegisterV1:input_type -> proto.keeper.v1.RegisterV1Request
	1, // 2: proto.keeper.v1.AuthService.LoginV1:output_type -> proto.keeper.v1.LoginV1Response
	3, // 3: proto.keeper.v1.AuthService.RegisterV1:output_type -> proto.keeper.v1.RegisterV1Response
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_keeper_v1_keeper_proto_init() }
func file_proto_keeper_v1_keeper_proto_init() {
	if File_proto_keeper_v1_keeper_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_keeper_v1_keeper_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*LoginV1Request); i {
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
		file_proto_keeper_v1_keeper_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*LoginV1Response); i {
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
		file_proto_keeper_v1_keeper_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*RegisterV1Request); i {
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
		file_proto_keeper_v1_keeper_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*RegisterV1Response); i {
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
			RawDescriptor: file_proto_keeper_v1_keeper_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_keeper_v1_keeper_proto_goTypes,
		DependencyIndexes: file_proto_keeper_v1_keeper_proto_depIdxs,
		MessageInfos:      file_proto_keeper_v1_keeper_proto_msgTypes,
	}.Build()
	File_proto_keeper_v1_keeper_proto = out.File
	file_proto_keeper_v1_keeper_proto_rawDesc = nil
	file_proto_keeper_v1_keeper_proto_goTypes = nil
	file_proto_keeper_v1_keeper_proto_depIdxs = nil
}
