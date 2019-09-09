// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/dappley/go-dappley/core/transaction/pb/transaction.proto

package transactionpb

import (
	fmt "fmt"
	pb "github.com/dappley/go-dappley/core/transactionbase/pb"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Transactions struct {
	Transactions         []*Transaction `protobuf:"bytes,1,rep,name=transactions,proto3" json:"transactions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Transactions) Reset()         { *m = Transactions{} }
func (m *Transactions) String() string { return proto.CompactTextString(m) }
func (*Transactions) ProtoMessage()    {}
func (*Transactions) Descriptor() ([]byte, []int) {
	return fileDescriptor_4138a4cf34c3b76a, []int{0}
}

func (m *Transactions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transactions.Unmarshal(m, b)
}
func (m *Transactions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transactions.Marshal(b, m, deterministic)
}
func (m *Transactions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transactions.Merge(m, src)
}
func (m *Transactions) XXX_Size() int {
	return xxx_messageInfo_Transactions.Size(m)
}
func (m *Transactions) XXX_DiscardUnknown() {
	xxx_messageInfo_Transactions.DiscardUnknown(m)
}

var xxx_messageInfo_Transactions proto.InternalMessageInfo

func (m *Transactions) GetTransactions() []*Transaction {
	if m != nil {
		return m.Transactions
	}
	return nil
}

type Transaction struct {
	Id                   []byte         `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Vin                  []*pb.TXInput  `protobuf:"bytes,2,rep,name=vin,proto3" json:"vin,omitempty"`
	Vout                 []*pb.TXOutput `protobuf:"bytes,3,rep,name=vout,proto3" json:"vout,omitempty"`
	Tip                  []byte         `protobuf:"bytes,4,opt,name=tip,proto3" json:"tip,omitempty"`
	GasLimit             []byte         `protobuf:"bytes,5,opt,name=gas_limit,json=gasLimit,proto3" json:"gas_limit,omitempty"`
	GasPrice             []byte         `protobuf:"bytes,6,opt,name=gas_price,json=gasPrice,proto3" json:"gas_price,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Transaction) Reset()         { *m = Transaction{} }
func (m *Transaction) String() string { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()    {}
func (*Transaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_4138a4cf34c3b76a, []int{1}
}

func (m *Transaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transaction.Unmarshal(m, b)
}
func (m *Transaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transaction.Marshal(b, m, deterministic)
}
func (m *Transaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transaction.Merge(m, src)
}
func (m *Transaction) XXX_Size() int {
	return xxx_messageInfo_Transaction.Size(m)
}
func (m *Transaction) XXX_DiscardUnknown() {
	xxx_messageInfo_Transaction.DiscardUnknown(m)
}

var xxx_messageInfo_Transaction proto.InternalMessageInfo

func (m *Transaction) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *Transaction) GetVin() []*pb.TXInput {
	if m != nil {
		return m.Vin
	}
	return nil
}

func (m *Transaction) GetVout() []*pb.TXOutput {
	if m != nil {
		return m.Vout
	}
	return nil
}

func (m *Transaction) GetTip() []byte {
	if m != nil {
		return m.Tip
	}
	return nil
}

func (m *Transaction) GetGasLimit() []byte {
	if m != nil {
		return m.GasLimit
	}
	return nil
}

func (m *Transaction) GetGasPrice() []byte {
	if m != nil {
		return m.GasPrice
	}
	return nil
}

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
	return fileDescriptor_4138a4cf34c3b76a, []int{2}
}

func (m *TransactionNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransactionNode.Unmarshal(m, b)
}
func (m *TransactionNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransactionNode.Marshal(b, m, deterministic)
}
func (m *TransactionNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransactionNode.Merge(m, src)
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

type TransactionJournal struct {
	Vout                 []*pb.TXOutput `protobuf:"bytes,1,rep,name=vout,proto3" json:"vout,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *TransactionJournal) Reset()         { *m = TransactionJournal{} }
func (m *TransactionJournal) String() string { return proto.CompactTextString(m) }
func (*TransactionJournal) ProtoMessage()    {}
func (*TransactionJournal) Descriptor() ([]byte, []int) {
	return fileDescriptor_4138a4cf34c3b76a, []int{3}
}

func (m *TransactionJournal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransactionJournal.Unmarshal(m, b)
}
func (m *TransactionJournal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransactionJournal.Marshal(b, m, deterministic)
}
func (m *TransactionJournal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransactionJournal.Merge(m, src)
}
func (m *TransactionJournal) XXX_Size() int {
	return xxx_messageInfo_TransactionJournal.Size(m)
}
func (m *TransactionJournal) XXX_DiscardUnknown() {
	xxx_messageInfo_TransactionJournal.DiscardUnknown(m)
}

var xxx_messageInfo_TransactionJournal proto.InternalMessageInfo

func (m *TransactionJournal) GetVout() []*pb.TXOutput {
	if m != nil {
		return m.Vout
	}
	return nil
}

func init() {
	proto.RegisterType((*Transactions)(nil), "transactionpb.Transactions")
	proto.RegisterType((*Transaction)(nil), "transactionpb.Transaction")
	proto.RegisterType((*TransactionNode)(nil), "transactionpb.TransactionNode")
	proto.RegisterMapType((map[string]*Transaction)(nil), "transactionpb.TransactionNode.ChildrenEntry")
	proto.RegisterType((*TransactionJournal)(nil), "transactionpb.TransactionJournal")
}

func init() {
	proto.RegisterFile("github.com/dappley/go-dappley/core/transaction/pb/transaction.proto", fileDescriptor_4138a4cf34c3b76a)
}

var fileDescriptor_4138a4cf34c3b76a = []byte{
	// 364 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0x4d, 0x4f, 0xf2, 0x40,
	0x14, 0x85, 0x33, 0x2d, 0x10, 0xb8, 0xc0, 0xfb, 0x9a, 0x59, 0x4d, 0x60, 0xd3, 0x74, 0xc5, 0x02,
	0x5b, 0xa3, 0x1b, 0xe3, 0xc2, 0x85, 0x84, 0x44, 0x0d, 0x41, 0xd3, 0x98, 0xe8, 0x8e, 0x4c, 0xdb,
	0x49, 0x99, 0x58, 0x3a, 0x93, 0xce, 0x94, 0x04, 0x7f, 0xa1, 0xbf, 0xc9, 0x95, 0x99, 0xf2, 0x61,
	0xf1, 0x23, 0xca, 0xee, 0xde, 0x39, 0xcf, 0x3d, 0x77, 0x7a, 0xa6, 0x30, 0x4a, 0xb8, 0x9e, 0x17,
	0xa1, 0x17, 0x89, 0x85, 0x1f, 0x53, 0x29, 0x53, 0xb6, 0xf2, 0x13, 0x71, 0xbc, 0x2d, 0x23, 0x91,
	0x33, 0x5f, 0xe7, 0x34, 0x53, 0x34, 0xd2, 0x5c, 0x64, 0xbe, 0x0c, 0xab, 0xad, 0x27, 0x73, 0xa1,
	0x05, 0xee, 0x56, 0x8e, 0x64, 0xd8, 0x9b, 0x1c, 0xe6, 0x39, 0x0b, 0xa9, 0x62, 0x9f, 0x8c, 0xaf,
	0xa8, 0x62, 0x6b, 0x73, 0x77, 0x0a, 0x9d, 0x87, 0x0f, 0x41, 0xe1, 0x4b, 0xe8, 0x54, 0x40, 0x45,
	0x90, 0x63, 0x0f, 0xda, 0xa7, 0x3d, 0x6f, 0xef, 0x0e, 0x5e, 0x65, 0x24, 0xd8, 0xe3, 0xdd, 0x57,
	0x04, 0xed, 0x8a, 0x8a, 0xff, 0x81, 0xc5, 0x63, 0x82, 0x1c, 0x34, 0xe8, 0x04, 0x16, 0x8f, 0xf1,
	0x10, 0xec, 0x25, 0xcf, 0x88, 0xf5, 0xd5, 0xd6, 0xdc, 0xd3, 0x58, 0x3f, 0xdd, 0x64, 0xb2, 0xd0,
	0x81, 0xc1, 0xb0, 0x0f, 0xb5, 0xa5, 0x28, 0x34, 0xb1, 0x4b, 0xbc, 0xff, 0x2d, 0x7e, 0x57, 0x68,
	0xc3, 0x97, 0x20, 0x3e, 0x02, 0x5b, 0x73, 0x49, 0x6a, 0xe5, 0x3e, 0x53, 0xe2, 0x3e, 0xb4, 0x12,
	0xaa, 0x66, 0x29, 0x5f, 0x70, 0x4d, 0xea, 0xe5, 0x79, 0x33, 0xa1, 0x6a, 0x62, 0xfa, 0xad, 0x28,
	0x73, 0x1e, 0x31, 0xd2, 0xd8, 0x89, 0xf7, 0xa6, 0x77, 0xdf, 0x10, 0xfc, 0xaf, 0x7c, 0xca, 0x54,
	0xc4, 0x0c, 0x5f, 0x43, 0x33, 0x9a, 0xf3, 0x34, 0xce, 0x59, 0xb6, 0x89, 0x66, 0xf8, 0x73, 0x34,
	0x66, 0xc2, 0x1b, 0x6d, 0xf0, 0x71, 0xa6, 0xf3, 0x55, 0xb0, 0x9b, 0xc6, 0x27, 0x50, 0x5f, 0xd2,
	0xb4, 0x60, 0xc4, 0x72, 0xd0, 0x2f, 0x09, 0xaf, 0x41, 0x8c, 0xa1, 0xa6, 0xf8, 0x0b, 0x23, 0xb6,
	0x83, 0x06, 0x76, 0x50, 0xd6, 0xbd, 0x47, 0xe8, 0xee, 0x2d, 0x30, 0x01, 0x3c, 0xb3, 0x55, 0x19,
	0x78, 0x2b, 0x30, 0xe5, 0xe1, 0x8b, 0x2e, 0xac, 0x73, 0xe4, 0x8e, 0x01, 0x57, 0x94, 0x5b, 0x51,
	0xe4, 0x19, 0x4d, 0x77, 0xef, 0x81, 0xfe, 0xf8, 0x1e, 0x61, 0xa3, 0xfc, 0xcb, 0xce, 0xde, 0x03,
	0x00, 0x00, 0xff, 0xff, 0x03, 0x82, 0x44, 0x1c, 0x09, 0x03, 0x00, 0x00,
}
