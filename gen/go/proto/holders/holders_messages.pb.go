// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/holders/holders_messages.proto

package tutorial_grpc_holders_v1

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Customer definition.
type Holder struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	FirstName            string   `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName             string   `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	Phone                string   `protobuf:"bytes,4,opt,name=phone,proto3" json:"phone,omitempty"`
	Email                string   `protobuf:"bytes,5,opt,name=email,proto3" json:"email,omitempty"`
	HeldBooks            []string `protobuf:"bytes,6,rep,name=held_books,json=heldBooks,proto3" json:"held_books,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Holder) Reset()         { *m = Holder{} }
func (m *Holder) String() string { return proto.CompactTextString(m) }
func (*Holder) ProtoMessage()    {}
func (*Holder) Descriptor() ([]byte, []int) {
	return fileDescriptor_087b08adb09aa2ea, []int{0}
}

func (m *Holder) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Holder.Unmarshal(m, b)
}
func (m *Holder) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Holder.Marshal(b, m, deterministic)
}
func (m *Holder) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Holder.Merge(m, src)
}
func (m *Holder) XXX_Size() int {
	return xxx_messageInfo_Holder.Size(m)
}
func (m *Holder) XXX_DiscardUnknown() {
	xxx_messageInfo_Holder.DiscardUnknown(m)
}

var xxx_messageInfo_Holder proto.InternalMessageInfo

func (m *Holder) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Holder) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *Holder) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *Holder) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *Holder) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *Holder) GetHeldBooks() []string {
	if m != nil {
		return m.HeldBooks
	}
	return nil
}

func init() {
	proto.RegisterType((*Holder)(nil), "tutorial.grpc.holders.v1.Holder")
}

func init() {
	proto.RegisterFile("proto/holders/holders_messages.proto", fileDescriptor_087b08adb09aa2ea)
}

var fileDescriptor_087b08adb09aa2ea = []byte{
	// 190 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x34, 0x8f, 0xb1, 0xaa, 0xc2, 0x30,
	0x14, 0x86, 0x69, 0x7b, 0x5b, 0x6e, 0x33, 0xdc, 0x21, 0xdc, 0x21, 0x20, 0x42, 0x11, 0x07, 0xa7,
	0x88, 0xf8, 0x06, 0x4e, 0x4e, 0x0e, 0xbe, 0x40, 0x49, 0xcd, 0xb1, 0x0d, 0x26, 0x3d, 0x25, 0x89,
	0xbe, 0x8e, 0xaf, 0x2a, 0x39, 0xb5, 0x53, 0xf8, 0xbe, 0x2f, 0xf0, 0x73, 0xd8, 0x76, 0xf2, 0x18,
	0x71, 0x3f, 0xa0, 0xd5, 0xe0, 0xc3, 0xf2, 0xb6, 0x0e, 0x42, 0x50, 0x3d, 0x04, 0x49, 0x99, 0x8b,
	0xf8, 0x8c, 0xe8, 0x8d, 0xb2, 0xb2, 0xf7, 0xd3, 0x4d, 0x7e, 0x7f, 0xc9, 0xd7, 0x61, 0xf3, 0xce,
	0x58, 0x75, 0x26, 0xe4, 0x7f, 0x2c, 0x37, 0x5a, 0x64, 0x4d, 0xb6, 0xab, 0xaf, 0xb9, 0xd1, 0x7c,
	0xcd, 0xd8, 0xdd, 0xf8, 0x10, 0xdb, 0x51, 0x39, 0x10, 0x39, 0xf9, 0x9a, 0xcc, 0x45, 0x39, 0xe0,
	0x2b, 0x56, 0x5b, 0xb5, 0xd4, 0x82, 0xea, 0x6f, 0x12, 0x14, 0xff, 0x59, 0x39, 0x0d, 0x38, 0x82,
	0xf8, 0xa1, 0x30, 0x43, 0xb2, 0xe0, 0x94, 0xb1, 0xa2, 0x9c, 0x2d, 0x41, 0xda, 0x19, 0xc0, 0xea,
	0xb6, 0x43, 0x7c, 0x04, 0x51, 0x35, 0x45, 0xda, 0x49, 0xe6, 0x94, 0x44, 0x57, 0xd1, 0x09, 0xc7,
	0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xd7, 0xb6, 0x3f, 0x8c, 0xea, 0x00, 0x00, 0x00,
}
