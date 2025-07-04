// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: discord/discord.proto

package discord

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Discord_GetGuildChannels_FullMethodName = "/discord.Discord/GetGuildChannels"
	Discord_GetGuildInfo_FullMethodName     = "/discord.Discord/GetGuildInfo"
	Discord_LeaveGuild_FullMethodName       = "/discord.Discord/LeaveGuild"
	Discord_GetGuildRoles_FullMethodName    = "/discord.Discord/GetGuildRoles"
)

// DiscordClient is the client API for Discord service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DiscordClient interface {
	GetGuildChannels(ctx context.Context, in *GetGuildChannelsRequest, opts ...grpc.CallOption) (*GetGuildChannelsResponse, error)
	GetGuildInfo(ctx context.Context, in *GetGuildInfoRequest, opts ...grpc.CallOption) (*GetGuildInfoResponse, error)
	LeaveGuild(ctx context.Context, in *LeaveGuildRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetGuildRoles(ctx context.Context, in *GetGuildRolesRequest, opts ...grpc.CallOption) (*GetGuildRolesResponse, error)
}

type discordClient struct {
	cc grpc.ClientConnInterface
}

func NewDiscordClient(cc grpc.ClientConnInterface) DiscordClient {
	return &discordClient{cc}
}

func (c *discordClient) GetGuildChannels(ctx context.Context, in *GetGuildChannelsRequest, opts ...grpc.CallOption) (*GetGuildChannelsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetGuildChannelsResponse)
	err := c.cc.Invoke(ctx, Discord_GetGuildChannels_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discordClient) GetGuildInfo(ctx context.Context, in *GetGuildInfoRequest, opts ...grpc.CallOption) (*GetGuildInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetGuildInfoResponse)
	err := c.cc.Invoke(ctx, Discord_GetGuildInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discordClient) LeaveGuild(ctx context.Context, in *LeaveGuildRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Discord_LeaveGuild_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discordClient) GetGuildRoles(ctx context.Context, in *GetGuildRolesRequest, opts ...grpc.CallOption) (*GetGuildRolesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetGuildRolesResponse)
	err := c.cc.Invoke(ctx, Discord_GetGuildRoles_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DiscordServer is the server API for Discord service.
// All implementations must embed UnimplementedDiscordServer
// for forward compatibility.
type DiscordServer interface {
	GetGuildChannels(context.Context, *GetGuildChannelsRequest) (*GetGuildChannelsResponse, error)
	GetGuildInfo(context.Context, *GetGuildInfoRequest) (*GetGuildInfoResponse, error)
	LeaveGuild(context.Context, *LeaveGuildRequest) (*emptypb.Empty, error)
	GetGuildRoles(context.Context, *GetGuildRolesRequest) (*GetGuildRolesResponse, error)
	mustEmbedUnimplementedDiscordServer()
}

// UnimplementedDiscordServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDiscordServer struct{}

func (UnimplementedDiscordServer) GetGuildChannels(context.Context, *GetGuildChannelsRequest) (*GetGuildChannelsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGuildChannels not implemented")
}
func (UnimplementedDiscordServer) GetGuildInfo(context.Context, *GetGuildInfoRequest) (*GetGuildInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGuildInfo not implemented")
}
func (UnimplementedDiscordServer) LeaveGuild(context.Context, *LeaveGuildRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaveGuild not implemented")
}
func (UnimplementedDiscordServer) GetGuildRoles(context.Context, *GetGuildRolesRequest) (*GetGuildRolesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGuildRoles not implemented")
}
func (UnimplementedDiscordServer) mustEmbedUnimplementedDiscordServer() {}
func (UnimplementedDiscordServer) testEmbeddedByValue()                 {}

// UnsafeDiscordServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DiscordServer will
// result in compilation errors.
type UnsafeDiscordServer interface {
	mustEmbedUnimplementedDiscordServer()
}

func RegisterDiscordServer(s grpc.ServiceRegistrar, srv DiscordServer) {
	// If the following call pancis, it indicates UnimplementedDiscordServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Discord_ServiceDesc, srv)
}

func _Discord_GetGuildChannels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGuildChannelsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscordServer).GetGuildChannels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Discord_GetGuildChannels_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscordServer).GetGuildChannels(ctx, req.(*GetGuildChannelsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Discord_GetGuildInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGuildInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscordServer).GetGuildInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Discord_GetGuildInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscordServer).GetGuildInfo(ctx, req.(*GetGuildInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Discord_LeaveGuild_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveGuildRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscordServer).LeaveGuild(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Discord_LeaveGuild_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscordServer).LeaveGuild(ctx, req.(*LeaveGuildRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Discord_GetGuildRoles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGuildRolesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscordServer).GetGuildRoles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Discord_GetGuildRoles_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscordServer).GetGuildRoles(ctx, req.(*GetGuildRolesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Discord_ServiceDesc is the grpc.ServiceDesc for Discord service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Discord_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "discord.Discord",
	HandlerType: (*DiscordServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetGuildChannels",
			Handler:    _Discord_GetGuildChannels_Handler,
		},
		{
			MethodName: "GetGuildInfo",
			Handler:    _Discord_GetGuildInfo_Handler,
		},
		{
			MethodName: "LeaveGuild",
			Handler:    _Discord_LeaveGuild_Handler,
		},
		{
			MethodName: "GetGuildRoles",
			Handler:    _Discord_GetGuildRoles_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "discord/discord.proto",
}
