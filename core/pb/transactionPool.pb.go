// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/dappley/go-dappley/core/pb/transactionPool.proto

package corepb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type TransactionNode struct {
	Children             map[string]*Transaction `protobuf:"bytes,1,rep,name=children,proto3" json:"children,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Value                *Transaction            `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Size                 int64                   `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *TransactionNode) Reset()         { *m = TransactionNode{} }
func (m *TransactionNode) String() string { return proto.CompactTextString(m) }
func (*TransactionNode) ProtoMessage()    {}
func (*TransactionNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_transactionPool_165430fdffed3bde, []int{0}
}
func (m *TransactionNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransactionNode.Unmarshal(m, b)
}
func (m *TransactionNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransactionNode.Marshal(b, m, deterministic)
}
func (dst *TransactionNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransactionNode.Merge(dst, src)
}
func (m *TransactionNode) XXX_Size() int {
	return xxx_messageInfo_TransactionNode.Size(m)
}
func (m *TransactionNode) XXX_DiscardUnknown() {
	xxx_messageInfo_TransactionNode.DiscardUnknown(m)
}

var xxx_messageInfo_TransactionNode proto.InternalMessageInfo

func (m *TransactionNode) GetChildren() map[string]*Transaction {
	if m != nil {
		return m.Children
	}
	return nil
}

func (m *TransactionNode) GetValue() *Transaction {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *TransactionNode) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

type TransactionPool struct {
	Txs                  map[string]*TransactionNode `protobuf:"bytes,1,rep,name=txs,proto3" json:"txs,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	TipOrder             []string                    `protobuf:"bytes,2,rep,name=tipOrder,proto3" json:"tipOrder,omitempty"`
	CurrSize             uint32                      `protobuf:"varint,3,opt,name=currSize,proto3" json:"currSize,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *TransactionPool) Reset()         { *m = TransactionPool{} }
func (m *TransactionPool) String() string { return proto.CompactTextString(m) }
func (*TransactionPool) ProtoMessage()    {}
func (*TransactionPool) Descriptor() ([]byte, []int) {
	return fileDescriptor_transactionPool_165430fdffed3bde, []int{1}
}
func (m *TransactionPool) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransactionPool.Unmarshal(m, b)
}
func (m *TransactionPool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransactionPool.Marshal(b, m, deterministic)
}
func (dst *TransactionPool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransactionPool.Merge(dst, src)
}
func (m *TransactionPool) XXX_Size() int {
	return xxx_messageInfo_TransactionPool.Size(m)
}
func (m *TransactionPool) XXX_DiscardUnknown() {
	xxx_messageInfo_TransactionPool.DiscardUnknown(m)
}

var xxx_messageInfo_TransactionPool proto.InternalMessageInfo

func (m *TransactionPool) GetTxs() map[string]*TransactionNode {
	if m != nil {
		return m.Txs
	}
	return nil
}

func (m *TransactionPool) GetTipOrder() []string {
	if m != nil {
		return m.TipOrder
	}
	return nil
}

func (m *TransactionPool) GetCurrSize() uint32 {
	if m != nil {
		return m.CurrSize
	}
	return 0
}

func init() {
	proto.RegisterType((*TransactionNode)(nil), "corepb.TransactionNode")
	proto.RegisterMapType((map[string]*Transaction)(nil), "corepb.TransactionNode.ChildrenEntry")
	proto.RegisterType((*TransactionPool)(nil), "corepb.TransactionPool")
	proto.RegisterMapType((map[string]*TransactionNode)(nil), "corepb.TransactionPool.TxsEntry")
}

func init() {
	proto.RegisterFile("github.com/dappley/go-dappley/core/pb/transactionPool.proto", fileDescriptor_transactionPool_165430fdffed3bde)
}

var fileDescriptor_transactionPool_165430fdffed3bde = []byte{
	// 294 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xcf, 0x4a, 0xf3, 0x40,
	0x14, 0xc5, 0x99, 0xcc, 0xf7, 0x95, 0xf4, 0x96, 0xa2, 0x8c, 0x0b, 0x43, 0x56, 0x43, 0x41, 0x88,
	0x8b, 0x26, 0x10, 0x17, 0x8a, 0xae, 0x44, 0xdc, 0xda, 0x32, 0xf6, 0x05, 0xf2, 0x67, 0x68, 0x83,
	0x31, 0x33, 0x4c, 0x26, 0xd2, 0xf8, 0x90, 0xbe, 0x85, 0xef, 0x21, 0x93, 0x34, 0xd1, 0x48, 0x5c,
	0x74, 0x77, 0xef, 0xdc, 0x73, 0xcf, 0xe5, 0x77, 0x06, 0xee, 0xb6, 0x99, 0xde, 0x55, 0xb1, 0x9f,
	0x88, 0xd7, 0x20, 0x8d, 0xa4, 0xcc, 0x79, 0x1d, 0x6c, 0xc5, 0xb2, 0x2b, 0x13, 0xa1, 0x78, 0x20,
	0xe3, 0x40, 0xab, 0xa8, 0x28, 0xa3, 0x44, 0x67, 0xa2, 0x58, 0x0b, 0x91, 0xfb, 0x52, 0x09, 0x2d,
	0xc8, 0xc4, 0x8c, 0x65, 0xec, 0x5e, 0x1f, 0x6d, 0xd2, 0x1a, 0x2c, 0x3e, 0x11, 0x9c, 0x6c, 0xbe,
	0x5f, 0x9f, 0x44, 0xca, 0xc9, 0x3d, 0xd8, 0xc9, 0x2e, 0xcb, 0x53, 0xc5, 0x0b, 0x07, 0x51, 0xec,
	0xcd, 0xc2, 0x0b, 0xbf, 0xbd, 0xe3, 0xff, 0x92, 0xfa, 0x0f, 0x07, 0xdd, 0x63, 0xa1, 0x55, 0xcd,
	0xfa, 0x35, 0x72, 0x09, 0xff, 0xdf, 0xa2, 0xbc, 0xe2, 0x8e, 0x45, 0x91, 0x37, 0x0b, 0xcf, 0x46,
	0xf6, 0x59, 0xab, 0x20, 0x04, 0xfe, 0x95, 0xd9, 0x3b, 0x77, 0x30, 0x45, 0x1e, 0x66, 0x4d, 0xed,
	0xae, 0x61, 0x3e, 0x70, 0x26, 0xa7, 0x80, 0x5f, 0x78, 0xed, 0x20, 0x8a, 0xbc, 0x29, 0x33, 0xe5,
	0x11, 0x17, 0x6e, 0xad, 0x1b, 0xb4, 0xf8, 0x18, 0x72, 0x9a, 0x08, 0x49, 0x08, 0x58, 0xef, 0xcb,
	0x03, 0x22, 0x1d, 0x31, 0x68, 0x82, 0xde, 0xec, 0xcb, 0x96, 0xce, 0x88, 0x89, 0x0b, 0xb6, 0xce,
	0xe4, 0x4a, 0xa5, 0x5c, 0x39, 0x16, 0xc5, 0xde, 0x94, 0xf5, 0xbd, 0x99, 0x25, 0x95, 0x52, 0xcf,
	0x1d, 0xcd, 0x9c, 0xf5, 0xbd, 0xbb, 0x02, 0xbb, 0x33, 0x1a, 0x81, 0x59, 0x0e, 0x61, 0xce, 0xff,
	0x88, 0xfb, 0x07, 0x50, 0x3c, 0x69, 0xfe, 0xef, 0xea, 0x2b, 0x00, 0x00, 0xff, 0xff, 0xc0, 0x10,
	0x1c, 0xd0, 0x3f, 0x02, 0x00, 0x00,
}