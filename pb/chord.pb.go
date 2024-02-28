// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.3
// source: pb/chord.proto

package pb

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

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_chord_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_pb_chord_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_pb_chord_proto_rawDescGZIP(), []int{0}
}

// Note all of the response parameters are "optional"
// because during stabilization, a node will only report
// back whether its alive or not, but not actually respond
// to incoming queries.
type HashFuncResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HashFuncName *string `protobuf:"bytes,1,opt,name=HashFuncName,proto3,oneof" json:"HashFuncName,omitempty"`
}

func (x *HashFuncResponse) Reset() {
	*x = HashFuncResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_chord_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HashFuncResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HashFuncResponse) ProtoMessage() {}

func (x *HashFuncResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_chord_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HashFuncResponse.ProtoReflect.Descriptor instead.
func (*HashFuncResponse) Descriptor() ([]byte, []int) {
	return file_pb_chord_proto_rawDescGZIP(), []int{1}
}

func (x *HashFuncResponse) GetHashFuncName() string {
	if x != nil && x.HashFuncName != nil {
		return *x.HashFuncName
	}
	return ""
}

type HashKeyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HashValue string `protobuf:"bytes,1,opt,name=HashValue,proto3" json:"HashValue,omitempty"`
}

func (x *HashKeyRequest) Reset() {
	*x = HashKeyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_chord_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HashKeyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HashKeyRequest) ProtoMessage() {}

func (x *HashKeyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_chord_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HashKeyRequest.ProtoReflect.Descriptor instead.
func (*HashKeyRequest) Descriptor() ([]byte, []int) {
	return file_pb_chord_proto_rawDescGZIP(), []int{2}
}

func (x *HashKeyRequest) GetHashValue() string {
	if x != nil {
		return x.HashValue
	}
	return ""
}

type AddressResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address *string `protobuf:"bytes,1,opt,name=Address,proto3,oneof" json:"Address,omitempty"`
}

func (x *AddressResponse) Reset() {
	*x = AddressResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_chord_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddressResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddressResponse) ProtoMessage() {}

func (x *AddressResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_chord_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddressResponse.ProtoReflect.Descriptor instead.
func (*AddressResponse) Descriptor() ([]byte, []int) {
	return file_pb_chord_proto_rawDescGZIP(), []int{3}
}

func (x *AddressResponse) GetAddress() string {
	if x != nil && x.Address != nil {
		return *x.Address
	}
	return ""
}

type AddressRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,2,opt,name=Address,proto3" json:"Address,omitempty"`
}

func (x *AddressRequest) Reset() {
	*x = AddressRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_chord_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddressRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddressRequest) ProtoMessage() {}

func (x *AddressRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_chord_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddressRequest.ProtoReflect.Descriptor instead.
func (*AddressRequest) Descriptor() ([]byte, []int) {
	return file_pb_chord_proto_rawDescGZIP(), []int{4}
}

func (x *AddressRequest) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

var File_pb_chord_proto protoreflect.FileDescriptor

var file_pb_chord_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x70, 0x62, 0x2f, 0x63, 0x68, 0x6f, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x02, 0x70, 0x62, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x4c, 0x0a,
	0x10, 0x48, 0x61, 0x73, 0x68, 0x46, 0x75, 0x6e, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x27, 0x0a, 0x0c, 0x48, 0x61, 0x73, 0x68, 0x46, 0x75, 0x6e, 0x63, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0c, 0x48, 0x61, 0x73, 0x68, 0x46,
	0x75, 0x6e, 0x63, 0x4e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x48,
	0x61, 0x73, 0x68, 0x46, 0x75, 0x6e, 0x63, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x2e, 0x0a, 0x0e, 0x48,
	0x61, 0x73, 0x68, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a,
	0x09, 0x48, 0x61, 0x73, 0x68, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x48, 0x61, 0x73, 0x68, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x3c, 0x0a, 0x0f, 0x41,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d,
	0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x00, 0x52, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x88, 0x01, 0x01, 0x42, 0x0a, 0x0a,
	0x08, 0x5f, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x2a, 0x0a, 0x0e, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x41,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x32, 0xe1, 0x02, 0x0a, 0x05, 0x43, 0x68, 0x6f, 0x72, 0x64, 0x12,
	0x32, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x48, 0x61, 0x73, 0x68, 0x46, 0x75, 0x6e, 0x63, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x14, 0x2e,
	0x70, 0x62, 0x2e, 0x48, 0x61, 0x73, 0x68, 0x46, 0x75, 0x6e, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x6f, 0x72, 0x12, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x13,
	0x2e, 0x70, 0x62, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x16, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x73, 0x74, 0x50, 0x72,
	0x65, 0x63, 0x65, 0x64, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x12, 0x12, 0x2e,
	0x70, 0x62, 0x2e, 0x48, 0x61, 0x73, 0x68, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x13, 0x2e, 0x70, 0x62, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x0d, 0x46, 0x69, 0x6e, 0x64, 0x53, 0x75,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x12, 0x12, 0x2e, 0x70, 0x62, 0x2e, 0x48, 0x61, 0x73,
	0x68, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x70, 0x62,
	0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x30, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x50, 0x72, 0x65, 0x64, 0x65, 0x63, 0x65, 0x73, 0x73,
	0x6f, 0x72, 0x12, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x13, 0x2e,
	0x70, 0x62, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x27, 0x0a, 0x06, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12, 0x12, 0x2e, 0x70,
	0x62, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x1c, 0x0a, 0x04, 0x50,
	0x69, 0x6e, 0x67, 0x12, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x09,
	0x2e, 0x70, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x1b, 0x5a, 0x19, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x79, 0x75, 0x64, 0x73, 0x2f, 0x63, 0x68,
	0x6f, 0x72, 0x64, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_chord_proto_rawDescOnce sync.Once
	file_pb_chord_proto_rawDescData = file_pb_chord_proto_rawDesc
)

func file_pb_chord_proto_rawDescGZIP() []byte {
	file_pb_chord_proto_rawDescOnce.Do(func() {
		file_pb_chord_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_chord_proto_rawDescData)
	})
	return file_pb_chord_proto_rawDescData
}

var file_pb_chord_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_pb_chord_proto_goTypes = []interface{}{
	(*Empty)(nil),            // 0: pb.Empty
	(*HashFuncResponse)(nil), // 1: pb.HashFuncResponse
	(*HashKeyRequest)(nil),   // 2: pb.HashKeyRequest
	(*AddressResponse)(nil),  // 3: pb.AddressResponse
	(*AddressRequest)(nil),   // 4: pb.AddressRequest
}
var file_pb_chord_proto_depIdxs = []int32{
	0, // 0: pb.Chord.GetHashFuncName:input_type -> pb.Empty
	0, // 1: pb.Chord.GetSuccessor:input_type -> pb.Empty
	2, // 2: pb.Chord.ClosestPrecedingFinger:input_type -> pb.HashKeyRequest
	2, // 3: pb.Chord.FindSuccessor:input_type -> pb.HashKeyRequest
	0, // 4: pb.Chord.GetPredecessor:input_type -> pb.Empty
	4, // 5: pb.Chord.Notify:input_type -> pb.AddressRequest
	0, // 6: pb.Chord.Ping:input_type -> pb.Empty
	1, // 7: pb.Chord.GetHashFuncName:output_type -> pb.HashFuncResponse
	3, // 8: pb.Chord.GetSuccessor:output_type -> pb.AddressResponse
	3, // 9: pb.Chord.ClosestPrecedingFinger:output_type -> pb.AddressResponse
	3, // 10: pb.Chord.FindSuccessor:output_type -> pb.AddressResponse
	3, // 11: pb.Chord.GetPredecessor:output_type -> pb.AddressResponse
	0, // 12: pb.Chord.Notify:output_type -> pb.Empty
	0, // 13: pb.Chord.Ping:output_type -> pb.Empty
	7, // [7:14] is the sub-list for method output_type
	0, // [0:7] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pb_chord_proto_init() }
func file_pb_chord_proto_init() {
	if File_pb_chord_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_chord_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_pb_chord_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HashFuncResponse); i {
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
		file_pb_chord_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HashKeyRequest); i {
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
		file_pb_chord_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddressResponse); i {
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
		file_pb_chord_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddressRequest); i {
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
	file_pb_chord_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_pb_chord_proto_msgTypes[3].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pb_chord_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_chord_proto_goTypes,
		DependencyIndexes: file_pb_chord_proto_depIdxs,
		MessageInfos:      file_pb_chord_proto_msgTypes,
	}.Build()
	File_pb_chord_proto = out.File
	file_pb_chord_proto_rawDesc = nil
	file_pb_chord_proto_goTypes = nil
	file_pb_chord_proto_depIdxs = nil
}
