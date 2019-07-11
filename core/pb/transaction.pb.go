// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/dappley/go-dappley/core/pb/transaction.proto

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
	return fileDescriptor_transaction_ed2615d3a1cee5a5, []int{0}
}
func (m *Transactions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transactions.Unmarshal(m, b)
}
func (m *Transactions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transactions.Marshal(b, m, deterministic)
}
func (dst *Transactions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transactions.Merge(dst, src)
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
	Id                   []byte      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Vin                  []*TXInput  `protobuf:"bytes,2,rep,name=vin,proto3" json:"vin,omitempty"`
	Vout                 []*TXOutput `protobuf:"bytes,3,rep,name=vout,proto3" json:"vout,omitempty"`
	Tip                  []byte      `protobuf:"bytes,4,opt,name=tip,proto3" json:"tip,omitempty"`
	GasLimit             []byte      `protobuf:"bytes,5,opt,name=gasLimit,proto3" json:"gasLimit,omitempty"`
	GasPrice             []byte      `protobuf:"bytes,6,opt,name=gasPrice,proto3" json:"gasPrice,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Transaction) Reset()         { *m = Transaction{} }
func (m *Transaction) String() string { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()    {}
func (*Transaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_ed2615d3a1cee5a5, []int{1}
}
func (m *Transaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transaction.Unmarshal(m, b)
}
func (m *Transaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transaction.Marshal(b, m, deterministic)
}
func (dst *Transaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transaction.Merge(dst, src)
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

func (m *Transaction) GetVin() []*TXInput {
	if m != nil {
		return m.Vin
	}
	return nil
}

func (m *Transaction) GetVout() []*TXOutput {
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

type TXInput struct {
	Txid                 []byte   `protobuf:"bytes,1,opt,name=txid,proto3" json:"txid,omitempty"`
	Vout                 int32    `protobuf:"varint,2,opt,name=vout,proto3" json:"vout,omitempty"`
	Signature            []byte   `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	PublicKey            []byte   `protobuf:"bytes,4,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TXInput) Reset()         { *m = TXInput{} }
func (m *TXInput) String() string { return proto.CompactTextString(m) }
func (*TXInput) ProtoMessage()    {}
func (*TXInput) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_ed2615d3a1cee5a5, []int{2}
}
func (m *TXInput) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TXInput.Unmarshal(m, b)
}
func (m *TXInput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TXInput.Marshal(b, m, deterministic)
}
func (dst *TXInput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TXInput.Merge(dst, src)
}
func (m *TXInput) XXX_Size() int {
	return xxx_messageInfo_TXInput.Size(m)
}
func (m *TXInput) XXX_DiscardUnknown() {
	xxx_messageInfo_TXInput.DiscardUnknown(m)
}

var xxx_messageInfo_TXInput proto.InternalMessageInfo

func (m *TXInput) GetTxid() []byte {
	if m != nil {
		return m.Txid
	}
	return nil
}

func (m *TXInput) GetVout() int32 {
	if m != nil {
		return m.Vout
	}
	return 0
}

func (m *TXInput) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *TXInput) GetPublicKey() []byte {
	if m != nil {
		return m.PublicKey
	}
	return nil
}

type TXOutput struct {
	Value                []byte   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	PublicKeyHash        []byte   `protobuf:"bytes,2,opt,name=public_key_hash,json=publicKeyHash,proto3" json:"public_key_hash,omitempty"`
	Contract             string   `protobuf:"bytes,3,opt,name=contract,proto3" json:"contract,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TXOutput) Reset()         { *m = TXOutput{} }
func (m *TXOutput) String() string { return proto.CompactTextString(m) }
func (*TXOutput) ProtoMessage()    {}
func (*TXOutput) Descriptor() ([]byte, []int) {
	return fileDescriptor_transaction_ed2615d3a1cee5a5, []int{3}
}
func (m *TXOutput) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TXOutput.Unmarshal(m, b)
}
func (m *TXOutput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TXOutput.Marshal(b, m, deterministic)
}
func (dst *TXOutput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TXOutput.Merge(dst, src)
}
func (m *TXOutput) XXX_Size() int {
	return xxx_messageInfo_TXOutput.Size(m)
}
func (m *TXOutput) XXX_DiscardUnknown() {
	xxx_messageInfo_TXOutput.DiscardUnknown(m)
}

var xxx_messageInfo_TXOutput proto.InternalMessageInfo

func (m *TXOutput) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *TXOutput) GetPublicKeyHash() []byte {
	if m != nil {
		return m.PublicKeyHash
	}
	return nil
}

func (m *TXOutput) GetContract() string {
	if m != nil {
		return m.Contract
	}
	return ""
}

func init() {
	proto.RegisterType((*Transactions)(nil), "corepb.Transactions")
	proto.RegisterType((*Transaction)(nil), "corepb.Transaction")
	proto.RegisterType((*TXInput)(nil), "corepb.TXInput")
	proto.RegisterType((*TXOutput)(nil), "corepb.TXOutput")
}

func init() {
	proto.RegisterFile("github.com/dappley/go-dappley/core/pb/transaction.proto", fileDescriptor_transaction_ed2615d3a1cee5a5)
}

var fileDescriptor_transaction_ed2615d3a1cee5a5 = []byte{
	// 329 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x91, 0x51, 0x4b, 0xc3, 0x30,
	0x14, 0x85, 0x69, 0xbb, 0xcd, 0xed, 0xae, 0xba, 0x11, 0x7d, 0x08, 0xa2, 0x30, 0x8b, 0xc8, 0x5e,
	0x6c, 0x41, 0x1f, 0xf6, 0x13, 0x54, 0x14, 0x94, 0xe2, 0x83, 0x6f, 0x23, 0x4d, 0x43, 0x1b, 0xdc,
	0x92, 0xd0, 0x24, 0xc3, 0xfd, 0x2b, 0x7f, 0xa2, 0x34, 0xeb, 0xda, 0xfa, 0x76, 0xef, 0x3d, 0x1f,
	0xe7, 0x9c, 0x10, 0x58, 0x15, 0xdc, 0x94, 0x36, 0x8b, 0xa9, 0xdc, 0x26, 0x39, 0x51, 0x6a, 0xc3,
	0xf6, 0x49, 0x21, 0xef, 0x8f, 0x23, 0x95, 0x15, 0x4b, 0x54, 0x96, 0x98, 0x8a, 0x08, 0x4d, 0xa8,
	0xe1, 0x52, 0xc4, 0xaa, 0x92, 0x46, 0xa2, 0x51, 0x2d, 0xa9, 0x2c, 0x7a, 0x82, 0xf0, 0xb3, 0x13,
	0x35, 0x5a, 0x41, 0xd8, 0x83, 0x35, 0xf6, 0x16, 0xc1, 0x72, 0xfa, 0x70, 0x1e, 0x1f, 0xf0, 0xb8,
	0xc7, 0xa6, 0xff, 0xc0, 0xe8, 0xd7, 0x83, 0x69, 0x4f, 0x45, 0x67, 0xe0, 0xf3, 0x1c, 0x7b, 0x0b,
	0x6f, 0x19, 0xa6, 0x3e, 0xcf, 0xd1, 0x0d, 0x04, 0x3b, 0x2e, 0xb0, 0xef, 0xfc, 0x66, 0xad, 0xdf,
	0xd7, 0x8b, 0x50, 0xd6, 0xa4, 0xb5, 0x86, 0x6e, 0x61, 0xb0, 0x93, 0xd6, 0xe0, 0xc0, 0x31, 0xf3,
	0x8e, 0x79, 0xb7, 0xa6, 0x86, 0x9c, 0x8a, 0xe6, 0x10, 0x18, 0xae, 0xf0, 0xc0, 0x39, 0xd7, 0x23,
	0xba, 0x84, 0x71, 0x41, 0xf4, 0x1b, 0xdf, 0x72, 0x83, 0x87, 0xee, 0xdc, 0xee, 0x8d, 0xf6, 0x51,
	0x71, 0xca, 0xf0, 0xa8, 0xd5, 0xdc, 0x1e, 0x09, 0x38, 0x69, 0xf2, 0x11, 0x82, 0x81, 0xf9, 0x69,
	0xfb, 0xba, 0xb9, 0xbe, 0xb9, 0x3a, 0xfe, 0xc2, 0x5b, 0x0e, 0x9b, 0xf0, 0x2b, 0x98, 0x68, 0x5e,
	0x08, 0x62, 0x6c, 0xc5, 0x70, 0xe0, 0xe0, 0xee, 0x80, 0xae, 0x01, 0x94, 0xcd, 0x36, 0x9c, 0xae,
	0xbf, 0xd9, 0xbe, 0x69, 0x38, 0x39, 0x5c, 0x5e, 0xd9, 0x3e, 0xca, 0x61, 0x7c, 0x7c, 0x0b, 0xba,
	0x80, 0xe1, 0x8e, 0x6c, 0x2c, 0x6b, 0x12, 0x0f, 0x0b, 0xba, 0x83, 0x59, 0x67, 0xb0, 0x2e, 0x89,
	0x2e, 0x5d, 0x7a, 0x98, 0x9e, 0xb6, 0x2e, 0xcf, 0x44, 0x97, 0xf5, 0xab, 0xa8, 0x14, 0xa6, 0x22,
	0xd4, 0xb8, 0x16, 0x93, 0xb4, 0xdd, 0xb3, 0x91, 0xfb, 0xe0, 0xc7, 0xbf, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xf3, 0x9c, 0xaf, 0xc0, 0x1b, 0x02, 0x00, 0x00,
}
