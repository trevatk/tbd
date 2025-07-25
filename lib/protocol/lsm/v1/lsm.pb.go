// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: lsm/v1/lsm.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Index struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	SortedStringTable string                 `protobuf:"bytes,1,opt,name=sorted_string_table,json=sortedStringTable,proto3" json:"sorted_string_table,omitempty"`
	Value             string                 `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *Index) Reset() {
	*x = Index{}
	mi := &file_lsm_v1_lsm_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Index) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Index) ProtoMessage() {}

func (x *Index) ProtoReflect() protoreflect.Message {
	mi := &file_lsm_v1_lsm_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Index.ProtoReflect.Descriptor instead.
func (*Index) Descriptor() ([]byte, []int) {
	return file_lsm_v1_lsm_proto_rawDescGZIP(), []int{0}
}

func (x *Index) GetSortedStringTable() string {
	if x != nil {
		return x.SortedStringTable
	}
	return ""
}

func (x *Index) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type KeyValue struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value         []byte                 `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Ttl           int64                  `protobuf:"varint,3,opt,name=ttl,proto3" json:"ttl,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *KeyValue) Reset() {
	*x = KeyValue{}
	mi := &file_lsm_v1_lsm_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *KeyValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyValue) ProtoMessage() {}

func (x *KeyValue) ProtoReflect() protoreflect.Message {
	mi := &file_lsm_v1_lsm_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyValue.ProtoReflect.Descriptor instead.
func (*KeyValue) Descriptor() ([]byte, []int) {
	return file_lsm_v1_lsm_proto_rawDescGZIP(), []int{1}
}

func (x *KeyValue) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *KeyValue) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *KeyValue) GetTtl() int64 {
	if x != nil {
		return x.Ttl
	}
	return 0
}

var File_lsm_v1_lsm_proto protoreflect.FileDescriptor

const file_lsm_v1_lsm_proto_rawDesc = "" +
	"\n" +
	"\x10lsm/v1/lsm.proto\x12\x06lsm.v1\"M\n" +
	"\x05Index\x12.\n" +
	"\x13sorted_string_table\x18\x01 \x01(\tR\x11sortedStringTable\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value\"D\n" +
	"\bKeyValue\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\fR\x05value\x12\x10\n" +
	"\x03ttl\x18\x03 \x01(\x03R\x03ttlB'Z%soft.structx.io/idp/api/gen/go/lsm/v1b\x06proto3"

var (
	file_lsm_v1_lsm_proto_rawDescOnce sync.Once
	file_lsm_v1_lsm_proto_rawDescData []byte
)

func file_lsm_v1_lsm_proto_rawDescGZIP() []byte {
	file_lsm_v1_lsm_proto_rawDescOnce.Do(func() {
		file_lsm_v1_lsm_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_lsm_v1_lsm_proto_rawDesc), len(file_lsm_v1_lsm_proto_rawDesc)))
	})
	return file_lsm_v1_lsm_proto_rawDescData
}

var file_lsm_v1_lsm_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_lsm_v1_lsm_proto_goTypes = []any{
	(*Index)(nil),    // 0: lsm.v1.Index
	(*KeyValue)(nil), // 1: lsm.v1.KeyValue
}
var file_lsm_v1_lsm_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_lsm_v1_lsm_proto_init() }
func file_lsm_v1_lsm_proto_init() {
	if File_lsm_v1_lsm_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_lsm_v1_lsm_proto_rawDesc), len(file_lsm_v1_lsm_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_lsm_v1_lsm_proto_goTypes,
		DependencyIndexes: file_lsm_v1_lsm_proto_depIdxs,
		MessageInfos:      file_lsm_v1_lsm_proto_msgTypes,
	}.Build()
	File_lsm_v1_lsm_proto = out.File
	file_lsm_v1_lsm_proto_goTypes = nil
	file_lsm_v1_lsm_proto_depIdxs = nil
}
