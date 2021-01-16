// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: saasyd_usability_business_manager_grpc.proto

package usability_business_manager

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
	Code_REQUEST_FAILED Code = 500
	Code_PARAM_ERROR    Code = 501
	Code_FAIL           Code = 1
)

// Enum value maps for Code.
var (
	Code_name = map[int32]string{
		0:   "SUCCESS",
		500: "REQUEST_FAILED",
		501: "PARAM_ERROR",
		1:   "FAIL",
	}
	Code_value = map[string]int32{
		"SUCCESS":        0,
		"REQUEST_FAILED": 500,
		"PARAM_ERROR":    501,
		"FAIL":           1,
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
	return file_saasyd_usability_business_manager_grpc_proto_enumTypes[0].Descriptor()
}

func (Code) Type() protoreflect.EnumType {
	return &file_saasyd_usability_business_manager_grpc_proto_enumTypes[0]
}

func (x Code) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Code.Descriptor instead.
func (Code) EnumDescriptor() ([]byte, []int) {
	return file_saasyd_usability_business_manager_grpc_proto_rawDescGZIP(), []int{0}
}

type BusinessTaskAddReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url  string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Freq int32  `protobuf:"varint,10,opt,name=freq,proto3" json:"freq,omitempty"`
}

func (x *BusinessTaskAddReq) Reset() {
	*x = BusinessTaskAddReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saasyd_usability_business_manager_grpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BusinessTaskAddReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BusinessTaskAddReq) ProtoMessage() {}

func (x *BusinessTaskAddReq) ProtoReflect() protoreflect.Message {
	mi := &file_saasyd_usability_business_manager_grpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BusinessTaskAddReq.ProtoReflect.Descriptor instead.
func (*BusinessTaskAddReq) Descriptor() ([]byte, []int) {
	return file_saasyd_usability_business_manager_grpc_proto_rawDescGZIP(), []int{0}
}

func (x *BusinessTaskAddReq) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *BusinessTaskAddReq) GetFreq() int32 {
	if x != nil {
		return x.Freq
	}
	return 0
}

type BusinessTaskAddRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    Code   `protobuf:"varint,1,opt,name=code,proto3,enum=usability_business_manager.Code" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *BusinessTaskAddRsp) Reset() {
	*x = BusinessTaskAddRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saasyd_usability_business_manager_grpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BusinessTaskAddRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BusinessTaskAddRsp) ProtoMessage() {}

func (x *BusinessTaskAddRsp) ProtoReflect() protoreflect.Message {
	mi := &file_saasyd_usability_business_manager_grpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BusinessTaskAddRsp.ProtoReflect.Descriptor instead.
func (*BusinessTaskAddRsp) Descriptor() ([]byte, []int) {
	return file_saasyd_usability_business_manager_grpc_proto_rawDescGZIP(), []int{1}
}

func (x *BusinessTaskAddRsp) GetCode() Code {
	if x != nil {
		return x.Code
	}
	return Code_SUCCESS
}

func (x *BusinessTaskAddRsp) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type BusinessTaskDelReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *BusinessTaskDelReq) Reset() {
	*x = BusinessTaskDelReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saasyd_usability_business_manager_grpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BusinessTaskDelReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BusinessTaskDelReq) ProtoMessage() {}

func (x *BusinessTaskDelReq) ProtoReflect() protoreflect.Message {
	mi := &file_saasyd_usability_business_manager_grpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BusinessTaskDelReq.ProtoReflect.Descriptor instead.
func (*BusinessTaskDelReq) Descriptor() ([]byte, []int) {
	return file_saasyd_usability_business_manager_grpc_proto_rawDescGZIP(), []int{2}
}

func (x *BusinessTaskDelReq) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type BusinessTaskDelRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    Code   `protobuf:"varint,1,opt,name=code,proto3,enum=usability_business_manager.Code" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *BusinessTaskDelRsp) Reset() {
	*x = BusinessTaskDelRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saasyd_usability_business_manager_grpc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BusinessTaskDelRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BusinessTaskDelRsp) ProtoMessage() {}

func (x *BusinessTaskDelRsp) ProtoReflect() protoreflect.Message {
	mi := &file_saasyd_usability_business_manager_grpc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BusinessTaskDelRsp.ProtoReflect.Descriptor instead.
func (*BusinessTaskDelRsp) Descriptor() ([]byte, []int) {
	return file_saasyd_usability_business_manager_grpc_proto_rawDescGZIP(), []int{3}
}

func (x *BusinessTaskDelRsp) GetCode() Code {
	if x != nil {
		return x.Code
	}
	return Code_SUCCESS
}

func (x *BusinessTaskDelRsp) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type BusinessTaskEditReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url  string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Freq int32  `protobuf:"varint,10,opt,name=freq,proto3" json:"freq,omitempty"`
}

func (x *BusinessTaskEditReq) Reset() {
	*x = BusinessTaskEditReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saasyd_usability_business_manager_grpc_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BusinessTaskEditReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BusinessTaskEditReq) ProtoMessage() {}

func (x *BusinessTaskEditReq) ProtoReflect() protoreflect.Message {
	mi := &file_saasyd_usability_business_manager_grpc_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BusinessTaskEditReq.ProtoReflect.Descriptor instead.
func (*BusinessTaskEditReq) Descriptor() ([]byte, []int) {
	return file_saasyd_usability_business_manager_grpc_proto_rawDescGZIP(), []int{4}
}

func (x *BusinessTaskEditReq) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *BusinessTaskEditReq) GetFreq() int32 {
	if x != nil {
		return x.Freq
	}
	return 0
}

type BusinessTaskEditRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    Code   `protobuf:"varint,1,opt,name=code,proto3,enum=usability_business_manager.Code" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *BusinessTaskEditRsp) Reset() {
	*x = BusinessTaskEditRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saasyd_usability_business_manager_grpc_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BusinessTaskEditRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BusinessTaskEditRsp) ProtoMessage() {}

func (x *BusinessTaskEditRsp) ProtoReflect() protoreflect.Message {
	mi := &file_saasyd_usability_business_manager_grpc_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BusinessTaskEditRsp.ProtoReflect.Descriptor instead.
func (*BusinessTaskEditRsp) Descriptor() ([]byte, []int) {
	return file_saasyd_usability_business_manager_grpc_proto_rawDescGZIP(), []int{5}
}

func (x *BusinessTaskEditRsp) GetCode() Code {
	if x != nil {
		return x.Code
	}
	return Code_SUCCESS
}

func (x *BusinessTaskEditRsp) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_saasyd_usability_business_manager_grpc_proto protoreflect.FileDescriptor

var file_saasyd_usability_business_manager_grpc_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x73, 0x61, 0x61, 0x73, 0x79, 0x64, 0x5f, 0x75, 0x73, 0x61, 0x62, 0x69, 0x6c, 0x69,
	0x74, 0x79, 0x5f, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x5f, 0x6d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x72, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a,
	0x75, 0x73, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65,
	0x73, 0x73, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x22, 0x3a, 0x0a, 0x12, 0x42, 0x75,
	0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x54, 0x61, 0x73, 0x6b, 0x41, 0x64, 0x64, 0x52, 0x65, 0x71,
	0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75,
	0x72, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x65, 0x71, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x04, 0x66, 0x72, 0x65, 0x71, 0x22, 0x64, 0x0a, 0x12, 0x42, 0x75, 0x73, 0x69, 0x6e, 0x65,
	0x73, 0x73, 0x54, 0x61, 0x73, 0x6b, 0x41, 0x64, 0x64, 0x52, 0x73, 0x70, 0x12, 0x34, 0x0a, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x75, 0x73, 0x61,
	0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x5f,
	0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x26, 0x0a, 0x12,
	0x42, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x54, 0x61, 0x73, 0x6b, 0x44, 0x65, 0x6c, 0x52,
	0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x75, 0x72, 0x6c, 0x22, 0x64, 0x0a, 0x12, 0x42, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73,
	0x54, 0x61, 0x73, 0x6b, 0x44, 0x65, 0x6c, 0x52, 0x73, 0x70, 0x12, 0x34, 0x0a, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x75, 0x73, 0x61, 0x62, 0x69,
	0x6c, 0x69, 0x74, 0x79, 0x5f, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x5f, 0x6d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x3b, 0x0a, 0x13, 0x42, 0x75,
	0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x54, 0x61, 0x73, 0x6b, 0x45, 0x64, 0x69, 0x74, 0x52, 0x65,
	0x71, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x75, 0x72, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x65, 0x71, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x04, 0x66, 0x72, 0x65, 0x71, 0x22, 0x65, 0x0a, 0x13, 0x42, 0x75, 0x73, 0x69, 0x6e,
	0x65, 0x73, 0x73, 0x54, 0x61, 0x73, 0x6b, 0x45, 0x64, 0x69, 0x74, 0x52, 0x73, 0x70, 0x12, 0x34,
	0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x75,
	0x73, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73,
	0x73, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2a, 0x44,
	0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53,
	0x53, 0x10, 0x00, 0x12, 0x13, 0x0a, 0x0e, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x5f, 0x46,
	0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0xf4, 0x03, 0x12, 0x10, 0x0a, 0x0b, 0x50, 0x41, 0x52, 0x41,
	0x4d, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0xf5, 0x03, 0x12, 0x08, 0x0a, 0x04, 0x46, 0x41,
	0x49, 0x4c, 0x10, 0x01, 0x32, 0xf7, 0x02, 0x0a, 0x13, 0x42, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73,
	0x73, 0x54, 0x61, 0x73, 0x6b, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x12, 0x73, 0x0a, 0x0f,
	0x42, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x54, 0x61, 0x73, 0x6b, 0x41, 0x64, 0x64, 0x12,
	0x2e, 0x2e, 0x75, 0x73, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x62, 0x75, 0x73, 0x69,
	0x6e, 0x65, 0x73, 0x73, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x42, 0x75, 0x73,
	0x69, 0x6e, 0x65, 0x73, 0x73, 0x54, 0x61, 0x73, 0x6b, 0x41, 0x64, 0x64, 0x52, 0x65, 0x71, 0x1a,
	0x2e, 0x2e, 0x75, 0x73, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x62, 0x75, 0x73, 0x69,
	0x6e, 0x65, 0x73, 0x73, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x42, 0x75, 0x73,
	0x69, 0x6e, 0x65, 0x73, 0x73, 0x54, 0x61, 0x73, 0x6b, 0x41, 0x64, 0x64, 0x52, 0x73, 0x70, 0x22,
	0x00, 0x12, 0x73, 0x0a, 0x0f, 0x42, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x54, 0x61, 0x73,
	0x6b, 0x44, 0x65, 0x6c, 0x12, 0x2e, 0x2e, 0x75, 0x73, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79,
	0x5f, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x72, 0x2e, 0x42, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x54, 0x61, 0x73, 0x6b, 0x44, 0x65,
	0x6c, 0x52, 0x65, 0x71, 0x1a, 0x2e, 0x2e, 0x75, 0x73, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79,
	0x5f, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x72, 0x2e, 0x42, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x54, 0x61, 0x73, 0x6b, 0x44, 0x65,
	0x6c, 0x52, 0x73, 0x70, 0x22, 0x00, 0x12, 0x76, 0x0a, 0x10, 0x42, 0x75, 0x73, 0x69, 0x6e, 0x65,
	0x73, 0x73, 0x54, 0x61, 0x73, 0x6b, 0x45, 0x64, 0x69, 0x74, 0x12, 0x2f, 0x2e, 0x75, 0x73, 0x61,
	0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x5f,
	0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x42, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73,
	0x54, 0x61, 0x73, 0x6b, 0x45, 0x64, 0x69, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x2f, 0x2e, 0x75, 0x73,
	0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73,
	0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x42, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73,
	0x73, 0x54, 0x61, 0x73, 0x6b, 0x45, 0x64, 0x69, 0x74, 0x52, 0x73, 0x70, 0x22, 0x00, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_saasyd_usability_business_manager_grpc_proto_rawDescOnce sync.Once
	file_saasyd_usability_business_manager_grpc_proto_rawDescData = file_saasyd_usability_business_manager_grpc_proto_rawDesc
)

func file_saasyd_usability_business_manager_grpc_proto_rawDescGZIP() []byte {
	file_saasyd_usability_business_manager_grpc_proto_rawDescOnce.Do(func() {
		file_saasyd_usability_business_manager_grpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_saasyd_usability_business_manager_grpc_proto_rawDescData)
	})
	return file_saasyd_usability_business_manager_grpc_proto_rawDescData
}

var file_saasyd_usability_business_manager_grpc_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_saasyd_usability_business_manager_grpc_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_saasyd_usability_business_manager_grpc_proto_goTypes = []interface{}{
	(Code)(0),                   // 0: usability_business_manager.Code
	(*BusinessTaskAddReq)(nil),  // 1: usability_business_manager.BusinessTaskAddReq
	(*BusinessTaskAddRsp)(nil),  // 2: usability_business_manager.BusinessTaskAddRsp
	(*BusinessTaskDelReq)(nil),  // 3: usability_business_manager.BusinessTaskDelReq
	(*BusinessTaskDelRsp)(nil),  // 4: usability_business_manager.BusinessTaskDelRsp
	(*BusinessTaskEditReq)(nil), // 5: usability_business_manager.BusinessTaskEditReq
	(*BusinessTaskEditRsp)(nil), // 6: usability_business_manager.BusinessTaskEditRsp
}
var file_saasyd_usability_business_manager_grpc_proto_depIdxs = []int32{
	0, // 0: usability_business_manager.BusinessTaskAddRsp.code:type_name -> usability_business_manager.Code
	0, // 1: usability_business_manager.BusinessTaskDelRsp.code:type_name -> usability_business_manager.Code
	0, // 2: usability_business_manager.BusinessTaskEditRsp.code:type_name -> usability_business_manager.Code
	1, // 3: usability_business_manager.BusinessTaskManager.BusinessTaskAdd:input_type -> usability_business_manager.BusinessTaskAddReq
	3, // 4: usability_business_manager.BusinessTaskManager.BusinessTaskDel:input_type -> usability_business_manager.BusinessTaskDelReq
	5, // 5: usability_business_manager.BusinessTaskManager.BusinessTaskEdit:input_type -> usability_business_manager.BusinessTaskEditReq
	2, // 6: usability_business_manager.BusinessTaskManager.BusinessTaskAdd:output_type -> usability_business_manager.BusinessTaskAddRsp
	4, // 7: usability_business_manager.BusinessTaskManager.BusinessTaskDel:output_type -> usability_business_manager.BusinessTaskDelRsp
	6, // 8: usability_business_manager.BusinessTaskManager.BusinessTaskEdit:output_type -> usability_business_manager.BusinessTaskEditRsp
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_saasyd_usability_business_manager_grpc_proto_init() }
func file_saasyd_usability_business_manager_grpc_proto_init() {
	if File_saasyd_usability_business_manager_grpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_saasyd_usability_business_manager_grpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BusinessTaskAddReq); i {
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
		file_saasyd_usability_business_manager_grpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BusinessTaskAddRsp); i {
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
		file_saasyd_usability_business_manager_grpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BusinessTaskDelReq); i {
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
		file_saasyd_usability_business_manager_grpc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BusinessTaskDelRsp); i {
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
		file_saasyd_usability_business_manager_grpc_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BusinessTaskEditReq); i {
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
		file_saasyd_usability_business_manager_grpc_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BusinessTaskEditRsp); i {
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
			RawDescriptor: file_saasyd_usability_business_manager_grpc_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_saasyd_usability_business_manager_grpc_proto_goTypes,
		DependencyIndexes: file_saasyd_usability_business_manager_grpc_proto_depIdxs,
		EnumInfos:         file_saasyd_usability_business_manager_grpc_proto_enumTypes,
		MessageInfos:      file_saasyd_usability_business_manager_grpc_proto_msgTypes,
	}.Build()
	File_saasyd_usability_business_manager_grpc_proto = out.File
	file_saasyd_usability_business_manager_grpc_proto_rawDesc = nil
	file_saasyd_usability_business_manager_grpc_proto_goTypes = nil
	file_saasyd_usability_business_manager_grpc_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// BusinessTaskManagerClient is the client API for BusinessTaskManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BusinessTaskManagerClient interface {
	BusinessTaskAdd(ctx context.Context, in *BusinessTaskAddReq, opts ...grpc.CallOption) (*BusinessTaskAddRsp, error)
	BusinessTaskDel(ctx context.Context, in *BusinessTaskDelReq, opts ...grpc.CallOption) (*BusinessTaskDelRsp, error)
	BusinessTaskEdit(ctx context.Context, in *BusinessTaskEditReq, opts ...grpc.CallOption) (*BusinessTaskEditRsp, error)
}

type businessTaskManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewBusinessTaskManagerClient(cc grpc.ClientConnInterface) BusinessTaskManagerClient {
	return &businessTaskManagerClient{cc}
}

func (c *businessTaskManagerClient) BusinessTaskAdd(ctx context.Context, in *BusinessTaskAddReq, opts ...grpc.CallOption) (*BusinessTaskAddRsp, error) {
	out := new(BusinessTaskAddRsp)
	err := c.cc.Invoke(ctx, "/usability_business_manager.BusinessTaskManager/BusinessTaskAdd", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *businessTaskManagerClient) BusinessTaskDel(ctx context.Context, in *BusinessTaskDelReq, opts ...grpc.CallOption) (*BusinessTaskDelRsp, error) {
	out := new(BusinessTaskDelRsp)
	err := c.cc.Invoke(ctx, "/usability_business_manager.BusinessTaskManager/BusinessTaskDel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *businessTaskManagerClient) BusinessTaskEdit(ctx context.Context, in *BusinessTaskEditReq, opts ...grpc.CallOption) (*BusinessTaskEditRsp, error) {
	out := new(BusinessTaskEditRsp)
	err := c.cc.Invoke(ctx, "/usability_business_manager.BusinessTaskManager/BusinessTaskEdit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BusinessTaskManagerServer is the server API for BusinessTaskManager service.
type BusinessTaskManagerServer interface {
	BusinessTaskAdd(context.Context, *BusinessTaskAddReq) (*BusinessTaskAddRsp, error)
	BusinessTaskDel(context.Context, *BusinessTaskDelReq) (*BusinessTaskDelRsp, error)
	BusinessTaskEdit(context.Context, *BusinessTaskEditReq) (*BusinessTaskEditRsp, error)
}

// UnimplementedBusinessTaskManagerServer can be embedded to have forward compatible implementations.
type UnimplementedBusinessTaskManagerServer struct {
}

func (*UnimplementedBusinessTaskManagerServer) BusinessTaskAdd(context.Context, *BusinessTaskAddReq) (*BusinessTaskAddRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BusinessTaskAdd not implemented")
}
func (*UnimplementedBusinessTaskManagerServer) BusinessTaskDel(context.Context, *BusinessTaskDelReq) (*BusinessTaskDelRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BusinessTaskDel not implemented")
}
func (*UnimplementedBusinessTaskManagerServer) BusinessTaskEdit(context.Context, *BusinessTaskEditReq) (*BusinessTaskEditRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BusinessTaskEdit not implemented")
}

func RegisterBusinessTaskManagerServer(s *grpc.Server, srv BusinessTaskManagerServer) {
	s.RegisterService(&_BusinessTaskManager_serviceDesc, srv)
}

func _BusinessTaskManager_BusinessTaskAdd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BusinessTaskAddReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BusinessTaskManagerServer).BusinessTaskAdd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/usability_business_manager.BusinessTaskManager/BusinessTaskAdd",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BusinessTaskManagerServer).BusinessTaskAdd(ctx, req.(*BusinessTaskAddReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BusinessTaskManager_BusinessTaskDel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BusinessTaskDelReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BusinessTaskManagerServer).BusinessTaskDel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/usability_business_manager.BusinessTaskManager/BusinessTaskDel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BusinessTaskManagerServer).BusinessTaskDel(ctx, req.(*BusinessTaskDelReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BusinessTaskManager_BusinessTaskEdit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BusinessTaskEditReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BusinessTaskManagerServer).BusinessTaskEdit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/usability_business_manager.BusinessTaskManager/BusinessTaskEdit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BusinessTaskManagerServer).BusinessTaskEdit(ctx, req.(*BusinessTaskEditReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _BusinessTaskManager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "usability_business_manager.BusinessTaskManager",
	HandlerType: (*BusinessTaskManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "BusinessTaskAdd",
			Handler:    _BusinessTaskManager_BusinessTaskAdd_Handler,
		},
		{
			MethodName: "BusinessTaskDel",
			Handler:    _BusinessTaskManager_BusinessTaskDel_Handler,
		},
		{
			MethodName: "BusinessTaskEdit",
			Handler:    _BusinessTaskManager_BusinessTaskEdit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "saasyd_usability_business_manager_grpc.proto",
}
