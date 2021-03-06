// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.1
// source: wordscore.proto

package pb

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

type MTimeStampInterval_MTimeFrameType int32

const (
	MTimeStampInterval_TFUnknown MTimeStampInterval_MTimeFrameType = 0
	MTimeStampInterval_TFWeek    MTimeStampInterval_MTimeFrameType = 1
	MTimeStampInterval_TFMonth   MTimeStampInterval_MTimeFrameType = 2
	MTimeStampInterval_TFQuarter MTimeStampInterval_MTimeFrameType = 3
	MTimeStampInterval_TFYear    MTimeStampInterval_MTimeFrameType = 4
	MTimeStampInterval_TFTerm    MTimeStampInterval_MTimeFrameType = 5
	MTimeStampInterval_TFSpan    MTimeStampInterval_MTimeFrameType = 6
)

// Enum value maps for MTimeStampInterval_MTimeFrameType.
var (
	MTimeStampInterval_MTimeFrameType_name = map[int32]string{
		0: "TFUnknown",
		1: "TFWeek",
		2: "TFMonth",
		3: "TFQuarter",
		4: "TFYear",
		5: "TFTerm",
		6: "TFSpan",
	}
	MTimeStampInterval_MTimeFrameType_value = map[string]int32{
		"TFUnknown": 0,
		"TFWeek":    1,
		"TFMonth":   2,
		"TFQuarter": 3,
		"TFYear":    4,
		"TFTerm":    5,
		"TFSpan":    6,
	}
)

func (x MTimeStampInterval_MTimeFrameType) Enum() *MTimeStampInterval_MTimeFrameType {
	p := new(MTimeStampInterval_MTimeFrameType)
	*p = x
	return p
}

func (x MTimeStampInterval_MTimeFrameType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MTimeStampInterval_MTimeFrameType) Descriptor() protoreflect.EnumDescriptor {
	return file_wordscore_proto_enumTypes[0].Descriptor()
}

func (MTimeStampInterval_MTimeFrameType) Type() protoreflect.EnumType {
	return &file_wordscore_proto_enumTypes[0]
}

func (x MTimeStampInterval_MTimeFrameType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MTimeStampInterval_MTimeFrameType.Descriptor instead.
func (MTimeStampInterval_MTimeFrameType) EnumDescriptor() ([]byte, []int) {
	return file_wordscore_proto_rawDescGZIP(), []int{1, 0}
}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int32  `protobuf:"varint,1,opt,name=Code,proto3" json:"Code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wordscore_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_wordscore_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_wordscore_proto_rawDescGZIP(), []int{0}
}

func (x *Error) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Error) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type MTimeStampInterval struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timeframetype MTimeStampInterval_MTimeFrameType `protobuf:"varint,1,opt,name=Timeframetype,proto3,enum=wordscore.MTimeStampInterval_MTimeFrameType" json:"Timeframetype,omitempty"`
	StartTime     *timestamppb.Timestamp            `protobuf:"bytes,2,opt,name=StartTime,proto3" json:"StartTime,omitempty"`
	EndTime       *timestamppb.Timestamp            `protobuf:"bytes,3,opt,name=EndTime,proto3" json:"EndTime,omitempty"`
}

func (x *MTimeStampInterval) Reset() {
	*x = MTimeStampInterval{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wordscore_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MTimeStampInterval) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MTimeStampInterval) ProtoMessage() {}

func (x *MTimeStampInterval) ProtoReflect() protoreflect.Message {
	mi := &file_wordscore_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MTimeStampInterval.ProtoReflect.Descriptor instead.
func (*MTimeStampInterval) Descriptor() ([]byte, []int) {
	return file_wordscore_proto_rawDescGZIP(), []int{1}
}

func (x *MTimeStampInterval) GetTimeframetype() MTimeStampInterval_MTimeFrameType {
	if x != nil {
		return x.Timeframetype
	}
	return MTimeStampInterval_TFUnknown
}

func (x *MTimeStampInterval) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *MTimeStampInterval) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

type TimeEventRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic             string              `protobuf:"bytes,1,opt,name=Topic,proto3" json:"Topic,omitempty"`
	Timestampinterval *MTimeStampInterval `protobuf:"bytes,2,opt,name=Timestampinterval,proto3" json:"Timestampinterval,omitempty"`
}

func (x *TimeEventRequest) Reset() {
	*x = TimeEventRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wordscore_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeEventRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeEventRequest) ProtoMessage() {}

func (x *TimeEventRequest) ProtoReflect() protoreflect.Message {
	mi := &file_wordscore_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeEventRequest.ProtoReflect.Descriptor instead.
func (*TimeEventRequest) Descriptor() ([]byte, []int) {
	return file_wordscore_proto_rawDescGZIP(), []int{2}
}

func (x *TimeEventRequest) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *TimeEventRequest) GetTimestampinterval() *MTimeStampInterval {
	if x != nil {
		return x.Timestampinterval
	}
	return nil
}

type TimeEventResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Completed bool   `protobuf:"varint,1,opt,name=Completed,proto3" json:"Completed,omitempty"`
	Error     *Error `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *TimeEventResponse) Reset() {
	*x = TimeEventResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wordscore_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeEventResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeEventResponse) ProtoMessage() {}

func (x *TimeEventResponse) ProtoReflect() protoreflect.Message {
	mi := &file_wordscore_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeEventResponse.ProtoReflect.Descriptor instead.
func (*TimeEventResponse) Descriptor() ([]byte, []int) {
	return file_wordscore_proto_rawDescGZIP(), []int{3}
}

func (x *TimeEventResponse) GetCompleted() bool {
	if x != nil {
		return x.Completed
	}
	return false
}

func (x *TimeEventResponse) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

type GetWordScoreRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Word         string              `protobuf:"bytes,1,opt,name=Word,proto3" json:"Word,omitempty"`
	Timeinterval *MTimeStampInterval `protobuf:"bytes,2,opt,name=Timeinterval,proto3" json:"Timeinterval,omitempty"`
}

func (x *GetWordScoreRequest) Reset() {
	*x = GetWordScoreRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wordscore_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetWordScoreRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetWordScoreRequest) ProtoMessage() {}

func (x *GetWordScoreRequest) ProtoReflect() protoreflect.Message {
	mi := &file_wordscore_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetWordScoreRequest.ProtoReflect.Descriptor instead.
func (*GetWordScoreRequest) Descriptor() ([]byte, []int) {
	return file_wordscore_proto_rawDescGZIP(), []int{4}
}

func (x *GetWordScoreRequest) GetWord() string {
	if x != nil {
		return x.Word
	}
	return ""
}

func (x *GetWordScoreRequest) GetTimeinterval() *MTimeStampInterval {
	if x != nil {
		return x.Timeinterval
	}
	return nil
}

type GetWordScoreResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           int32               `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Word         string              `protobuf:"bytes,2,opt,name=Word,proto3" json:"Word,omitempty"`
	Timeinterval *MTimeStampInterval `protobuf:"bytes,3,opt,name=Timeinterval,proto3" json:"Timeinterval,omitempty"`
	Density      float32             `protobuf:"fixed32,4,opt,name=Density,proto3" json:"Density,omitempty"`
	Linkage      float32             `protobuf:"fixed32,5,opt,name=Linkage,proto3" json:"Linkage,omitempty"`
	Growth       float32             `protobuf:"fixed32,6,opt,name=Growth,proto3" json:"Growth,omitempty"`
	Score        float32             `protobuf:"fixed32,7,opt,name=Score,proto3" json:"Score,omitempty"`
}

func (x *GetWordScoreResponse) Reset() {
	*x = GetWordScoreResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wordscore_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetWordScoreResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetWordScoreResponse) ProtoMessage() {}

func (x *GetWordScoreResponse) ProtoReflect() protoreflect.Message {
	mi := &file_wordscore_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetWordScoreResponse.ProtoReflect.Descriptor instead.
func (*GetWordScoreResponse) Descriptor() ([]byte, []int) {
	return file_wordscore_proto_rawDescGZIP(), []int{5}
}

func (x *GetWordScoreResponse) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *GetWordScoreResponse) GetWord() string {
	if x != nil {
		return x.Word
	}
	return ""
}

func (x *GetWordScoreResponse) GetTimeinterval() *MTimeStampInterval {
	if x != nil {
		return x.Timeinterval
	}
	return nil
}

func (x *GetWordScoreResponse) GetDensity() float32 {
	if x != nil {
		return x.Density
	}
	return 0
}

func (x *GetWordScoreResponse) GetLinkage() float32 {
	if x != nil {
		return x.Linkage
	}
	return 0
}

func (x *GetWordScoreResponse) GetGrowth() float32 {
	if x != nil {
		return x.Growth
	}
	return 0
}

func (x *GetWordScoreResponse) GetScore() float32 {
	if x != nil {
		return x.Score
	}
	return 0
}

type CreateWordScoreRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Word         string              `protobuf:"bytes,1,opt,name=Word,proto3" json:"Word,omitempty"`
	Timeinterval *MTimeStampInterval `protobuf:"bytes,2,opt,name=Timeinterval,proto3" json:"Timeinterval,omitempty"`
}

func (x *CreateWordScoreRequest) Reset() {
	*x = CreateWordScoreRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wordscore_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateWordScoreRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateWordScoreRequest) ProtoMessage() {}

func (x *CreateWordScoreRequest) ProtoReflect() protoreflect.Message {
	mi := &file_wordscore_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateWordScoreRequest.ProtoReflect.Descriptor instead.
func (*CreateWordScoreRequest) Descriptor() ([]byte, []int) {
	return file_wordscore_proto_rawDescGZIP(), []int{6}
}

func (x *CreateWordScoreRequest) GetWord() string {
	if x != nil {
		return x.Word
	}
	return ""
}

func (x *CreateWordScoreRequest) GetTimeinterval() *MTimeStampInterval {
	if x != nil {
		return x.Timeinterval
	}
	return nil
}

type CreateWordScoreResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error *Error `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *CreateWordScoreResponse) Reset() {
	*x = CreateWordScoreResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wordscore_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateWordScoreResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateWordScoreResponse) ProtoMessage() {}

func (x *CreateWordScoreResponse) ProtoReflect() protoreflect.Message {
	mi := &file_wordscore_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateWordScoreResponse.ProtoReflect.Descriptor instead.
func (*CreateWordScoreResponse) Descriptor() ([]byte, []int) {
	return file_wordscore_proto_rawDescGZIP(), []int{7}
}

func (x *CreateWordScoreResponse) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

var File_wordscore_proto protoreflect.FileDescriptor

var file_wordscore_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x35, 0x0a,
	0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0xc5, 0x02, 0x0a, 0x12, 0x4d, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74,
	0x61, 0x6d, 0x70, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x52, 0x0a, 0x0d, 0x54,
	0x69, 0x6d, 0x65, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x2c, 0x2e, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x4d,
	0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61,
	0x6c, 0x2e, 0x4d, 0x54, 0x69, 0x6d, 0x65, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x0d, 0x54, 0x69, 0x6d, 0x65, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x38, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x34, 0x0a, 0x07, 0x45, 0x6e, 0x64,
	0x54, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x45, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x22,
	0x6b, 0x0a, 0x0e, 0x4d, 0x54, 0x69, 0x6d, 0x65, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x0d, 0x0a, 0x09, 0x54, 0x46, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x00,
	0x12, 0x0a, 0x0a, 0x06, 0x54, 0x46, 0x57, 0x65, 0x65, 0x6b, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07,
	0x54, 0x46, 0x4d, 0x6f, 0x6e, 0x74, 0x68, 0x10, 0x02, 0x12, 0x0d, 0x0a, 0x09, 0x54, 0x46, 0x51,
	0x75, 0x61, 0x72, 0x74, 0x65, 0x72, 0x10, 0x03, 0x12, 0x0a, 0x0a, 0x06, 0x54, 0x46, 0x59, 0x65,
	0x61, 0x72, 0x10, 0x04, 0x12, 0x0a, 0x0a, 0x06, 0x54, 0x46, 0x54, 0x65, 0x72, 0x6d, 0x10, 0x05,
	0x12, 0x0a, 0x0a, 0x06, 0x54, 0x46, 0x53, 0x70, 0x61, 0x6e, 0x10, 0x06, 0x22, 0x75, 0x0a, 0x10,
	0x54, 0x69, 0x6d, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x4b, 0x0a, 0x11, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1d, 0x2e, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x4d, 0x54,
	0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c,
	0x52, 0x11, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x76, 0x61, 0x6c, 0x22, 0x59, 0x0a, 0x11, 0x54, 0x69, 0x6d, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x43, 0x6f, 0x6d, 0x70,
	0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x43, 0x6f, 0x6d,
	0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x12, 0x26, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x6c,
	0x0a, 0x13, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x64, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x57, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x57, 0x6f, 0x72, 0x64, 0x12, 0x41, 0x0a, 0x0c, 0x54, 0x69, 0x6d,
	0x65, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1d, 0x2e, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x4d, 0x54, 0x69, 0x6d,
	0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x52, 0x0c,
	0x54, 0x69, 0x6d, 0x65, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x22, 0xdf, 0x01, 0x0a,
	0x14, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x64, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x02, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x57, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x57, 0x6f, 0x72, 0x64, 0x12, 0x41, 0x0a, 0x0c, 0x54, 0x69, 0x6d,
	0x65, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1d, 0x2e, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x4d, 0x54, 0x69, 0x6d,
	0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x52, 0x0c,
	0x54, 0x69, 0x6d, 0x65, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x18, 0x0a, 0x07,
	0x44, 0x65, 0x6e, 0x73, 0x69, 0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x52, 0x07, 0x44,
	0x65, 0x6e, 0x73, 0x69, 0x74, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x4c, 0x69, 0x6e, 0x6b, 0x61, 0x67,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x02, 0x52, 0x07, 0x4c, 0x69, 0x6e, 0x6b, 0x61, 0x67, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x47, 0x72, 0x6f, 0x77, 0x74, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28, 0x02,
	0x52, 0x06, 0x47, 0x72, 0x6f, 0x77, 0x74, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x53, 0x63, 0x6f, 0x72,
	0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x22, 0x6f,
	0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x64, 0x53, 0x63, 0x6f, 0x72,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x57, 0x6f, 0x72, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x57, 0x6f, 0x72, 0x64, 0x12, 0x41, 0x0a, 0x0c,
	0x54, 0x69, 0x6d, 0x65, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x4d,
	0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61,
	0x6c, 0x52, 0x0c, 0x54, 0x69, 0x6d, 0x65, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x22,
	0x41, 0x0a, 0x17, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x64, 0x53, 0x63, 0x6f,
	0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x05, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x77, 0x6f, 0x72, 0x64,
	0x73, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x05, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x32, 0xcd, 0x01, 0x0a, 0x1c, 0x57, 0x6f, 0x72, 0x64, 0x53, 0x63, 0x6f, 0x72, 0x65,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x70, 0x63, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66,
	0x61, 0x63, 0x65, 0x12, 0x51, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x64, 0x53, 0x63,
	0x6f, 0x72, 0x65, 0x12, 0x1e, 0x2e, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x64, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x64, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x5a, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x57, 0x6f, 0x72, 0x64, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x21, 0x2e, 0x77, 0x6f, 0x72, 0x64,
	0x73, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x64,
	0x53, 0x63, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x77,
	0x6f, 0x72, 0x64, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x57,
	0x6f, 0x72, 0x64, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_wordscore_proto_rawDescOnce sync.Once
	file_wordscore_proto_rawDescData = file_wordscore_proto_rawDesc
)

func file_wordscore_proto_rawDescGZIP() []byte {
	file_wordscore_proto_rawDescOnce.Do(func() {
		file_wordscore_proto_rawDescData = protoimpl.X.CompressGZIP(file_wordscore_proto_rawDescData)
	})
	return file_wordscore_proto_rawDescData
}

var file_wordscore_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_wordscore_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_wordscore_proto_goTypes = []interface{}{
	(MTimeStampInterval_MTimeFrameType)(0), // 0: wordscore.MTimeStampInterval.MTimeFrameType
	(*Error)(nil),                          // 1: wordscore.Error
	(*MTimeStampInterval)(nil),             // 2: wordscore.MTimeStampInterval
	(*TimeEventRequest)(nil),               // 3: wordscore.TimeEventRequest
	(*TimeEventResponse)(nil),              // 4: wordscore.TimeEventResponse
	(*GetWordScoreRequest)(nil),            // 5: wordscore.GetWordScoreRequest
	(*GetWordScoreResponse)(nil),           // 6: wordscore.GetWordScoreResponse
	(*CreateWordScoreRequest)(nil),         // 7: wordscore.CreateWordScoreRequest
	(*CreateWordScoreResponse)(nil),        // 8: wordscore.CreateWordScoreResponse
	(*timestamppb.Timestamp)(nil),          // 9: google.protobuf.Timestamp
}
var file_wordscore_proto_depIdxs = []int32{
	0,  // 0: wordscore.MTimeStampInterval.Timeframetype:type_name -> wordscore.MTimeStampInterval.MTimeFrameType
	9,  // 1: wordscore.MTimeStampInterval.StartTime:type_name -> google.protobuf.Timestamp
	9,  // 2: wordscore.MTimeStampInterval.EndTime:type_name -> google.protobuf.Timestamp
	2,  // 3: wordscore.TimeEventRequest.Timestampinterval:type_name -> wordscore.MTimeStampInterval
	1,  // 4: wordscore.TimeEventResponse.error:type_name -> wordscore.Error
	2,  // 5: wordscore.GetWordScoreRequest.Timeinterval:type_name -> wordscore.MTimeStampInterval
	2,  // 6: wordscore.GetWordScoreResponse.Timeinterval:type_name -> wordscore.MTimeStampInterval
	2,  // 7: wordscore.CreateWordScoreRequest.Timeinterval:type_name -> wordscore.MTimeStampInterval
	1,  // 8: wordscore.CreateWordScoreResponse.error:type_name -> wordscore.Error
	5,  // 9: wordscore.WordScoreServiceRpcInterface.GetWordScore:input_type -> wordscore.GetWordScoreRequest
	7,  // 10: wordscore.WordScoreServiceRpcInterface.CreateWordScore:input_type -> wordscore.CreateWordScoreRequest
	6,  // 11: wordscore.WordScoreServiceRpcInterface.GetWordScore:output_type -> wordscore.GetWordScoreResponse
	8,  // 12: wordscore.WordScoreServiceRpcInterface.CreateWordScore:output_type -> wordscore.CreateWordScoreResponse
	11, // [11:13] is the sub-list for method output_type
	9,  // [9:11] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_wordscore_proto_init() }
func file_wordscore_proto_init() {
	if File_wordscore_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_wordscore_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Error); i {
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
		file_wordscore_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MTimeStampInterval); i {
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
		file_wordscore_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimeEventRequest); i {
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
		file_wordscore_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimeEventResponse); i {
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
		file_wordscore_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetWordScoreRequest); i {
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
		file_wordscore_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetWordScoreResponse); i {
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
		file_wordscore_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateWordScoreRequest); i {
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
		file_wordscore_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateWordScoreResponse); i {
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
			RawDescriptor: file_wordscore_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_wordscore_proto_goTypes,
		DependencyIndexes: file_wordscore_proto_depIdxs,
		EnumInfos:         file_wordscore_proto_enumTypes,
		MessageInfos:      file_wordscore_proto_msgTypes,
	}.Build()
	File_wordscore_proto = out.File
	file_wordscore_proto_rawDesc = nil
	file_wordscore_proto_goTypes = nil
	file_wordscore_proto_depIdxs = nil
}
