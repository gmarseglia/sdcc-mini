// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.12
// source: cs-mw.proto

package proto

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

type ChoiceBiRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Option1 string `protobuf:"bytes,1,opt,name=Option1,proto3" json:"Option1,omitempty"`
	Option2 string `protobuf:"bytes,2,opt,name=Option2,proto3" json:"Option2,omitempty"`
}

func (x *ChoiceBiRequest) Reset() {
	*x = ChoiceBiRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cs_mw_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChoiceBiRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChoiceBiRequest) ProtoMessage() {}

func (x *ChoiceBiRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cs_mw_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChoiceBiRequest.ProtoReflect.Descriptor instead.
func (*ChoiceBiRequest) Descriptor() ([]byte, []int) {
	return file_cs_mw_proto_rawDescGZIP(), []int{0}
}

func (x *ChoiceBiRequest) GetOption1() string {
	if x != nil {
		return x.Option1
	}
	return ""
}

func (x *ChoiceBiRequest) GetOption2() string {
	if x != nil {
		return x.Option2
	}
	return ""
}

type ChoiceReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Option  string `protobuf:"bytes,1,opt,name=Option,proto3" json:"Option,omitempty"`
	ReplyID int32  `protobuf:"varint,2,opt,name=ReplyID,proto3" json:"ReplyID,omitempty"`
}

func (x *ChoiceReply) Reset() {
	*x = ChoiceReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cs_mw_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChoiceReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChoiceReply) ProtoMessage() {}

func (x *ChoiceReply) ProtoReflect() protoreflect.Message {
	mi := &file_cs_mw_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChoiceReply.ProtoReflect.Descriptor instead.
func (*ChoiceReply) Descriptor() ([]byte, []int) {
	return file_cs_mw_proto_rawDescGZIP(), []int{1}
}

func (x *ChoiceReply) GetOption() string {
	if x != nil {
		return x.Option
	}
	return ""
}

func (x *ChoiceReply) GetReplyID() int32 {
	if x != nil {
		return x.ReplyID
	}
	return 0
}

type NotifyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WorkerAddress string `protobuf:"bytes,1,opt,name=WorkerAddress,proto3" json:"WorkerAddress,omitempty"`
}

func (x *NotifyRequest) Reset() {
	*x = NotifyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cs_mw_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyRequest) ProtoMessage() {}

func (x *NotifyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cs_mw_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyRequest.ProtoReflect.Descriptor instead.
func (*NotifyRequest) Descriptor() ([]byte, []int) {
	return file_cs_mw_proto_rawDescGZIP(), []int{2}
}

func (x *NotifyRequest) GetWorkerAddress() string {
	if x != nil {
		return x.WorkerAddress
	}
	return ""
}

type NotifyReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result string `protobuf:"bytes,1,opt,name=Result,proto3" json:"Result,omitempty"`
}

func (x *NotifyReply) Reset() {
	*x = NotifyReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cs_mw_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyReply) ProtoMessage() {}

func (x *NotifyReply) ProtoReflect() protoreflect.Message {
	mi := &file_cs_mw_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyReply.ProtoReflect.Descriptor instead.
func (*NotifyReply) Descriptor() ([]byte, []int) {
	return file_cs_mw_proto_rawDescGZIP(), []int{3}
}

func (x *NotifyReply) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

var File_cs_mw_proto protoreflect.FileDescriptor

var file_cs_mw_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x63, 0x73, 0x2d, 0x6d, 0x77, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x6d,
	0x69, 0x6e, 0x69, 0x22, 0x45, 0x0a, 0x0f, 0x43, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x42, 0x69, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x31, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x31,
	0x12, 0x18, 0x0a, 0x07, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x32, 0x22, 0x3f, 0x0a, 0x0b, 0x43, 0x68,
	0x6f, 0x69, 0x63, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x18, 0x0a, 0x07, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x07, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x49, 0x44, 0x22, 0x35, 0x0a, 0x0d, 0x4e,
	0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x0d,
	0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x22, 0x25, 0x0a, 0x0b, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x32, 0x3d, 0x0a, 0x05, 0x46, 0x72, 0x6f,
	0x6e, 0x74, 0x12, 0x34, 0x0a, 0x06, 0x43, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x12, 0x15, 0x2e, 0x6d,
	0x69, 0x6e, 0x69, 0x2e, 0x43, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x42, 0x69, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x6d, 0x69, 0x6e, 0x69, 0x2e, 0x43, 0x68, 0x6f, 0x69, 0x63,
	0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x32, 0x3c, 0x0a, 0x04, 0x42, 0x61, 0x63, 0x6b,
	0x12, 0x34, 0x0a, 0x06, 0x43, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x12, 0x15, 0x2e, 0x6d, 0x69, 0x6e,
	0x69, 0x2e, 0x43, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x42, 0x69, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x11, 0x2e, 0x6d, 0x69, 0x6e, 0x69, 0x2e, 0x43, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x32, 0x8a, 0x01, 0x0a, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65,
	0x72, 0x12, 0x3e, 0x0a, 0x12, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x41, 0x63, 0x74, 0x69, 0x76,
	0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x12, 0x13, 0x2e, 0x6d, 0x69, 0x6e, 0x69, 0x2e, 0x4e,
	0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x6d,
	0x69, 0x6e, 0x69, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22,
	0x00, 0x12, 0x40, 0x0a, 0x14, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x44, 0x65, 0x61, 0x63, 0x74,
	0x69, 0x76, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x12, 0x13, 0x2e, 0x6d, 0x69, 0x6e, 0x69,
	0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11,
	0x2e, 0x6d, 0x69, 0x6e, 0x69, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x22, 0x00, 0x42, 0x0c, 0x5a, 0x0a, 0x6d, 0x69, 0x6e, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cs_mw_proto_rawDescOnce sync.Once
	file_cs_mw_proto_rawDescData = file_cs_mw_proto_rawDesc
)

func file_cs_mw_proto_rawDescGZIP() []byte {
	file_cs_mw_proto_rawDescOnce.Do(func() {
		file_cs_mw_proto_rawDescData = protoimpl.X.CompressGZIP(file_cs_mw_proto_rawDescData)
	})
	return file_cs_mw_proto_rawDescData
}

var file_cs_mw_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_cs_mw_proto_goTypes = []interface{}{
	(*ChoiceBiRequest)(nil), // 0: mini.ChoiceBiRequest
	(*ChoiceReply)(nil),     // 1: mini.ChoiceReply
	(*NotifyRequest)(nil),   // 2: mini.NotifyRequest
	(*NotifyReply)(nil),     // 3: mini.NotifyReply
}
var file_cs_mw_proto_depIdxs = []int32{
	0, // 0: mini.Front.Choice:input_type -> mini.ChoiceBiRequest
	0, // 1: mini.Back.Choice:input_type -> mini.ChoiceBiRequest
	2, // 2: mini.Master.NotifyActiveWorker:input_type -> mini.NotifyRequest
	2, // 3: mini.Master.NotifyDeactiveWorker:input_type -> mini.NotifyRequest
	1, // 4: mini.Front.Choice:output_type -> mini.ChoiceReply
	1, // 5: mini.Back.Choice:output_type -> mini.ChoiceReply
	3, // 6: mini.Master.NotifyActiveWorker:output_type -> mini.NotifyReply
	3, // 7: mini.Master.NotifyDeactiveWorker:output_type -> mini.NotifyReply
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_cs_mw_proto_init() }
func file_cs_mw_proto_init() {
	if File_cs_mw_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cs_mw_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChoiceBiRequest); i {
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
		file_cs_mw_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChoiceReply); i {
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
		file_cs_mw_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyRequest); i {
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
		file_cs_mw_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyReply); i {
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
			RawDescriptor: file_cs_mw_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   3,
		},
		GoTypes:           file_cs_mw_proto_goTypes,
		DependencyIndexes: file_cs_mw_proto_depIdxs,
		MessageInfos:      file_cs_mw_proto_msgTypes,
	}.Build()
	File_cs_mw_proto = out.File
	file_cs_mw_proto_rawDesc = nil
	file_cs_mw_proto_goTypes = nil
	file_cs_mw_proto_depIdxs = nil
}