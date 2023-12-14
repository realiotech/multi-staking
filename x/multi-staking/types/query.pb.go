// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: multistaking/v1/query.proto

package types

import (
	context "context"
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	_ "github.com/cosmos/cosmos-sdk/x/staking/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type QueryBondTokenWeightRequest struct {
	TokenDenom string `protobuf:"bytes,1,opt,name=token_denom,json=tokenDenom,proto3" json:"token_denom,omitempty"`
}

func (m *QueryBondTokenWeightRequest) Reset()         { *m = QueryBondTokenWeightRequest{} }
func (m *QueryBondTokenWeightRequest) String() string { return proto.CompactTextString(m) }
func (*QueryBondTokenWeightRequest) ProtoMessage()    {}
func (*QueryBondTokenWeightRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_82d174b604da394d, []int{0}
}
func (m *QueryBondTokenWeightRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryBondTokenWeightRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryBondTokenWeightRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryBondTokenWeightRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryBondTokenWeightRequest.Merge(m, src)
}
func (m *QueryBondTokenWeightRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryBondTokenWeightRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryBondTokenWeightRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryBondTokenWeightRequest proto.InternalMessageInfo

func (m *QueryBondTokenWeightRequest) GetTokenDenom() string {
	if m != nil {
		return m.TokenDenom
	}
	return ""
}

type QueryBondTokenWeightResponse struct {
	Weight cosmossdk_io_math.LegacyDec `protobuf:"bytes,1,opt,name=weight,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"weight"`
}

func (m *QueryBondTokenWeightResponse) Reset()         { *m = QueryBondTokenWeightResponse{} }
func (m *QueryBondTokenWeightResponse) String() string { return proto.CompactTextString(m) }
func (*QueryBondTokenWeightResponse) ProtoMessage()    {}
func (*QueryBondTokenWeightResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_82d174b604da394d, []int{1}
}
func (m *QueryBondTokenWeightResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryBondTokenWeightResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryBondTokenWeightResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryBondTokenWeightResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryBondTokenWeightResponse.Merge(m, src)
}
func (m *QueryBondTokenWeightResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryBondTokenWeightResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryBondTokenWeightResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryBondTokenWeightResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*QueryBondTokenWeightRequest)(nil), "multistaking.v1.QueryBondTokenWeightRequest")
	proto.RegisterType((*QueryBondTokenWeightResponse)(nil), "multistaking.v1.QueryBondTokenWeightResponse")
}

func init() { proto.RegisterFile("multistaking/v1/query.proto", fileDescriptor_82d174b604da394d) }

var fileDescriptor_82d174b604da394d = []byte{
	// 395 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0x3f, 0x6b, 0xdb, 0x40,
	0x14, 0x97, 0x0a, 0x35, 0xf4, 0x3a, 0x18, 0x44, 0xa1, 0xad, 0x6d, 0xe4, 0xe2, 0x76, 0xe8, 0x50,
	0xeb, 0x50, 0x3b, 0x14, 0x53, 0x68, 0xc0, 0x78, 0xcc, 0x12, 0x27, 0x10, 0x48, 0x06, 0x73, 0x92,
	0x8e, 0xd3, 0x61, 0xeb, 0x9e, 0xec, 0x3b, 0x39, 0x31, 0x21, 0x4b, 0x3e, 0x41, 0x20, 0x1f, 0x24,
	0x73, 0xbe, 0x81, 0x47, 0x43, 0x96, 0x90, 0xc1, 0x04, 0x3b, 0x1f, 0x24, 0x48, 0x77, 0x26, 0x46,
	0x84, 0x90, 0x4d, 0xef, 0xf7, 0xe7, 0xbd, 0xdf, 0x7b, 0x3a, 0x54, 0x4f, 0xb2, 0x91, 0xe2, 0x52,
	0x91, 0x21, 0x17, 0x0c, 0x4f, 0x7d, 0x3c, 0xce, 0xe8, 0x64, 0xe6, 0xa5, 0x13, 0x50, 0xe0, 0x54,
	0xb7, 0x49, 0x6f, 0xea, 0xd7, 0x3e, 0x31, 0x60, 0x50, 0x70, 0x38, 0xff, 0xd2, 0xb2, 0x5a, 0x83,
	0x01, 0xb0, 0x11, 0xc5, 0x24, 0xe5, 0x98, 0x08, 0x01, 0x8a, 0x28, 0x0e, 0x42, 0x1a, 0xf6, 0x6b,
	0x08, 0x32, 0x01, 0x39, 0xd0, 0x36, 0x5d, 0x18, 0xca, 0xd5, 0x15, 0x0e, 0x88, 0xa4, 0x78, 0xea,
	0x07, 0x54, 0x11, 0x1f, 0x87, 0xc0, 0x85, 0xe1, 0x7f, 0x18, 0xfe, 0x39, 0x9e, 0x96, 0x6c, 0x12,
	0x69, 0xd5, 0x67, 0xa3, 0x4a, 0x64, 0xb1, 0x40, 0x22, 0x0d, 0xd1, 0xfa, 0x8f, 0xea, 0x7b, 0xf9,
	0x36, 0x5d, 0x10, 0xd1, 0x01, 0x0c, 0xa9, 0x38, 0xa4, 0x9c, 0xc5, 0xaa, 0x4f, 0xc7, 0x19, 0x95,
	0xca, 0x69, 0xa2, 0x8f, 0x2a, 0x47, 0x07, 0x11, 0x15, 0x90, 0x7c, 0xb1, 0xbf, 0xd9, 0x3f, 0x3f,
	0xf4, 0x51, 0x01, 0xf5, 0x72, 0xa4, 0x75, 0x8c, 0x1a, 0x2f, 0xfb, 0x65, 0x0a, 0x42, 0x52, 0xe7,
	0x1f, 0xaa, 0x9c, 0x14, 0x88, 0xf6, 0x76, 0xbf, 0xcf, 0x97, 0x4d, 0xeb, 0x7e, 0xd9, 0xac, 0xeb,
	0x40, 0x32, 0x1a, 0x7a, 0x1c, 0x70, 0x42, 0x54, 0xec, 0xed, 0x52, 0x46, 0xc2, 0x59, 0x8f, 0x86,
	0x7d, 0x63, 0xf9, 0x7d, 0x63, 0xa3, 0xf7, 0x45, 0x77, 0xe7, 0xda, 0x46, 0xd5, 0xd2, 0x08, 0xe7,
	0x97, 0x57, 0x3a, 0xbd, 0xf7, 0xca, 0x26, 0xb5, 0xf6, 0x1b, 0xd5, 0x3a, 0x77, 0x6b, 0xe7, 0xe2,
	0xf6, 0xf1, 0xea, 0x5d, 0xc7, 0xf9, 0x8b, 0x27, 0x94, 0x8c, 0xf2, 0x88, 0xa5, 0x37, 0x10, 0x80,
	0x88, 0x06, 0xfa, 0x36, 0x3a, 0x2d, 0x3e, 0xdb, 0xba, 0xd4, 0x79, 0x77, 0x7f, 0xbe, 0x72, 0xed,
	0xc5, 0xca, 0xb5, 0x1f, 0x56, 0xae, 0x7d, 0xb9, 0x76, 0xad, 0xc5, 0xda, 0xb5, 0xee, 0xd6, 0xae,
	0x75, 0xd4, 0x61, 0x5c, 0xc5, 0x59, 0xe0, 0x85, 0x90, 0x98, 0xe6, 0x8a, 0x86, 0xb1, 0x1e, 0xd0,
	0xde, 0x4c, 0x38, 0x2d, 0xd5, 0x6a, 0x96, 0x52, 0x19, 0x54, 0x8a, 0x9f, 0xf6, 0xe7, 0x29, 0x00,
	0x00, 0xff, 0xff, 0x56, 0x32, 0x73, 0x0c, 0x92, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	BondTokenWeight(ctx context.Context, in *QueryBondTokenWeightRequest, opts ...grpc.CallOption) (*QueryBondTokenWeightResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) BondTokenWeight(ctx context.Context, in *QueryBondTokenWeightRequest, opts ...grpc.CallOption) (*QueryBondTokenWeightResponse, error) {
	out := new(QueryBondTokenWeightResponse)
	err := c.cc.Invoke(ctx, "/multistaking.v1.Query/BondTokenWeight", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	BondTokenWeight(context.Context, *QueryBondTokenWeightRequest) (*QueryBondTokenWeightResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) BondTokenWeight(ctx context.Context, req *QueryBondTokenWeightRequest) (*QueryBondTokenWeightResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BondTokenWeight not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_BondTokenWeight_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryBondTokenWeightRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).BondTokenWeight(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/multistaking.v1.Query/BondTokenWeight",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).BondTokenWeight(ctx, req.(*QueryBondTokenWeightRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "multistaking.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "BondTokenWeight",
			Handler:    _Query_BondTokenWeight_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "multistaking/v1/query.proto",
}

func (m *QueryBondTokenWeightRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryBondTokenWeightRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryBondTokenWeightRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TokenDenom) > 0 {
		i -= len(m.TokenDenom)
		copy(dAtA[i:], m.TokenDenom)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.TokenDenom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryBondTokenWeightResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryBondTokenWeightResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryBondTokenWeightResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Weight.Size()
		i -= size
		if _, err := m.Weight.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryBondTokenWeightRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.TokenDenom)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryBondTokenWeightResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Weight.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryBondTokenWeightRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryBondTokenWeightRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryBondTokenWeightRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TokenDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryBondTokenWeightResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryBondTokenWeightResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryBondTokenWeightResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Weight", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Weight.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
