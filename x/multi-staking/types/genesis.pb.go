// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: multistaking/v1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/codec/types"
	_ "github.com/cosmos/cosmos-sdk/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	types "github.com/cosmos/cosmos-sdk/x/staking/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/protobuf/types/known/timestamppb"
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

type GenesisState struct {
	MultiStakingLocks          []MultiStakingLock          `protobuf:"bytes,1,rep,name=multi_staking_locks,json=multiStakingLocks,proto3" json:"multi_staking_locks"`
	MultiStakingUnlocks        []MultiStakingUnlock        `protobuf:"bytes,2,rep,name=multi_staking_unlocks,json=multiStakingUnlocks,proto3" json:"multi_staking_unlocks"`
	MultiStakingCoinInfo       []MultiStakingCoinInfo      `protobuf:"bytes,3,rep,name=multi_staking_coin_info,json=multiStakingCoinInfo,proto3" json:"multi_staking_coin_info"`
	ValidatorMultiStakingCoins []ValidatorMultiStakingCoin `protobuf:"bytes,4,rep,name=validator_multi_staking_coins,json=validatorMultiStakingCoins,proto3" json:"validator_multi_staking_coins"`
	IntermediaryDelegators     []string                    `protobuf:"bytes,5,rep,name=IntermediaryDelegators,proto3" json:"IntermediaryDelegators,omitempty"`
	StakingGenesisState        types.GenesisState          `protobuf:"bytes,6,opt,name=staking_genesis_state,json=stakingGenesisState,proto3" json:"staking_genesis_state"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f95a201ebed173c, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetMultiStakingLocks() []MultiStakingLock {
	if m != nil {
		return m.MultiStakingLocks
	}
	return nil
}

func (m *GenesisState) GetMultiStakingUnlocks() []MultiStakingUnlock {
	if m != nil {
		return m.MultiStakingUnlocks
	}
	return nil
}

func (m *GenesisState) GetMultiStakingCoinInfo() []MultiStakingCoinInfo {
	if m != nil {
		return m.MultiStakingCoinInfo
	}
	return nil
}

func (m *GenesisState) GetValidatorMultiStakingCoins() []ValidatorMultiStakingCoin {
	if m != nil {
		return m.ValidatorMultiStakingCoins
	}
	return nil
}

func (m *GenesisState) GetIntermediaryDelegators() []string {
	if m != nil {
		return m.IntermediaryDelegators
	}
	return nil
}

func (m *GenesisState) GetStakingGenesisState() types.GenesisState {
	if m != nil {
		return m.StakingGenesisState
	}
	return types.GenesisState{}
}

type MultiStakingCoinInfo struct {
	Denom      string                                 `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	BondWeight github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=bond_weight,json=bondWeight,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"bond_weight"`
}

func (m *MultiStakingCoinInfo) Reset()         { *m = MultiStakingCoinInfo{} }
func (m *MultiStakingCoinInfo) String() string { return proto.CompactTextString(m) }
func (*MultiStakingCoinInfo) ProtoMessage()    {}
func (*MultiStakingCoinInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f95a201ebed173c, []int{1}
}
func (m *MultiStakingCoinInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MultiStakingCoinInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MultiStakingCoinInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MultiStakingCoinInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MultiStakingCoinInfo.Merge(m, src)
}
func (m *MultiStakingCoinInfo) XXX_Size() int {
	return m.Size()
}
func (m *MultiStakingCoinInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_MultiStakingCoinInfo.DiscardUnknown(m)
}

var xxx_messageInfo_MultiStakingCoinInfo proto.InternalMessageInfo

func (m *MultiStakingCoinInfo) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "multistaking.v1.GenesisState")
	proto.RegisterType((*MultiStakingCoinInfo)(nil), "multistaking.v1.MultiStakingCoinInfo")
}

func init() { proto.RegisterFile("multistaking/v1/genesis.proto", fileDescriptor_8f95a201ebed173c) }

var fileDescriptor_8f95a201ebed173c = []byte{
	// 534 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x53, 0xd1, 0x6a, 0xdb, 0x30,
	0x14, 0x8d, 0x9b, 0xb4, 0x50, 0x65, 0x30, 0xe6, 0xa6, 0xab, 0x1b, 0xa8, 0x93, 0xb5, 0xdd, 0x08,
	0x83, 0xd8, 0xa4, 0x83, 0xc1, 0x60, 0x4f, 0x59, 0x60, 0x14, 0xb6, 0x97, 0x84, 0xad, 0x30, 0xe8,
	0x8c, 0x6c, 0x2b, 0x8e, 0x88, 0x25, 0x05, 0x4b, 0xf1, 0x96, 0x6f, 0xd8, 0xcb, 0x3e, 0x66, 0xb0,
	0x5f, 0xe8, 0x63, 0xd9, 0xd3, 0xd8, 0x43, 0x19, 0xc9, 0x8f, 0x14, 0x4b, 0x32, 0x4d, 0x9c, 0x36,
	0x4f, 0xd6, 0xd1, 0x39, 0x3e, 0xe7, 0xea, 0x72, 0x2f, 0x38, 0x22, 0xd3, 0x58, 0x60, 0x2e, 0xe0,
	0x18, 0xd3, 0xc8, 0x4d, 0x3b, 0x6e, 0x84, 0x28, 0xe2, 0x98, 0x3b, 0x93, 0x84, 0x09, 0x66, 0x3e,
	0x5e, 0xa6, 0x9d, 0xb4, 0x53, 0x3f, 0x8c, 0x18, 0x8b, 0x62, 0xe4, 0x4a, 0xda, 0x9f, 0x0e, 0x5d,
	0x48, 0x67, 0x4a, 0x5b, 0x6f, 0x14, 0x29, 0x81, 0x09, 0xe2, 0x02, 0x92, 0x89, 0x16, 0xd4, 0x22,
	0x16, 0x31, 0x79, 0x74, 0xb3, 0x93, 0xbe, 0x3d, 0x0c, 0x18, 0x27, 0x8c, 0x7b, 0x8a, 0x50, 0x40,
	0x53, 0xb6, 0x42, 0xae, 0x0f, 0x39, 0x72, 0xd3, 0x8e, 0x8f, 0x04, 0xec, 0xb8, 0x01, 0xc3, 0x54,
	0xf3, 0xa7, 0x9a, 0xbf, 0x2b, 0x5f, 0x49, 0x56, 0xde, 0x50, 0x3f, 0x29, 0x3e, 0x51, 0x62, 0x2f,
	0x7f, 0x94, 0x12, 0x1d, 0x68, 0x2b, 0xc2, 0x95, 0x84, 0x6b, 0xe2, 0xf8, 0x77, 0x05, 0x3c, 0x7a,
	0xaf, 0xfc, 0x06, 0x02, 0x0a, 0x64, 0x5e, 0x80, 0xbd, 0x15, 0x03, 0x2f, 0x66, 0xc1, 0x98, 0x5b,
	0x46, 0xb3, 0xdc, 0xaa, 0x9e, 0x3d, 0x73, 0x0a, 0x0d, 0x73, 0x3e, 0x66, 0x78, 0xa0, 0xf0, 0x07,
	0x16, 0x8c, 0xbb, 0x95, 0xab, 0x9b, 0x46, 0xa9, 0xff, 0x84, 0x14, 0xee, 0xb9, 0x79, 0x09, 0xf6,
	0x57, 0x8d, 0xa7, 0x54, 0x59, 0x6f, 0x49, 0xeb, 0x93, 0x8d, 0xd6, 0x9f, 0xa4, 0x56, 0x9b, 0xef,
	0x91, 0x35, 0x86, 0x9b, 0x3e, 0x38, 0x58, 0xb5, 0xcf, 0x1a, 0xe9, 0x61, 0x3a, 0x64, 0x56, 0x59,
	0x06, 0x3c, 0xdf, 0x18, 0xf0, 0x8e, 0x61, 0x7a, 0x4e, 0x87, 0x4c, 0x47, 0xd4, 0xc8, 0x3d, 0x9c,
	0xc9, 0xc1, 0x51, 0x0a, 0x63, 0x1c, 0x42, 0xc1, 0x12, 0x6f, 0x3d, 0x8d, 0x5b, 0x15, 0x99, 0xf4,
	0x72, 0x2d, 0xe9, 0x73, 0xfe, 0x57, 0x31, 0x52, 0xc7, 0xd5, 0xd3, 0x87, 0x04, 0xdc, 0x7c, 0x0d,
	0x9e, 0x9e, 0x53, 0x81, 0x12, 0x82, 0x42, 0x0c, 0x93, 0x59, 0x0f, 0xc5, 0x28, 0xca, 0x94, 0xdc,
	0xda, 0x6e, 0x96, 0x5b, 0xbb, 0xfd, 0x07, 0x58, 0xf3, 0x2b, 0xd8, 0xcf, 0x8b, 0xd3, 0x03, 0x93,
	0x15, 0x2b, 0x90, 0xb5, 0xd3, 0x34, 0x5a, 0xd5, 0xb3, 0x53, 0x47, 0xcf, 0xe2, 0x5d, 0x99, 0x72,
	0xba, 0x9c, 0xe5, 0x69, 0xc8, 0x1b, 0xae, 0x35, 0xcb, 0xd4, 0xf1, 0x0f, 0x03, 0xd4, 0xee, 0xeb,
	0xa0, 0x59, 0x03, 0xdb, 0x21, 0xa2, 0x8c, 0x58, 0x46, 0xd3, 0x68, 0xed, 0xf6, 0x15, 0x30, 0x2f,
	0x41, 0xd5, 0x67, 0x34, 0xf4, 0xbe, 0x21, 0x1c, 0x8d, 0x84, 0xb5, 0x95, 0x71, 0xdd, 0xb7, 0x99,
	0xfd, 0xbf, 0x9b, 0xc6, 0x8b, 0x08, 0x8b, 0xd1, 0xd4, 0x77, 0x02, 0x46, 0xf4, 0x8a, 0xe8, 0x4f,
	0x9b, 0x87, 0x63, 0x57, 0xcc, 0x26, 0x88, 0x3b, 0x3d, 0x14, 0xfc, 0xf9, 0xd5, 0x06, 0xba, 0xea,
	0x1e, 0x0a, 0xfa, 0x20, 0x33, 0xbc, 0x90, 0x7e, 0xdd, 0xc1, 0xd5, 0xdc, 0x36, 0xae, 0xe7, 0xb6,
	0xf1, 0x7f, 0x6e, 0x1b, 0x3f, 0x17, 0x76, 0xe9, 0x7a, 0x61, 0x97, 0xfe, 0x2e, 0xec, 0xd2, 0x97,
	0x37, 0x4b, 0xde, 0x09, 0x82, 0x31, 0x66, 0x02, 0x05, 0x23, 0xb5, 0x25, 0xed, 0x7c, 0x6d, 0xbe,
	0x17, 0xb0, 0x8c, 0xf4, 0x77, 0xe4, 0x8e, 0xbc, 0xba, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x47, 0x41,
	0x5a, 0xce, 0x46, 0x04, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.StakingGenesisState.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	if len(m.IntermediaryDelegators) > 0 {
		for iNdEx := len(m.IntermediaryDelegators) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.IntermediaryDelegators[iNdEx])
			copy(dAtA[i:], m.IntermediaryDelegators[iNdEx])
			i = encodeVarintGenesis(dAtA, i, uint64(len(m.IntermediaryDelegators[iNdEx])))
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.ValidatorMultiStakingCoins) > 0 {
		for iNdEx := len(m.ValidatorMultiStakingCoins) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ValidatorMultiStakingCoins[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.MultiStakingCoinInfo) > 0 {
		for iNdEx := len(m.MultiStakingCoinInfo) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.MultiStakingCoinInfo[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.MultiStakingUnlocks) > 0 {
		for iNdEx := len(m.MultiStakingUnlocks) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.MultiStakingUnlocks[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.MultiStakingLocks) > 0 {
		for iNdEx := len(m.MultiStakingLocks) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.MultiStakingLocks[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *MultiStakingCoinInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MultiStakingCoinInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MultiStakingCoinInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.BondWeight.Size()
		i -= size
		if _, err := m.BondWeight.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.MultiStakingLocks) > 0 {
		for _, e := range m.MultiStakingLocks {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.MultiStakingUnlocks) > 0 {
		for _, e := range m.MultiStakingUnlocks {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.MultiStakingCoinInfo) > 0 {
		for _, e := range m.MultiStakingCoinInfo {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.ValidatorMultiStakingCoins) > 0 {
		for _, e := range m.ValidatorMultiStakingCoins {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.IntermediaryDelegators) > 0 {
		for _, s := range m.IntermediaryDelegators {
			l = len(s)
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	l = m.StakingGenesisState.Size()
	n += 1 + l + sovGenesis(uint64(l))
	return n
}

func (m *MultiStakingCoinInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = m.BondWeight.Size()
	n += 1 + l + sovGenesis(uint64(l))
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MultiStakingLocks", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MultiStakingLocks = append(m.MultiStakingLocks, MultiStakingLock{})
			if err := m.MultiStakingLocks[len(m.MultiStakingLocks)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MultiStakingUnlocks", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MultiStakingUnlocks = append(m.MultiStakingUnlocks, MultiStakingUnlock{})
			if err := m.MultiStakingUnlocks[len(m.MultiStakingUnlocks)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MultiStakingCoinInfo", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MultiStakingCoinInfo = append(m.MultiStakingCoinInfo, MultiStakingCoinInfo{})
			if err := m.MultiStakingCoinInfo[len(m.MultiStakingCoinInfo)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorMultiStakingCoins", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidatorMultiStakingCoins = append(m.ValidatorMultiStakingCoins, ValidatorMultiStakingCoin{})
			if err := m.ValidatorMultiStakingCoins[len(m.ValidatorMultiStakingCoins)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IntermediaryDelegators", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IntermediaryDelegators = append(m.IntermediaryDelegators, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StakingGenesisState", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.StakingGenesisState.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *MultiStakingCoinInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: MultiStakingCoinInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MultiStakingCoinInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BondWeight", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.BondWeight.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)