// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: install.proto

package v1 // import "github.com/solo-io/supergloo/pkg/api/v1"

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import core "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

import bytes "bytes"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

//
// @solo-kit:resource.short_name=install
// @solo-kit:resource.plural_name=installs
// @solo-kit:resource.resource_groups=install.supergloo.solo.io
type Install struct {
	// Status indicates the validation status of this resource.
	// Status is read-only by clients, and set by gloo during validation
	Status core.Status `protobuf:"bytes,1,opt,name=status" json:"status" testdiff:"ignore"`
	// Metadata contains the object metadata for this resource
	Metadata core.Metadata `protobuf:"bytes,2,opt,name=metadata" json:"metadata"`
	// mesh-specific configuration
	//
	// Types that are valid to be assigned to MeshType:
	//	*Install_Istio
	//	*Install_Linkerd2
	//	*Install_Consul
	MeshType     isInstall_MeshType `protobuf_oneof:"mesh_type"`
	ChartLocator *HelmChartLocator  `protobuf:"bytes,6,opt,name=chartLocator" json:"chartLocator,omitempty"`
	Encryption   *Encryption        `protobuf:"bytes,7,opt,name=encryption" json:"encryption,omitempty"`
	// whether or not this install should be enabled
	// if disabled, corresponding resources will be uninstalled
	Enabled              bool     `protobuf:"varint,12,opt,name=enabled,proto3" json:"enabled,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Install) Reset()         { *m = Install{} }
func (m *Install) String() string { return proto.CompactTextString(m) }
func (*Install) ProtoMessage()    {}
func (*Install) Descriptor() ([]byte, []int) {
	return fileDescriptor_install_da8b87227e98e8d3, []int{0}
}
func (m *Install) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Install.Unmarshal(m, b)
}
func (m *Install) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Install.Marshal(b, m, deterministic)
}
func (dst *Install) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Install.Merge(dst, src)
}
func (m *Install) XXX_Size() int {
	return xxx_messageInfo_Install.Size(m)
}
func (m *Install) XXX_DiscardUnknown() {
	xxx_messageInfo_Install.DiscardUnknown(m)
}

var xxx_messageInfo_Install proto.InternalMessageInfo

type isInstall_MeshType interface {
	isInstall_MeshType()
	Equal(interface{}) bool
}

type Install_Istio struct {
	Istio *Istio `protobuf:"bytes,10,opt,name=istio,oneof"`
}
type Install_Linkerd2 struct {
	Linkerd2 *Linkerd2 `protobuf:"bytes,20,opt,name=linkerd2,oneof"`
}
type Install_Consul struct {
	Consul *Consul `protobuf:"bytes,30,opt,name=consul,oneof"`
}

func (*Install_Istio) isInstall_MeshType()    {}
func (*Install_Linkerd2) isInstall_MeshType() {}
func (*Install_Consul) isInstall_MeshType()   {}

func (m *Install) GetMeshType() isInstall_MeshType {
	if m != nil {
		return m.MeshType
	}
	return nil
}

func (m *Install) GetStatus() core.Status {
	if m != nil {
		return m.Status
	}
	return core.Status{}
}

func (m *Install) GetMetadata() core.Metadata {
	if m != nil {
		return m.Metadata
	}
	return core.Metadata{}
}

func (m *Install) GetIstio() *Istio {
	if x, ok := m.GetMeshType().(*Install_Istio); ok {
		return x.Istio
	}
	return nil
}

func (m *Install) GetLinkerd2() *Linkerd2 {
	if x, ok := m.GetMeshType().(*Install_Linkerd2); ok {
		return x.Linkerd2
	}
	return nil
}

func (m *Install) GetConsul() *Consul {
	if x, ok := m.GetMeshType().(*Install_Consul); ok {
		return x.Consul
	}
	return nil
}

func (m *Install) GetChartLocator() *HelmChartLocator {
	if m != nil {
		return m.ChartLocator
	}
	return nil
}

func (m *Install) GetEncryption() *Encryption {
	if m != nil {
		return m.Encryption
	}
	return nil
}

func (m *Install) GetEnabled() bool {
	if m != nil {
		return m.Enabled
	}
	return false
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Install) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Install_OneofMarshaler, _Install_OneofUnmarshaler, _Install_OneofSizer, []interface{}{
		(*Install_Istio)(nil),
		(*Install_Linkerd2)(nil),
		(*Install_Consul)(nil),
	}
}

func _Install_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Install)
	// mesh_type
	switch x := m.MeshType.(type) {
	case *Install_Istio:
		_ = b.EncodeVarint(10<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Istio); err != nil {
			return err
		}
	case *Install_Linkerd2:
		_ = b.EncodeVarint(20<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Linkerd2); err != nil {
			return err
		}
	case *Install_Consul:
		_ = b.EncodeVarint(30<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Consul); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Install.MeshType has unexpected type %T", x)
	}
	return nil
}

func _Install_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Install)
	switch tag {
	case 10: // mesh_type.istio
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Istio)
		err := b.DecodeMessage(msg)
		m.MeshType = &Install_Istio{msg}
		return true, err
	case 20: // mesh_type.linkerd2
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Linkerd2)
		err := b.DecodeMessage(msg)
		m.MeshType = &Install_Linkerd2{msg}
		return true, err
	case 30: // mesh_type.consul
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Consul)
		err := b.DecodeMessage(msg)
		m.MeshType = &Install_Consul{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Install_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Install)
	// mesh_type
	switch x := m.MeshType.(type) {
	case *Install_Istio:
		s := proto.Size(x.Istio)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Install_Linkerd2:
		s := proto.Size(x.Linkerd2)
		n += 2 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Install_Consul:
		s := proto.Size(x.Consul)
		n += 2 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type HelmChartLocator struct {
	// Types that are valid to be assigned to Kind:
	//	*HelmChartLocator_ChartPath
	Kind                 isHelmChartLocator_Kind `protobuf_oneof:"kind"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *HelmChartLocator) Reset()         { *m = HelmChartLocator{} }
func (m *HelmChartLocator) String() string { return proto.CompactTextString(m) }
func (*HelmChartLocator) ProtoMessage()    {}
func (*HelmChartLocator) Descriptor() ([]byte, []int) {
	return fileDescriptor_install_da8b87227e98e8d3, []int{1}
}
func (m *HelmChartLocator) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HelmChartLocator.Unmarshal(m, b)
}
func (m *HelmChartLocator) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HelmChartLocator.Marshal(b, m, deterministic)
}
func (dst *HelmChartLocator) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelmChartLocator.Merge(dst, src)
}
func (m *HelmChartLocator) XXX_Size() int {
	return xxx_messageInfo_HelmChartLocator.Size(m)
}
func (m *HelmChartLocator) XXX_DiscardUnknown() {
	xxx_messageInfo_HelmChartLocator.DiscardUnknown(m)
}

var xxx_messageInfo_HelmChartLocator proto.InternalMessageInfo

type isHelmChartLocator_Kind interface {
	isHelmChartLocator_Kind()
	Equal(interface{}) bool
}

type HelmChartLocator_ChartPath struct {
	ChartPath *HelmChartPath `protobuf:"bytes,1,opt,name=chartPath,oneof"`
}

func (*HelmChartLocator_ChartPath) isHelmChartLocator_Kind() {}

func (m *HelmChartLocator) GetKind() isHelmChartLocator_Kind {
	if m != nil {
		return m.Kind
	}
	return nil
}

func (m *HelmChartLocator) GetChartPath() *HelmChartPath {
	if x, ok := m.GetKind().(*HelmChartLocator_ChartPath); ok {
		return x.ChartPath
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*HelmChartLocator) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _HelmChartLocator_OneofMarshaler, _HelmChartLocator_OneofUnmarshaler, _HelmChartLocator_OneofSizer, []interface{}{
		(*HelmChartLocator_ChartPath)(nil),
	}
}

func _HelmChartLocator_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*HelmChartLocator)
	// kind
	switch x := m.Kind.(type) {
	case *HelmChartLocator_ChartPath:
		_ = b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ChartPath); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("HelmChartLocator.Kind has unexpected type %T", x)
	}
	return nil
}

func _HelmChartLocator_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*HelmChartLocator)
	switch tag {
	case 1: // kind.chartPath
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(HelmChartPath)
		err := b.DecodeMessage(msg)
		m.Kind = &HelmChartLocator_ChartPath{msg}
		return true, err
	default:
		return false, nil
	}
}

func _HelmChartLocator_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*HelmChartLocator)
	// kind
	switch x := m.Kind.(type) {
	case *HelmChartLocator_ChartPath:
		s := proto.Size(x.ChartPath)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type HelmChartPath struct {
	Path                 string   `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HelmChartPath) Reset()         { *m = HelmChartPath{} }
func (m *HelmChartPath) String() string { return proto.CompactTextString(m) }
func (*HelmChartPath) ProtoMessage()    {}
func (*HelmChartPath) Descriptor() ([]byte, []int) {
	return fileDescriptor_install_da8b87227e98e8d3, []int{2}
}
func (m *HelmChartPath) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HelmChartPath.Unmarshal(m, b)
}
func (m *HelmChartPath) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HelmChartPath.Marshal(b, m, deterministic)
}
func (dst *HelmChartPath) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelmChartPath.Merge(dst, src)
}
func (m *HelmChartPath) XXX_Size() int {
	return xxx_messageInfo_HelmChartPath.Size(m)
}
func (m *HelmChartPath) XXX_DiscardUnknown() {
	xxx_messageInfo_HelmChartPath.DiscardUnknown(m)
}

var xxx_messageInfo_HelmChartPath proto.InternalMessageInfo

func (m *HelmChartPath) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func init() {
	proto.RegisterType((*Install)(nil), "supergloo.solo.io.Install")
	proto.RegisterType((*HelmChartLocator)(nil), "supergloo.solo.io.HelmChartLocator")
	proto.RegisterType((*HelmChartPath)(nil), "supergloo.solo.io.HelmChartPath")
}
func (this *Install) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Install)
	if !ok {
		that2, ok := that.(Install)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.Status.Equal(&that1.Status) {
		return false
	}
	if !this.Metadata.Equal(&that1.Metadata) {
		return false
	}
	if that1.MeshType == nil {
		if this.MeshType != nil {
			return false
		}
	} else if this.MeshType == nil {
		return false
	} else if !this.MeshType.Equal(that1.MeshType) {
		return false
	}
	if !this.ChartLocator.Equal(that1.ChartLocator) {
		return false
	}
	if !this.Encryption.Equal(that1.Encryption) {
		return false
	}
	if this.Enabled != that1.Enabled {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *Install_Istio) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Install_Istio)
	if !ok {
		that2, ok := that.(Install_Istio)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.Istio.Equal(that1.Istio) {
		return false
	}
	return true
}
func (this *Install_Linkerd2) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Install_Linkerd2)
	if !ok {
		that2, ok := that.(Install_Linkerd2)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.Linkerd2.Equal(that1.Linkerd2) {
		return false
	}
	return true
}
func (this *Install_Consul) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Install_Consul)
	if !ok {
		that2, ok := that.(Install_Consul)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.Consul.Equal(that1.Consul) {
		return false
	}
	return true
}
func (this *HelmChartLocator) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HelmChartLocator)
	if !ok {
		that2, ok := that.(HelmChartLocator)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if that1.Kind == nil {
		if this.Kind != nil {
			return false
		}
	} else if this.Kind == nil {
		return false
	} else if !this.Kind.Equal(that1.Kind) {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *HelmChartLocator_ChartPath) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HelmChartLocator_ChartPath)
	if !ok {
		that2, ok := that.(HelmChartLocator_ChartPath)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.ChartPath.Equal(that1.ChartPath) {
		return false
	}
	return true
}
func (this *HelmChartPath) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HelmChartPath)
	if !ok {
		that2, ok := that.(HelmChartPath)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Path != that1.Path {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}

func init() { proto.RegisterFile("install.proto", fileDescriptor_install_da8b87227e98e8d3) }

var fileDescriptor_install_da8b87227e98e8d3 = []byte{
	// 448 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xd1, 0x6a, 0xd4, 0x40,
	0x14, 0x86, 0xb3, 0xba, 0x66, 0x77, 0x4f, 0x5b, 0x68, 0x87, 0x45, 0xa6, 0x15, 0xdb, 0x25, 0xbd,
	0xb0, 0x37, 0x4d, 0x6c, 0x7b, 0xa3, 0x82, 0x20, 0x29, 0xd2, 0x14, 0x2a, 0xc8, 0x78, 0x27, 0x82,
	0xcc, 0x26, 0xd3, 0xec, 0xb0, 0xd9, 0x9c, 0x30, 0x33, 0x2b, 0xf4, 0x8d, 0x7c, 0x14, 0x5f, 0xc2,
	0x5e, 0xf8, 0x08, 0x3e, 0x81, 0x64, 0x32, 0xd9, 0x76, 0x35, 0x48, 0xaf, 0x32, 0xc9, 0xf9, 0xbe,
	0x3f, 0x27, 0x39, 0x07, 0xb6, 0x64, 0xa9, 0x0d, 0x2f, 0x8a, 0xb0, 0x52, 0x68, 0x90, 0xec, 0xe8,
	0x65, 0x25, 0x54, 0x5e, 0x20, 0x86, 0x1a, 0x0b, 0x0c, 0x25, 0xee, 0x8d, 0x73, 0xcc, 0xd1, 0x56,
	0xa3, 0xfa, 0xd4, 0x80, 0x7b, 0x27, 0xb9, 0x34, 0xb3, 0xe5, 0x34, 0x4c, 0x71, 0x11, 0xd5, 0xe4,
	0xb1, 0xc4, 0xe6, 0x3a, 0x97, 0x26, 0xe2, 0x95, 0x8c, 0xbe, 0x9d, 0x44, 0x0b, 0x61, 0x78, 0xc6,
	0x0d, 0x77, 0x4a, 0xf4, 0x00, 0x45, 0x1b, 0x6e, 0x96, 0xda, 0x09, 0xdb, 0xa2, 0x4c, 0xd5, 0x4d,
	0x65, 0x24, 0x96, 0xee, 0x09, 0x2c, 0x84, 0x9e, 0x35, 0xe7, 0xe0, 0xe7, 0x63, 0x18, 0x5c, 0x36,
	0xcd, 0x93, 0x0b, 0xf0, 0x1b, 0x93, 0xf6, 0x26, 0xbd, 0xa3, 0x8d, 0xd3, 0x71, 0x98, 0xa2, 0x12,
	0xed, 0x27, 0x84, 0x9f, 0x6c, 0x2d, 0xde, 0xfd, 0x71, 0x7b, 0xe0, 0xfd, 0xbe, 0x3d, 0xd8, 0x31,
	0x42, 0x9b, 0x4c, 0x5e, 0x5f, 0xbf, 0x09, 0x64, 0x5e, 0xa2, 0x12, 0x01, 0x73, 0x3a, 0x79, 0x05,
	0xc3, 0xb6, 0x6b, 0xfa, 0xc8, 0x46, 0x3d, 0x5d, 0x8f, 0xfa, 0xe0, 0xaa, 0x71, 0xbf, 0x0e, 0x63,
	0x2b, 0x9a, 0xbc, 0x84, 0x27, 0x52, 0x1b, 0x89, 0x14, 0xac, 0x46, 0xc3, 0x7f, 0xfe, 0x64, 0x78,
	0x59, 0xd7, 0x13, 0x8f, 0x35, 0x20, 0x79, 0x0d, 0xc3, 0x42, 0x96, 0x73, 0xa1, 0xb2, 0x53, 0x3a,
	0xb6, 0xd2, 0xb3, 0x0e, 0xe9, 0xca, 0x21, 0x89, 0xc7, 0x56, 0x38, 0x39, 0x03, 0x3f, 0xc5, 0x52,
	0x2f, 0x0b, 0xba, 0x6f, 0xc5, 0xdd, 0x0e, 0xf1, 0xdc, 0x02, 0x89, 0xc7, 0x1c, 0x4a, 0x2e, 0x60,
	0x33, 0x9d, 0x71, 0x65, 0xae, 0x30, 0xe5, 0x06, 0x15, 0xf5, 0xad, 0x7a, 0xd8, 0xa1, 0x26, 0xa2,
	0x58, 0x9c, 0xdf, 0x43, 0xd9, 0x9a, 0x48, 0xde, 0x02, 0xdc, 0x4d, 0x86, 0x0e, 0x6c, 0xcc, 0xf3,
	0x8e, 0x98, 0xf7, 0x2b, 0x88, 0xdd, 0x13, 0x08, 0x85, 0x81, 0x28, 0xf9, 0xb4, 0x10, 0x19, 0xdd,
	0x9c, 0xf4, 0x8e, 0x86, 0xac, 0xbd, 0x8d, 0x37, 0x60, 0x54, 0x0f, 0xf8, 0xab, 0xb9, 0xa9, 0x44,
	0xf0, 0x05, 0xb6, 0xff, 0xee, 0x83, 0xbc, 0x83, 0x91, 0xed, 0xe4, 0x23, 0x37, 0x33, 0x37, 0xea,
	0xc9, 0xff, 0xfa, 0xaf, 0xb9, 0xc4, 0x63, 0x77, 0x52, 0xec, 0x43, 0x7f, 0x2e, 0xcb, 0x2c, 0x38,
	0x84, 0xad, 0x35, 0x8a, 0x10, 0xe8, 0x57, 0x6d, 0xea, 0x88, 0xd9, 0x73, 0x7c, 0xfc, 0xfd, 0xd7,
	0x7e, 0xef, 0xf3, 0x8b, 0xae, 0xbd, 0x6d, 0xdf, 0x19, 0x55, 0xf3, 0xdc, 0x2d, 0xef, 0xd4, 0xb7,
	0x8b, 0x79, 0xf6, 0x27, 0x00, 0x00, 0xff, 0xff, 0xd2, 0x30, 0x5e, 0x0e, 0x54, 0x03, 0x00, 0x00,
}
