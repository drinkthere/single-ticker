// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: path.proto

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

type BestPath struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SourceIP string `protobuf:"bytes,1,opt,name=sourceIP,proto3" json:"sourceIP,omitempty"`
	TargetIP string `protobuf:"bytes,2,opt,name=targetIP,proto3" json:"targetIP,omitempty"`
}

func (x *BestPath) Reset() {
	*x = BestPath{}
	if protoimpl.UnsafeEnabled {
		mi := &file_path_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BestPath) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BestPath) ProtoMessage() {}

func (x *BestPath) ProtoReflect() protoreflect.Message {
	mi := &file_path_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BestPath.ProtoReflect.Descriptor instead.
func (*BestPath) Descriptor() ([]byte, []int) {
	return file_path_proto_rawDescGZIP(), []int{0}
}

func (x *BestPath) GetSourceIP() string {
	if x != nil {
		return x.SourceIP
	}
	return ""
}

func (x *BestPath) GetTargetIP() string {
	if x != nil {
		return x.TargetIP
	}
	return ""
}

var File_path_proto protoreflect.FileDescriptor

var file_path_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x70, 0x61, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x42, 0x0a, 0x08,
	0x42, 0x65, 0x73, 0x74, 0x50, 0x61, 0x74, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x49, 0x50, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x49, 0x50, 0x12, 0x1a, 0x0a, 0x08, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x49, 0x50,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x49, 0x50,
	0x42, 0x0d, 0x5a, 0x0b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_path_proto_rawDescOnce sync.Once
	file_path_proto_rawDescData = file_path_proto_rawDesc
)

func file_path_proto_rawDescGZIP() []byte {
	file_path_proto_rawDescOnce.Do(func() {
		file_path_proto_rawDescData = protoimpl.X.CompressGZIP(file_path_proto_rawDescData)
	})
	return file_path_proto_rawDescData
}

var file_path_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_path_proto_goTypes = []interface{}{
	(*BestPath)(nil), // 0: BestPath
}
var file_path_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_path_proto_init() }
func file_path_proto_init() {
	if File_path_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_path_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BestPath); i {
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
			RawDescriptor: file_path_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_path_proto_goTypes,
		DependencyIndexes: file_path_proto_depIdxs,
		MessageInfos:      file_path_proto_msgTypes,
	}.Build()
	File_path_proto = out.File
	file_path_proto_rawDesc = nil
	file_path_proto_goTypes = nil
	file_path_proto_depIdxs = nil
}
