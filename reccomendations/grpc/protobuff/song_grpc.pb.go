// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: song.proto

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	SongService_GetSong_FullMethodName = "/songRecommendation.SongService/GetSong"
)

// SongServiceClient is the client API for SongService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SongServiceClient interface {
	// takes in userID and returns an array of posssible songs of the day, specific to that user based on prefrences
	GetSong(ctx context.Context, in *SongRequest, opts ...grpc.CallOption) (*SongResponse, error)
}

type songServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSongServiceClient(cc grpc.ClientConnInterface) SongServiceClient {
	return &songServiceClient{cc}
}

func (c *songServiceClient) GetSong(ctx context.Context, in *SongRequest, opts ...grpc.CallOption) (*SongResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SongResponse)
	err := c.cc.Invoke(ctx, SongService_GetSong_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SongServiceServer is the server API for SongService service.
// All implementations must embed UnimplementedSongServiceServer
// for forward compatibility.
type SongServiceServer interface {
	// takes in userID and returns an array of posssible songs of the day, specific to that user based on prefrences
	GetSong(context.Context, *SongRequest) (*SongResponse, error)
	mustEmbedUnimplementedSongServiceServer()
}

// UnimplementedSongServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSongServiceServer struct{}

func (UnimplementedSongServiceServer) GetSong(context.Context, *SongRequest) (*SongResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSong not implemented")
}
func (UnimplementedSongServiceServer) mustEmbedUnimplementedSongServiceServer() {}
func (UnimplementedSongServiceServer) testEmbeddedByValue()                     {}

// UnsafeSongServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SongServiceServer will
// result in compilation errors.
type UnsafeSongServiceServer interface {
	mustEmbedUnimplementedSongServiceServer()
}

func RegisterSongServiceServer(s grpc.ServiceRegistrar, srv SongServiceServer) {
	// If the following call pancis, it indicates UnimplementedSongServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SongService_ServiceDesc, srv)
}

func _SongService_GetSong_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SongRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SongServiceServer).GetSong(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SongService_GetSong_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SongServiceServer).GetSong(ctx, req.(*SongRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SongService_ServiceDesc is the grpc.ServiceDesc for SongService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SongService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "songRecommendation.SongService",
	HandlerType: (*SongServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSong",
			Handler:    _SongService_GetSong_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "song.proto",
}
