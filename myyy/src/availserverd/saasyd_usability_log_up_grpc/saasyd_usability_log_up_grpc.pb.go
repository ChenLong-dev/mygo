// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: saasyd_usability_log_up_grpc.proto

package usability_log_up_l

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Code int32

const (
	Code_SUCCESS        Code = 0
	Code_FAIL           Code = 1
	Code_REQUEST_FAILED Code = 500
	Code_PARAM_ERROR    Code = 501
)

// Enum value maps for Code.
var (
	Code_name = map[int32]string{
		0:   "SUCCESS",
		1:   "FAIL",
		500: "REQUEST_FAILED",
		501: "PARAM_ERROR",
	}
	Code_value = map[string]int32{
		"SUCCESS":        0,
		"FAIL":           1,
		"REQUEST_FAILED": 500,
		"PARAM_ERROR":    501,
	}
)

func (x Code) Enum() *Code {
	p := new(Code)
	*p = x
	return p
}

func (x Code) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Code) Descriptor() protoreflect.EnumDescriptor {
	return file_saasyd_usability_log_up_grpc_proto_enumTypes[0].Descriptor()
}

func (Code) Type() protoreflect.EnumType {
	return &file_saasyd_usability_log_up_grpc_proto_enumTypes[0]
}

func (x Code) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Code.Descriptor instead.
func (Code) EnumDescriptor() ([]byte, []int) {
	return file_saasyd_usability_log_up_grpc_proto_rawDescGZIP(), []int{0}
}

type UsabilityLogReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url        string                      `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Result     string                      `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"`
	HappenTime int64                       `protobuf:"varint,3,opt,name=happen_time,json=happenTime,proto3" json:"happen_time,omitempty"`
	NodeList   []*UsabilityLogReq_NodeList `protobuf:"bytes,4,rep,name=node_list,json=nodeList,proto3" json:"node_list,omitempty"`
}

func (x *UsabilityLogReq) Reset() {
	*x = UsabilityLogReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saasyd_usability_log_up_grpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UsabilityLogReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UsabilityLogReq) ProtoMessage() {}

func (x *UsabilityLogReq) ProtoReflect() protoreflect.Message {
	mi := &file_saasyd_usability_log_up_grpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UsabilityLogReq.ProtoReflect.Descriptor instead.
func (*UsabilityLogReq) Descriptor() ([]byte, []int) {
	return file_saasyd_usability_log_up_grpc_proto_rawDescGZIP(), []int{0}
}

func (x *UsabilityLogReq) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *UsabilityLogReq) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

func (x *UsabilityLogReq) GetHappenTime() int64 {
	if x != nil {
		return x.HappenTime
	}
	return 0
}

func (x *UsabilityLogReq) GetNodeList() []*UsabilityLogReq_NodeList {
	if x != nil {
		return x.NodeList
	}
	return nil
}

type UsabilityLogRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    Code   `protobuf:"varint,1,opt,name=code,proto3,enum=usability_log_up_l.Code" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *UsabilityLogRsp) Reset() {
	*x = UsabilityLogRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saasyd_usability_log_up_grpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UsabilityLogRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UsabilityLogRsp) ProtoMessage() {}

func (x *UsabilityLogRsp) ProtoReflect() protoreflect.Message {
	mi := &file_saasyd_usability_log_up_grpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UsabilityLogRsp.ProtoReflect.Descriptor instead.
func (*UsabilityLogRsp) Descriptor() ([]byte, []int) {
	return file_saasyd_usability_log_up_grpc_proto_rawDescGZIP(), []int{1}
}

func (x *UsabilityLogRsp) GetCode() Code {
	if x != nil {
		return x.Code
	}
	return Code_SUCCESS
}

func (x *UsabilityLogRsp) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type UsabilityLogReq_NodeList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Node       string                                 `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
	LineResult []*UsabilityLogReq_NodeList_LineResult `protobuf:"bytes,2,rep,name=line_result,json=lineResult,proto3" json:"line_result,omitempty"`
}

func (x *UsabilityLogReq_NodeList) Reset() {
	*x = UsabilityLogReq_NodeList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saasyd_usability_log_up_grpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UsabilityLogReq_NodeList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UsabilityLogReq_NodeList) ProtoMessage() {}

func (x *UsabilityLogReq_NodeList) ProtoReflect() protoreflect.Message {
	mi := &file_saasyd_usability_log_up_grpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UsabilityLogReq_NodeList.ProtoReflect.Descriptor instead.
func (*UsabilityLogReq_NodeList) Descriptor() ([]byte, []int) {
	return file_saasyd_usability_log_up_grpc_proto_rawDescGZIP(), []int{0, 0}
}

func (x *UsabilityLogReq_NodeList) GetNode() string {
	if x != nil {
		return x.Node
	}
	return ""
}

func (x *UsabilityLogReq_NodeList) GetLineResult() []*UsabilityLogReq_NodeList_LineResult {
	if x != nil {
		return x.LineResult
	}
	return nil
}

type UsabilityLogReq_NodeList_LineResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addr   string `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	Isp    string `protobuf:"bytes,2,opt,name=isp,proto3" json:"isp,omitempty"`
	Code   int32  `protobuf:"varint,3,opt,name=code,proto3" json:"code,omitempty"`
	Status string `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *UsabilityLogReq_NodeList_LineResult) Reset() {
	*x = UsabilityLogReq_NodeList_LineResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saasyd_usability_log_up_grpc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UsabilityLogReq_NodeList_LineResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UsabilityLogReq_NodeList_LineResult) ProtoMessage() {}

func (x *UsabilityLogReq_NodeList_LineResult) ProtoReflect() protoreflect.Message {
	mi := &file_saasyd_usability_log_up_grpc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UsabilityLogReq_NodeList_LineResult.ProtoReflect.Descriptor instead.
func (*UsabilityLogReq_NodeList_LineResult) Descriptor() ([]byte, []int) {
	return file_saasyd_usability_log_up_grpc_proto_rawDescGZIP(), []int{0, 0, 0}
}

func (x *UsabilityLogReq_NodeList_LineResult) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *UsabilityLogReq_NodeList_LineResult) GetIsp() string {
	if x != nil {
		return x.Isp
	}
	return ""
}

func (x *UsabilityLogReq_NodeList_LineResult) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *UsabilityLogReq_NodeList_LineResult) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

var File_saasyd_usability_log_up_grpc_proto protoreflect.FileDescriptor

var file_saasyd_usability_log_up_grpc_proto_rawDesc = []byte{
	0x0a, 0x22, 0x73, 0x61, 0x61, 0x73, 0x79, 0x64, 0x5f, 0x75, 0x73, 0x61, 0x62, 0x69, 0x6c, 0x69,
	0x74, 0x79, 0x5f, 0x6c, 0x6f, 0x67, 0x5f, 0x75, 0x70, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x75, 0x73, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x5f,
	0x6c, 0x6f, 0x67, 0x5f, 0x75, 0x70, 0x5f, 0x6c, 0x22, 0x82, 0x03, 0x0a, 0x0f, 0x55, 0x73, 0x61,
	0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x16,
	0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x68, 0x61, 0x70, 0x70, 0x65, 0x6e,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x68, 0x61, 0x70,
	0x70, 0x65, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x49, 0x0a, 0x09, 0x6e, 0x6f, 0x64, 0x65, 0x5f,
	0x6c, 0x69, 0x73, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x75, 0x73, 0x61,
	0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x6c, 0x6f, 0x67, 0x5f, 0x75, 0x70, 0x5f, 0x6c, 0x2e,
	0x55, 0x73, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x71, 0x2e,
	0x4e, 0x6f, 0x64, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x4c, 0x69,
	0x73, 0x74, 0x1a, 0xd8, 0x01, 0x0a, 0x08, 0x4e, 0x6f, 0x64, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x6f, 0x64, 0x65, 0x12, 0x58, 0x0a, 0x0b, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x37, 0x2e, 0x75, 0x73, 0x61, 0x62, 0x69,
	0x6c, 0x69, 0x74, 0x79, 0x5f, 0x6c, 0x6f, 0x67, 0x5f, 0x75, 0x70, 0x5f, 0x6c, 0x2e, 0x55, 0x73,
	0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x71, 0x2e, 0x4e, 0x6f,
	0x64, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x2e, 0x4c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x52, 0x0a, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x1a, 0x5e, 0x0a,
	0x0a, 0x4c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x61,
	0x64, 0x64, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12,
	0x10, 0x0a, 0x03, 0x69, 0x73, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x69, 0x73,
	0x70, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x59, 0x0a,
	0x0f, 0x55, 0x73, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x4c, 0x6f, 0x67, 0x52, 0x73, 0x70,
	0x12, 0x2c, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x18,
	0x2e, 0x75, 0x73, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x6c, 0x6f, 0x67, 0x5f, 0x75,
	0x70, 0x5f, 0x6c, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2a, 0x44, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x0b, 0x0a, 0x07, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x00, 0x12, 0x08, 0x0a,
	0x04, 0x46, 0x41, 0x49, 0x4c, 0x10, 0x01, 0x12, 0x13, 0x0a, 0x0e, 0x52, 0x45, 0x51, 0x55, 0x45,
	0x53, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0xf4, 0x03, 0x12, 0x10, 0x0a, 0x0b,
	0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0xf5, 0x03, 0x32, 0x72,
	0x0a, 0x0e, 0x55, 0x73, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x4c, 0x6f, 0x67, 0x55, 0x70,
	0x12, 0x60, 0x0a, 0x12, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x55, 0x73, 0x61, 0x62, 0x69, 0x6c,
	0x69, 0x74, 0x79, 0x4c, 0x6f, 0x67, 0x12, 0x23, 0x2e, 0x75, 0x73, 0x61, 0x62, 0x69, 0x6c, 0x69,
	0x74, 0x79, 0x5f, 0x6c, 0x6f, 0x67, 0x5f, 0x75, 0x70, 0x5f, 0x6c, 0x2e, 0x55, 0x73, 0x61, 0x62,
	0x69, 0x6c, 0x69, 0x74, 0x79, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x71, 0x1a, 0x23, 0x2e, 0x75, 0x73,
	0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x6c, 0x6f, 0x67, 0x5f, 0x75, 0x70, 0x5f, 0x6c,
	0x2e, 0x55, 0x73, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x4c, 0x6f, 0x67, 0x52, 0x73, 0x70,
	0x22, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_saasyd_usability_log_up_grpc_proto_rawDescOnce sync.Once
	file_saasyd_usability_log_up_grpc_proto_rawDescData = file_saasyd_usability_log_up_grpc_proto_rawDesc
)

func file_saasyd_usability_log_up_grpc_proto_rawDescGZIP() []byte {
	file_saasyd_usability_log_up_grpc_proto_rawDescOnce.Do(func() {
		file_saasyd_usability_log_up_grpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_saasyd_usability_log_up_grpc_proto_rawDescData)
	})
	return file_saasyd_usability_log_up_grpc_proto_rawDescData
}

var file_saasyd_usability_log_up_grpc_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_saasyd_usability_log_up_grpc_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_saasyd_usability_log_up_grpc_proto_goTypes = []interface{}{
	(Code)(0),                                   // 0: usability_log_up_l.Code
	(*UsabilityLogReq)(nil),                     // 1: usability_log_up_l.UsabilityLogReq
	(*UsabilityLogRsp)(nil),                     // 2: usability_log_up_l.UsabilityLogRsp
	(*UsabilityLogReq_NodeList)(nil),            // 3: usability_log_up_l.UsabilityLogReq.NodeList
	(*UsabilityLogReq_NodeList_LineResult)(nil), // 4: usability_log_up_l.UsabilityLogReq.NodeList.LineResult
}
var file_saasyd_usability_log_up_grpc_proto_depIdxs = []int32{
	3, // 0: usability_log_up_l.UsabilityLogReq.node_list:type_name -> usability_log_up_l.UsabilityLogReq.NodeList
	0, // 1: usability_log_up_l.UsabilityLogRsp.code:type_name -> usability_log_up_l.Code
	4, // 2: usability_log_up_l.UsabilityLogReq.NodeList.line_result:type_name -> usability_log_up_l.UsabilityLogReq.NodeList.LineResult
	1, // 3: usability_log_up_l.UsabilityLogUp.uploadUsabilityLog:input_type -> usability_log_up_l.UsabilityLogReq
	2, // 4: usability_log_up_l.UsabilityLogUp.uploadUsabilityLog:output_type -> usability_log_up_l.UsabilityLogRsp
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_saasyd_usability_log_up_grpc_proto_init() }
func file_saasyd_usability_log_up_grpc_proto_init() {
	if File_saasyd_usability_log_up_grpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_saasyd_usability_log_up_grpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UsabilityLogReq); i {
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
		file_saasyd_usability_log_up_grpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UsabilityLogRsp); i {
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
		file_saasyd_usability_log_up_grpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UsabilityLogReq_NodeList); i {
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
		file_saasyd_usability_log_up_grpc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UsabilityLogReq_NodeList_LineResult); i {
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
			RawDescriptor: file_saasyd_usability_log_up_grpc_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_saasyd_usability_log_up_grpc_proto_goTypes,
		DependencyIndexes: file_saasyd_usability_log_up_grpc_proto_depIdxs,
		EnumInfos:         file_saasyd_usability_log_up_grpc_proto_enumTypes,
		MessageInfos:      file_saasyd_usability_log_up_grpc_proto_msgTypes,
	}.Build()
	File_saasyd_usability_log_up_grpc_proto = out.File
	file_saasyd_usability_log_up_grpc_proto_rawDesc = nil
	file_saasyd_usability_log_up_grpc_proto_goTypes = nil
	file_saasyd_usability_log_up_grpc_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// UsabilityLogUpClient is the client API for UsabilityLogUp service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UsabilityLogUpClient interface {
	UploadUsabilityLog(ctx context.Context, in *UsabilityLogReq, opts ...grpc.CallOption) (*UsabilityLogRsp, error)
}

type usabilityLogUpClient struct {
	cc grpc.ClientConnInterface
}

func NewUsabilityLogUpClient(cc grpc.ClientConnInterface) UsabilityLogUpClient {
	return &usabilityLogUpClient{cc}
}

func (c *usabilityLogUpClient) UploadUsabilityLog(ctx context.Context, in *UsabilityLogReq, opts ...grpc.CallOption) (*UsabilityLogRsp, error) {
	out := new(UsabilityLogRsp)
	err := c.cc.Invoke(ctx, "/usability_log_up_l.UsabilityLogUp/uploadUsabilityLog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UsabilityLogUpServer is the server API for UsabilityLogUp service.
type UsabilityLogUpServer interface {
	UploadUsabilityLog(context.Context, *UsabilityLogReq) (*UsabilityLogRsp, error)
}

// UnimplementedUsabilityLogUpServer can be embedded to have forward compatible implementations.
type UnimplementedUsabilityLogUpServer struct {
}

func (*UnimplementedUsabilityLogUpServer) UploadUsabilityLog(context.Context, *UsabilityLogReq) (*UsabilityLogRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadUsabilityLog not implemented")
}

func RegisterUsabilityLogUpServer(s *grpc.Server, srv UsabilityLogUpServer) {
	s.RegisterService(&_UsabilityLogUp_serviceDesc, srv)
}

func _UsabilityLogUp_UploadUsabilityLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UsabilityLogReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsabilityLogUpServer).UploadUsabilityLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/usability_log_up_l.UsabilityLogUp/UploadUsabilityLog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsabilityLogUpServer).UploadUsabilityLog(ctx, req.(*UsabilityLogReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _UsabilityLogUp_serviceDesc = grpc.ServiceDesc{
	ServiceName: "usability_log_up_l.UsabilityLogUp",
	HandlerType: (*UsabilityLogUpServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "uploadUsabilityLog",
			Handler:    _UsabilityLogUp_UploadUsabilityLog_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "saasyd_usability_log_up_grpc.proto",
}