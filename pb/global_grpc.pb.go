// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.1
// source: global.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GlobalClient is the client API for Global service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GlobalClient interface {
	SysInfo(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SysInfoResponse, error)
	Summary(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SummaryResponse, error)
	Shutdown(ctx context.Context, in *RequestWithId, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Restart(ctx context.Context, in *RequestWithId, opts ...grpc.CallOption) (*emptypb.Empty, error)
	StreamList(ctx context.Context, in *StreamListRequest, opts ...grpc.CallOption) (*StreamListResponse, error)
	WaitList(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StreamWaitListResponse, error)
	StreamInfo(ctx context.Context, in *StreamSnapRequest, opts ...grpc.CallOption) (*StreamInfoResponse, error)
	GetSubscribers(ctx context.Context, in *SubscribersRequest, opts ...grpc.CallOption) (*SubscribersResponse, error)
	AudioTrackSnap(ctx context.Context, in *StreamSnapRequest, opts ...grpc.CallOption) (*TrackSnapShotResponse, error)
	VideoTrackSnap(ctx context.Context, in *StreamSnapRequest, opts ...grpc.CallOption) (*TrackSnapShotResponse, error)
	ChangeSubscribe(ctx context.Context, in *ChangeSubscribeRequest, opts ...grpc.CallOption) (*SuccessResponse, error)
	StopSubscribe(ctx context.Context, in *RequestWithId, opts ...grpc.CallOption) (*SuccessResponse, error)
	GetConfig(ctx context.Context, in *GetConfigRequest, opts ...grpc.CallOption) (*GetConfigResponse, error)
	GetFormily(ctx context.Context, in *GetConfigRequest, opts ...grpc.CallOption) (*GetConfigResponse, error)
	ModifyConfig(ctx context.Context, in *ModifyConfigRequest, opts ...grpc.CallOption) (*SuccessResponse, error)
}

type globalClient struct {
	cc grpc.ClientConnInterface
}

func NewGlobalClient(cc grpc.ClientConnInterface) GlobalClient {
	return &globalClient{cc}
}

func (c *globalClient) SysInfo(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SysInfoResponse, error) {
	out := new(SysInfoResponse)
	err := c.cc.Invoke(ctx, "/m7s.Global/SysInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *globalClient) Summary(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SummaryResponse, error) {
	out := new(SummaryResponse)
	err := c.cc.Invoke(ctx, "/m7s.Global/Summary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *globalClient) Shutdown(ctx context.Context, in *RequestWithId, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/m7s.Global/Shutdown", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *globalClient) Restart(ctx context.Context, in *RequestWithId, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/m7s.Global/Restart", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *globalClient) StreamList(ctx context.Context, in *StreamListRequest, opts ...grpc.CallOption) (*StreamListResponse, error) {
	out := new(StreamListResponse)
	err := c.cc.Invoke(ctx, "/m7s.Global/StreamList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *globalClient) WaitList(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StreamWaitListResponse, error) {
	out := new(StreamWaitListResponse)
	err := c.cc.Invoke(ctx, "/m7s.Global/WaitList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *globalClient) StreamInfo(ctx context.Context, in *StreamSnapRequest, opts ...grpc.CallOption) (*StreamInfoResponse, error) {
	out := new(StreamInfoResponse)
	err := c.cc.Invoke(ctx, "/m7s.Global/StreamInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *globalClient) GetSubscribers(ctx context.Context, in *SubscribersRequest, opts ...grpc.CallOption) (*SubscribersResponse, error) {
	out := new(SubscribersResponse)
	err := c.cc.Invoke(ctx, "/m7s.Global/GetSubscribers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *globalClient) AudioTrackSnap(ctx context.Context, in *StreamSnapRequest, opts ...grpc.CallOption) (*TrackSnapShotResponse, error) {
	out := new(TrackSnapShotResponse)
	err := c.cc.Invoke(ctx, "/m7s.Global/AudioTrackSnap", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *globalClient) VideoTrackSnap(ctx context.Context, in *StreamSnapRequest, opts ...grpc.CallOption) (*TrackSnapShotResponse, error) {
	out := new(TrackSnapShotResponse)
	err := c.cc.Invoke(ctx, "/m7s.Global/VideoTrackSnap", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *globalClient) ChangeSubscribe(ctx context.Context, in *ChangeSubscribeRequest, opts ...grpc.CallOption) (*SuccessResponse, error) {
	out := new(SuccessResponse)
	err := c.cc.Invoke(ctx, "/m7s.Global/ChangeSubscribe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *globalClient) StopSubscribe(ctx context.Context, in *RequestWithId, opts ...grpc.CallOption) (*SuccessResponse, error) {
	out := new(SuccessResponse)
	err := c.cc.Invoke(ctx, "/m7s.Global/StopSubscribe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *globalClient) GetConfig(ctx context.Context, in *GetConfigRequest, opts ...grpc.CallOption) (*GetConfigResponse, error) {
	out := new(GetConfigResponse)
	err := c.cc.Invoke(ctx, "/m7s.Global/GetConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *globalClient) GetFormily(ctx context.Context, in *GetConfigRequest, opts ...grpc.CallOption) (*GetConfigResponse, error) {
	out := new(GetConfigResponse)
	err := c.cc.Invoke(ctx, "/m7s.Global/GetFormily", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *globalClient) ModifyConfig(ctx context.Context, in *ModifyConfigRequest, opts ...grpc.CallOption) (*SuccessResponse, error) {
	out := new(SuccessResponse)
	err := c.cc.Invoke(ctx, "/m7s.Global/ModifyConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GlobalServer is the server API for Global service.
// All implementations must embed UnimplementedGlobalServer
// for forward compatibility
type GlobalServer interface {
	SysInfo(context.Context, *emptypb.Empty) (*SysInfoResponse, error)
	Summary(context.Context, *emptypb.Empty) (*SummaryResponse, error)
	Shutdown(context.Context, *RequestWithId) (*emptypb.Empty, error)
	Restart(context.Context, *RequestWithId) (*emptypb.Empty, error)
	StreamList(context.Context, *StreamListRequest) (*StreamListResponse, error)
	WaitList(context.Context, *emptypb.Empty) (*StreamWaitListResponse, error)
	StreamInfo(context.Context, *StreamSnapRequest) (*StreamInfoResponse, error)
	GetSubscribers(context.Context, *SubscribersRequest) (*SubscribersResponse, error)
	AudioTrackSnap(context.Context, *StreamSnapRequest) (*TrackSnapShotResponse, error)
	VideoTrackSnap(context.Context, *StreamSnapRequest) (*TrackSnapShotResponse, error)
	ChangeSubscribe(context.Context, *ChangeSubscribeRequest) (*SuccessResponse, error)
	StopSubscribe(context.Context, *RequestWithId) (*SuccessResponse, error)
	GetConfig(context.Context, *GetConfigRequest) (*GetConfigResponse, error)
	GetFormily(context.Context, *GetConfigRequest) (*GetConfigResponse, error)
	ModifyConfig(context.Context, *ModifyConfigRequest) (*SuccessResponse, error)
	mustEmbedUnimplementedGlobalServer()
}

// UnimplementedGlobalServer must be embedded to have forward compatible implementations.
type UnimplementedGlobalServer struct {
}

func (UnimplementedGlobalServer) SysInfo(context.Context, *emptypb.Empty) (*SysInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SysInfo not implemented")
}
func (UnimplementedGlobalServer) Summary(context.Context, *emptypb.Empty) (*SummaryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Summary not implemented")
}
func (UnimplementedGlobalServer) Shutdown(context.Context, *RequestWithId) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Shutdown not implemented")
}
func (UnimplementedGlobalServer) Restart(context.Context, *RequestWithId) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Restart not implemented")
}
func (UnimplementedGlobalServer) StreamList(context.Context, *StreamListRequest) (*StreamListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StreamList not implemented")
}
func (UnimplementedGlobalServer) WaitList(context.Context, *emptypb.Empty) (*StreamWaitListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WaitList not implemented")
}
func (UnimplementedGlobalServer) StreamInfo(context.Context, *StreamSnapRequest) (*StreamInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StreamInfo not implemented")
}
func (UnimplementedGlobalServer) GetSubscribers(context.Context, *SubscribersRequest) (*SubscribersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSubscribers not implemented")
}
func (UnimplementedGlobalServer) AudioTrackSnap(context.Context, *StreamSnapRequest) (*TrackSnapShotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AudioTrackSnap not implemented")
}
func (UnimplementedGlobalServer) VideoTrackSnap(context.Context, *StreamSnapRequest) (*TrackSnapShotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VideoTrackSnap not implemented")
}
func (UnimplementedGlobalServer) ChangeSubscribe(context.Context, *ChangeSubscribeRequest) (*SuccessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeSubscribe not implemented")
}
func (UnimplementedGlobalServer) StopSubscribe(context.Context, *RequestWithId) (*SuccessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopSubscribe not implemented")
}
func (UnimplementedGlobalServer) GetConfig(context.Context, *GetConfigRequest) (*GetConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfig not implemented")
}
func (UnimplementedGlobalServer) GetFormily(context.Context, *GetConfigRequest) (*GetConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFormily not implemented")
}
func (UnimplementedGlobalServer) ModifyConfig(context.Context, *ModifyConfigRequest) (*SuccessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModifyConfig not implemented")
}
func (UnimplementedGlobalServer) mustEmbedUnimplementedGlobalServer() {}

// UnsafeGlobalServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GlobalServer will
// result in compilation errors.
type UnsafeGlobalServer interface {
	mustEmbedUnimplementedGlobalServer()
}

func RegisterGlobalServer(s grpc.ServiceRegistrar, srv GlobalServer) {
	s.RegisterService(&Global_ServiceDesc, srv)
}

func _Global_SysInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GlobalServer).SysInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m7s.Global/SysInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GlobalServer).SysInfo(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Global_Summary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GlobalServer).Summary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m7s.Global/Summary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GlobalServer).Summary(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Global_Shutdown_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestWithId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GlobalServer).Shutdown(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m7s.Global/Shutdown",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GlobalServer).Shutdown(ctx, req.(*RequestWithId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Global_Restart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestWithId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GlobalServer).Restart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m7s.Global/Restart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GlobalServer).Restart(ctx, req.(*RequestWithId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Global_StreamList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StreamListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GlobalServer).StreamList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m7s.Global/StreamList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GlobalServer).StreamList(ctx, req.(*StreamListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Global_WaitList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GlobalServer).WaitList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m7s.Global/WaitList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GlobalServer).WaitList(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Global_StreamInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StreamSnapRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GlobalServer).StreamInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m7s.Global/StreamInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GlobalServer).StreamInfo(ctx, req.(*StreamSnapRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Global_GetSubscribers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubscribersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GlobalServer).GetSubscribers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m7s.Global/GetSubscribers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GlobalServer).GetSubscribers(ctx, req.(*SubscribersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Global_AudioTrackSnap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StreamSnapRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GlobalServer).AudioTrackSnap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m7s.Global/AudioTrackSnap",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GlobalServer).AudioTrackSnap(ctx, req.(*StreamSnapRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Global_VideoTrackSnap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StreamSnapRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GlobalServer).VideoTrackSnap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m7s.Global/VideoTrackSnap",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GlobalServer).VideoTrackSnap(ctx, req.(*StreamSnapRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Global_ChangeSubscribe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeSubscribeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GlobalServer).ChangeSubscribe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m7s.Global/ChangeSubscribe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GlobalServer).ChangeSubscribe(ctx, req.(*ChangeSubscribeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Global_StopSubscribe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestWithId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GlobalServer).StopSubscribe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m7s.Global/StopSubscribe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GlobalServer).StopSubscribe(ctx, req.(*RequestWithId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Global_GetConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GlobalServer).GetConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m7s.Global/GetConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GlobalServer).GetConfig(ctx, req.(*GetConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Global_GetFormily_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GlobalServer).GetFormily(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m7s.Global/GetFormily",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GlobalServer).GetFormily(ctx, req.(*GetConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Global_ModifyConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModifyConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GlobalServer).ModifyConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m7s.Global/ModifyConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GlobalServer).ModifyConfig(ctx, req.(*ModifyConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Global_ServiceDesc is the grpc.ServiceDesc for Global service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Global_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "m7s.Global",
	HandlerType: (*GlobalServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SysInfo",
			Handler:    _Global_SysInfo_Handler,
		},
		{
			MethodName: "Summary",
			Handler:    _Global_Summary_Handler,
		},
		{
			MethodName: "Shutdown",
			Handler:    _Global_Shutdown_Handler,
		},
		{
			MethodName: "Restart",
			Handler:    _Global_Restart_Handler,
		},
		{
			MethodName: "StreamList",
			Handler:    _Global_StreamList_Handler,
		},
		{
			MethodName: "WaitList",
			Handler:    _Global_WaitList_Handler,
		},
		{
			MethodName: "StreamInfo",
			Handler:    _Global_StreamInfo_Handler,
		},
		{
			MethodName: "GetSubscribers",
			Handler:    _Global_GetSubscribers_Handler,
		},
		{
			MethodName: "AudioTrackSnap",
			Handler:    _Global_AudioTrackSnap_Handler,
		},
		{
			MethodName: "VideoTrackSnap",
			Handler:    _Global_VideoTrackSnap_Handler,
		},
		{
			MethodName: "ChangeSubscribe",
			Handler:    _Global_ChangeSubscribe_Handler,
		},
		{
			MethodName: "StopSubscribe",
			Handler:    _Global_StopSubscribe_Handler,
		},
		{
			MethodName: "GetConfig",
			Handler:    _Global_GetConfig_Handler,
		},
		{
			MethodName: "GetFormily",
			Handler:    _Global_GetFormily_Handler,
		},
		{
			MethodName: "ModifyConfig",
			Handler:    _Global_ModifyConfig_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "global.proto",
}
