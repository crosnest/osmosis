// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: osmosis/concentrated-liquidity/tickInfo.proto

// this is a legacy package that requires additional migration logic
// in order to use the correct packge. Decision made to use legacy package path
// until clear steps for migration logic and the unknowns for state breaking are
// investigated for changing proto package.

package model

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type TickInfo struct {
	LiquidityGross   github_com_cosmos_cosmos_sdk_types.Dec      `protobuf:"bytes,1,opt,name=liquidity_gross,json=liquidityGross,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"liquidity_gross" yaml:"liquidity_gross"`
	LiquidityNet     github_com_cosmos_cosmos_sdk_types.Dec      `protobuf:"bytes,2,opt,name=liquidity_net,json=liquidityNet,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"liquidity_net" yaml:"liquidity_net"`
	FeeGrowthOutside github_com_cosmos_cosmos_sdk_types.DecCoins `protobuf:"bytes,3,rep,name=fee_growth_outside,json=feeGrowthOutside,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.DecCoins" json:"fee_growth_outside"`
}

func (m *TickInfo) Reset()         { *m = TickInfo{} }
func (m *TickInfo) String() string { return proto.CompactTextString(m) }
func (*TickInfo) ProtoMessage()    {}
func (*TickInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_1ccb7e45032b943a, []int{0}
}
func (m *TickInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TickInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TickInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TickInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TickInfo.Merge(m, src)
}
func (m *TickInfo) XXX_Size() int {
	return m.Size()
}
func (m *TickInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_TickInfo.DiscardUnknown(m)
}

var xxx_messageInfo_TickInfo proto.InternalMessageInfo

func (m *TickInfo) GetFeeGrowthOutside() github_com_cosmos_cosmos_sdk_types.DecCoins {
	if m != nil {
		return m.FeeGrowthOutside
	}
	return nil
}

func init() {
	proto.RegisterType((*TickInfo)(nil), "osmosis.concentratedliquidity.v1beta1.TickInfo")
}

func init() {
	proto.RegisterFile("osmosis/concentrated-liquidity/tickInfo.proto", fileDescriptor_1ccb7e45032b943a)
}

var fileDescriptor_1ccb7e45032b943a = []byte{
	// 374 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0xcd, 0x2f, 0xce, 0xcd,
	0x2f, 0xce, 0x2c, 0xd6, 0x4f, 0xce, 0xcf, 0x4b, 0x4e, 0xcd, 0x2b, 0x29, 0x4a, 0x2c, 0x49, 0x4d,
	0xd1, 0xcd, 0xc9, 0x2c, 0x2c, 0xcd, 0x4c, 0xc9, 0x2c, 0xa9, 0xd4, 0x2f, 0xc9, 0x4c, 0xce, 0xf6,
	0xcc, 0x4b, 0xcb, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x52, 0x85, 0x2a, 0xd7, 0x43, 0x56,
	0x0e, 0x57, 0xad, 0x57, 0x66, 0x98, 0x94, 0x5a, 0x92, 0x68, 0x28, 0x25, 0x99, 0x0c, 0x56, 0x17,
	0x0f, 0xd6, 0xa4, 0x0f, 0xe1, 0x40, 0x4c, 0x90, 0x12, 0x49, 0xcf, 0x4f, 0xcf, 0x87, 0x88, 0x83,
	0x58, 0x50, 0x51, 0x39, 0x88, 0x1a, 0xfd, 0xa4, 0xc4, 0xe2, 0x54, 0x7d, 0xa8, 0x29, 0xfa, 0xc9,
	0xf9, 0x99, 0x79, 0x10, 0x79, 0xa5, 0xd7, 0x4c, 0x5c, 0x1c, 0x21, 0x50, 0xa7, 0x08, 0x15, 0x72,
	0xf1, 0xc3, 0xad, 0x8c, 0x4f, 0x2f, 0xca, 0x2f, 0x2e, 0x96, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x74,
	0xf2, 0x38, 0x71, 0x4f, 0x9e, 0xe1, 0xd6, 0x3d, 0x79, 0xb5, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24,
	0xbd, 0xe4, 0xfc, 0x5c, 0xa8, 0xe5, 0x50, 0x4a, 0xb7, 0x38, 0x25, 0x5b, 0xbf, 0xa4, 0xb2, 0x20,
	0xb5, 0x58, 0xcf, 0x25, 0x35, 0xf9, 0xd3, 0x3d, 0x79, 0xb1, 0xca, 0xc4, 0xdc, 0x1c, 0x2b, 0x25,
	0x34, 0xe3, 0x94, 0x82, 0xf8, 0xe0, 0x22, 0xee, 0x20, 0x01, 0xa1, 0x6c, 0x2e, 0x5e, 0x84, 0x9a,
	0xbc, 0xd4, 0x12, 0x09, 0x26, 0xb0, 0x85, 0x6e, 0x24, 0x5b, 0x28, 0x82, 0x6e, 0x61, 0x5e, 0x6a,
	0x89, 0x52, 0x10, 0x0f, 0x9c, 0xef, 0x97, 0x5a, 0x22, 0x54, 0xcf, 0x25, 0x94, 0x96, 0x9a, 0x0a,
	0x72, 0x4a, 0x79, 0x49, 0x46, 0x7c, 0x7e, 0x69, 0x49, 0x71, 0x66, 0x4a, 0xaa, 0x04, 0xb3, 0x02,
	0xb3, 0x06, 0xb7, 0x91, 0x8c, 0x1e, 0x34, 0x34, 0x41, 0x21, 0x05, 0x0b, 0x6f, 0x90, 0xd9, 0xce,
	0xf9, 0x99, 0x79, 0x4e, 0xc6, 0x20, 0xf7, 0xac, 0xba, 0x2f, 0xaf, 0x4d, 0x9c, 0x7b, 0x40, 0x7a,
	0x8a, 0x83, 0x04, 0xd2, 0x52, 0x53, 0xdd, 0xc1, 0x76, 0xf9, 0x43, 0xac, 0x72, 0x8a, 0x39, 0xf1,
	0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8,
	0xf0, 0x58, 0x8e, 0xe1, 0xc6, 0x63, 0x39, 0x86, 0x28, 0x27, 0x24, 0x83, 0xa1, 0x49, 0x41, 0x37,
	0x27, 0x31, 0xa9, 0x18, 0xc6, 0xd1, 0x2f, 0x33, 0x34, 0xd1, 0xaf, 0xc0, 0x95, 0x98, 0x72, 0xf3,
	0x53, 0x52, 0x73, 0x92, 0xd8, 0xc0, 0x51, 0x6a, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x93, 0x5e,
	0x70, 0x0c, 0x7b, 0x02, 0x00, 0x00,
}

func (m *TickInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TickInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TickInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.FeeGrowthOutside) > 0 {
		for iNdEx := len(m.FeeGrowthOutside) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.FeeGrowthOutside[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTickInfo(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	{
		size := m.LiquidityNet.Size()
		i -= size
		if _, err := m.LiquidityNet.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTickInfo(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.LiquidityGross.Size()
		i -= size
		if _, err := m.LiquidityGross.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTickInfo(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintTickInfo(dAtA []byte, offset int, v uint64) int {
	offset -= sovTickInfo(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TickInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.LiquidityGross.Size()
	n += 1 + l + sovTickInfo(uint64(l))
	l = m.LiquidityNet.Size()
	n += 1 + l + sovTickInfo(uint64(l))
	if len(m.FeeGrowthOutside) > 0 {
		for _, e := range m.FeeGrowthOutside {
			l = e.Size()
			n += 1 + l + sovTickInfo(uint64(l))
		}
	}
	return n
}

func sovTickInfo(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTickInfo(x uint64) (n int) {
	return sovTickInfo(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TickInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTickInfo
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TickInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TickInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LiquidityGross", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTickInfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTickInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTickInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LiquidityGross.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LiquidityNet", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTickInfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTickInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTickInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LiquidityNet.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeGrowthOutside", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTickInfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTickInfo
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTickInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FeeGrowthOutside = append(m.FeeGrowthOutside, types.DecCoin{})
			if err := m.FeeGrowthOutside[len(m.FeeGrowthOutside)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTickInfo(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTickInfo
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTickInfo(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTickInfo
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTickInfo
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTickInfo
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTickInfo
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTickInfo
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTickInfo
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTickInfo        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTickInfo          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTickInfo = fmt.Errorf("proto: unexpected end of group")
)
