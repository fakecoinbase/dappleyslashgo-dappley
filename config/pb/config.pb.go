// Code generated by protoc-gen-go. DO NOT EDIT.
// source: config/pb/config.proto

package configpb

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

type Config struct {
	ConsensusConfig      *ConsensusConfig `protobuf:"bytes,1,opt,name=consensusConfig,proto3" json:"consensusConfig,omitempty"`
	NodeConfig           *NodeConfig      `protobuf:"bytes,2,opt,name=nodeConfig,proto3" json:"nodeConfig,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_2cddfe0b24fd900c, []int{0}
}
func (m *Config) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config.Unmarshal(m, b)
}
func (m *Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config.Marshal(b, m, deterministic)
}
func (dst *Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config.Merge(dst, src)
}
func (m *Config) XXX_Size() int {
	return xxx_messageInfo_Config.Size(m)
}
func (m *Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Config proto.InternalMessageInfo

func (m *Config) GetConsensusConfig() *ConsensusConfig {
	if m != nil {
		return m.ConsensusConfig
	}
	return nil
}

func (m *Config) GetNodeConfig() *NodeConfig {
	if m != nil {
		return m.NodeConfig
	}
	return nil
}

type ConsensusConfig struct {
	MinerAddr            string   `protobuf:"bytes,1,opt,name=minerAddr,proto3" json:"minerAddr,omitempty"`
	PrivKey              string   `protobuf:"bytes,2,opt,name=privKey,proto3" json:"privKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConsensusConfig) Reset()         { *m = ConsensusConfig{} }
func (m *ConsensusConfig) String() string { return proto.CompactTextString(m) }
func (*ConsensusConfig) ProtoMessage()    {}
func (*ConsensusConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_2cddfe0b24fd900c, []int{1}
}
func (m *ConsensusConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConsensusConfig.Unmarshal(m, b)
}
func (m *ConsensusConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConsensusConfig.Marshal(b, m, deterministic)
}
func (dst *ConsensusConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConsensusConfig.Merge(dst, src)
}
func (m *ConsensusConfig) XXX_Size() int {
	return xxx_messageInfo_ConsensusConfig.Size(m)
}
func (m *ConsensusConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_ConsensusConfig.DiscardUnknown(m)
}

var xxx_messageInfo_ConsensusConfig proto.InternalMessageInfo

func (m *ConsensusConfig) GetMinerAddr() string {
	if m != nil {
		return m.MinerAddr
	}
	return ""
}

func (m *ConsensusConfig) GetPrivKey() string {
	if m != nil {
		return m.PrivKey
	}
	return ""
}

type NodeConfig struct {
	Port                 uint32   `protobuf:"varint,1,opt,name=port,proto3" json:"port,omitempty"`
	Seed                 []string `protobuf:"bytes,2,rep,name=seed,proto3" json:"seed,omitempty"`
	DbPath               string   `protobuf:"bytes,3,opt,name=dbPath,proto3" json:"dbPath,omitempty"`
	RpcPort              uint32   `protobuf:"varint,4,opt,name=rpcPort,proto3" json:"rpcPort,omitempty"`
	KeyPath              string   `protobuf:"bytes,5,opt,name=keyPath,proto3" json:"keyPath,omitempty"`
	TxPoolLimit          uint32   `protobuf:"varint,6,opt,name=txPoolLimit,proto3" json:"txPoolLimit,omitempty"`
	NodeAddr             string   `protobuf:"bytes,7,opt,name=nodeAddr,proto3" json:"nodeAddr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeConfig) Reset()         { *m = NodeConfig{} }
func (m *NodeConfig) String() string { return proto.CompactTextString(m) }
func (*NodeConfig) ProtoMessage()    {}
func (*NodeConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_2cddfe0b24fd900c, []int{2}
}
func (m *NodeConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeConfig.Unmarshal(m, b)
}
func (m *NodeConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeConfig.Marshal(b, m, deterministic)
}
func (dst *NodeConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeConfig.Merge(dst, src)
}
func (m *NodeConfig) XXX_Size() int {
	return xxx_messageInfo_NodeConfig.Size(m)
}
func (m *NodeConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeConfig.DiscardUnknown(m)
}

var xxx_messageInfo_NodeConfig proto.InternalMessageInfo

func (m *NodeConfig) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *NodeConfig) GetSeed() []string {
	if m != nil {
		return m.Seed
	}
	return nil
}

func (m *NodeConfig) GetDbPath() string {
	if m != nil {
		return m.DbPath
	}
	return ""
}

func (m *NodeConfig) GetRpcPort() uint32 {
	if m != nil {
		return m.RpcPort
	}
	return 0
}

func (m *NodeConfig) GetKeyPath() string {
	if m != nil {
		return m.KeyPath
	}
	return ""
}

func (m *NodeConfig) GetTxPoolLimit() uint32 {
	if m != nil {
		return m.TxPoolLimit
	}
	return 0
}

func (m *NodeConfig) GetNodeAddr() string {
	if m != nil {
		return m.NodeAddr
	}
	return ""
}

type DynastyConfig struct {
	Producers            []string `protobuf:"bytes,1,rep,name=producers,proto3" json:"producers,omitempty"`
	MaxProducers         uint32   `protobuf:"varint,2,opt,name=maxProducers,proto3" json:"maxProducers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DynastyConfig) Reset()         { *m = DynastyConfig{} }
func (m *DynastyConfig) String() string { return proto.CompactTextString(m) }
func (*DynastyConfig) ProtoMessage()    {}
func (*DynastyConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_2cddfe0b24fd900c, []int{3}
}
func (m *DynastyConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DynastyConfig.Unmarshal(m, b)
}
func (m *DynastyConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DynastyConfig.Marshal(b, m, deterministic)
}
func (dst *DynastyConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DynastyConfig.Merge(dst, src)
}
func (m *DynastyConfig) XXX_Size() int {
	return xxx_messageInfo_DynastyConfig.Size(m)
}
func (m *DynastyConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_DynastyConfig.DiscardUnknown(m)
}

var xxx_messageInfo_DynastyConfig proto.InternalMessageInfo

func (m *DynastyConfig) GetProducers() []string {
	if m != nil {
		return m.Producers
	}
	return nil
}

func (m *DynastyConfig) GetMaxProducers() uint32 {
	if m != nil {
		return m.MaxProducers
	}
	return 0
}

type CliConfig struct {
	Port                 uint32   `protobuf:"varint,1,opt,name=port,proto3" json:"port,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CliConfig) Reset()         { *m = CliConfig{} }
func (m *CliConfig) String() string { return proto.CompactTextString(m) }
func (*CliConfig) ProtoMessage()    {}
func (*CliConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_2cddfe0b24fd900c, []int{4}
}
func (m *CliConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CliConfig.Unmarshal(m, b)
}
func (m *CliConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CliConfig.Marshal(b, m, deterministic)
}
func (dst *CliConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CliConfig.Merge(dst, src)
}
func (m *CliConfig) XXX_Size() int {
	return xxx_messageInfo_CliConfig.Size(m)
}
func (m *CliConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_CliConfig.DiscardUnknown(m)
}

var xxx_messageInfo_CliConfig proto.InternalMessageInfo

func (m *CliConfig) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *CliConfig) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func init() {
	proto.RegisterType((*Config)(nil), "configpb.Config")
	proto.RegisterType((*ConsensusConfig)(nil), "configpb.ConsensusConfig")
	proto.RegisterType((*NodeConfig)(nil), "configpb.NodeConfig")
	proto.RegisterType((*DynastyConfig)(nil), "configpb.DynastyConfig")
	proto.RegisterType((*CliConfig)(nil), "configpb.CliConfig")
}

func init() { proto.RegisterFile("config/pb/config.proto", fileDescriptor_config_2cddfe0b24fd900c) }

var fileDescriptor_config_2cddfe0b24fd900c = []byte{
	// 320 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x52, 0x3d, 0x4f, 0xc3, 0x30,
	0x10, 0x55, 0xda, 0x92, 0x36, 0x57, 0xaa, 0x4a, 0x16, 0xaa, 0x42, 0xc5, 0x50, 0x65, 0xea, 0xd4,
	0x4a, 0xc0, 0xc6, 0x84, 0xc2, 0x82, 0x40, 0x28, 0xf8, 0x1f, 0x24, 0xb1, 0x01, 0x8b, 0xd6, 0xb6,
	0x6c, 0x17, 0x9a, 0x99, 0xff, 0xc5, 0x6f, 0x43, 0xb9, 0x7c, 0xb5, 0x1d, 0xd8, 0xee, 0x3d, 0xdf,
	0x7b, 0xef, 0xee, 0x64, 0x98, 0xe5, 0x4a, 0xbe, 0x89, 0xf7, 0xb5, 0xce, 0xd6, 0x55, 0xb5, 0xd2,
	0x46, 0x39, 0x45, 0x46, 0x15, 0xd2, 0x59, 0xf4, 0xe3, 0x81, 0x1f, 0x23, 0x20, 0x31, 0x4c, 0x73,
	0x25, 0x2d, 0x97, 0x76, 0x67, 0x2b, 0x2a, 0xf4, 0x16, 0xde, 0x72, 0x7c, 0x7d, 0xb9, 0x6a, 0xda,
	0x57, 0xf1, 0x71, 0x03, 0x3d, 0x55, 0x90, 0x5b, 0x00, 0xa9, 0x18, 0xaf, 0xf5, 0x3d, 0xd4, 0x5f,
	0x74, 0xfa, 0x97, 0xf6, 0x8d, 0x1e, 0xf4, 0x45, 0x8f, 0x30, 0x3d, 0x71, 0x26, 0x57, 0x10, 0x6c,
	0x85, 0xe4, 0xe6, 0x9e, 0x31, 0x83, 0x73, 0x04, 0xb4, 0x23, 0x48, 0x08, 0x43, 0x6d, 0xc4, 0xd7,
	0x13, 0x2f, 0x30, 0x23, 0xa0, 0x0d, 0x8c, 0x7e, 0x3d, 0x80, 0x2e, 0x85, 0x10, 0x18, 0x68, 0x65,
	0x1c, 0x3a, 0x4c, 0x28, 0xd6, 0x25, 0x67, 0x39, 0x67, 0x61, 0x6f, 0xd1, 0x5f, 0x06, 0x14, 0x6b,
	0x32, 0x03, 0x9f, 0x65, 0x49, 0xea, 0x3e, 0xc2, 0x3e, 0xfa, 0xd5, 0xa8, 0x0c, 0x32, 0x3a, 0x4f,
	0x4a, 0x8b, 0x01, 0x5a, 0x34, 0xb0, 0x7c, 0xf9, 0xe4, 0x05, 0x4a, 0xce, 0xaa, 0x11, 0x6a, 0x48,
	0x16, 0x30, 0x76, 0xfb, 0x44, 0xa9, 0xcd, 0xb3, 0xd8, 0x0a, 0x17, 0xfa, 0xa8, 0x3b, 0xa4, 0xc8,
	0x1c, 0x46, 0xe5, 0xf6, 0xb8, 0xdb, 0x10, 0xc5, 0x2d, 0x8e, 0x5e, 0x61, 0xf2, 0x50, 0xc8, 0xd4,
	0xba, 0xa2, 0xbb, 0x84, 0x36, 0x8a, 0xed, 0x72, 0x6e, 0x6c, 0xe8, 0xe1, 0xcc, 0x1d, 0x41, 0x22,
	0x38, 0xdf, 0xa6, 0xfb, 0xa4, 0x6d, 0xe8, 0x61, 0xda, 0x11, 0x17, 0xdd, 0x41, 0x10, 0x6f, 0xc4,
	0x3f, 0x17, 0x99, 0xc3, 0x48, 0xa7, 0xd6, 0x7e, 0x2b, 0xc3, 0xea, 0x7b, 0xb6, 0x38, 0xf3, 0xf1,
	0xcb, 0xdc, 0xfc, 0x05, 0x00, 0x00, 0xff, 0xff, 0x7c, 0x85, 0xf9, 0x71, 0x4c, 0x02, 0x00, 0x00,
}
