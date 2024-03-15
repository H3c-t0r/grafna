// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.2
// source: entity.proto

package entity

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	EntityStore_Read_FullMethodName      = "/entity.EntityStore/Read"
	EntityStore_BatchRead_FullMethodName = "/entity.EntityStore/BatchRead"
	EntityStore_Create_FullMethodName    = "/entity.EntityStore/Create"
	EntityStore_Update_FullMethodName    = "/entity.EntityStore/Update"
	EntityStore_Delete_FullMethodName    = "/entity.EntityStore/Delete"
	EntityStore_History_FullMethodName   = "/entity.EntityStore/History"
	EntityStore_List_FullMethodName      = "/entity.EntityStore/List"
	EntityStore_Watch_FullMethodName     = "/entity.EntityStore/Watch"
)

// EntityStoreClient is the client API for EntityStore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EntityStoreClient interface {
	Read(ctx context.Context, in *ReadEntityRequest, opts ...grpc.CallOption) (*Entity, error)
	BatchRead(ctx context.Context, in *BatchReadEntityRequest, opts ...grpc.CallOption) (*BatchReadEntityResponse, error)
	Create(ctx context.Context, in *CreateEntityRequest, opts ...grpc.CallOption) (*CreateEntityResponse, error)
	Update(ctx context.Context, in *UpdateEntityRequest, opts ...grpc.CallOption) (*UpdateEntityResponse, error)
	Delete(ctx context.Context, in *DeleteEntityRequest, opts ...grpc.CallOption) (*DeleteEntityResponse, error)
	History(ctx context.Context, in *EntityHistoryRequest, opts ...grpc.CallOption) (*EntityHistoryResponse, error)
	List(ctx context.Context, in *EntityListRequest, opts ...grpc.CallOption) (*EntityListResponse, error)
	Watch(ctx context.Context, in *EntityWatchRequest, opts ...grpc.CallOption) (EntityStore_WatchClient, error)
}

type entityStoreClient struct {
	cc grpc.ClientConnInterface
}

func NewEntityStoreClient(cc grpc.ClientConnInterface) EntityStoreClient {
	return &entityStoreClient{cc}
}

func (c *entityStoreClient) Read(ctx context.Context, in *ReadEntityRequest, opts ...grpc.CallOption) (*Entity, error) {
	out := new(Entity)
	err := c.cc.Invoke(ctx, EntityStore_Read_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *entityStoreClient) BatchRead(ctx context.Context, in *BatchReadEntityRequest, opts ...grpc.CallOption) (*BatchReadEntityResponse, error) {
	out := new(BatchReadEntityResponse)
	err := c.cc.Invoke(ctx, EntityStore_BatchRead_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *entityStoreClient) Create(ctx context.Context, in *CreateEntityRequest, opts ...grpc.CallOption) (*CreateEntityResponse, error) {
	out := new(CreateEntityResponse)
	err := c.cc.Invoke(ctx, EntityStore_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *entityStoreClient) Update(ctx context.Context, in *UpdateEntityRequest, opts ...grpc.CallOption) (*UpdateEntityResponse, error) {
	out := new(UpdateEntityResponse)
	err := c.cc.Invoke(ctx, EntityStore_Update_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *entityStoreClient) Delete(ctx context.Context, in *DeleteEntityRequest, opts ...grpc.CallOption) (*DeleteEntityResponse, error) {
	out := new(DeleteEntityResponse)
	err := c.cc.Invoke(ctx, EntityStore_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *entityStoreClient) History(ctx context.Context, in *EntityHistoryRequest, opts ...grpc.CallOption) (*EntityHistoryResponse, error) {
	out := new(EntityHistoryResponse)
	err := c.cc.Invoke(ctx, EntityStore_History_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *entityStoreClient) List(ctx context.Context, in *EntityListRequest, opts ...grpc.CallOption) (*EntityListResponse, error) {
	out := new(EntityListResponse)
	err := c.cc.Invoke(ctx, EntityStore_List_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *entityStoreClient) Watch(ctx context.Context, in *EntityWatchRequest, opts ...grpc.CallOption) (EntityStore_WatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &EntityStore_ServiceDesc.Streams[0], EntityStore_Watch_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &entityStoreWatchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type EntityStore_WatchClient interface {
	Recv() (*EntityWatchResponse, error)
	grpc.ClientStream
}

type entityStoreWatchClient struct {
	grpc.ClientStream
}

func (x *entityStoreWatchClient) Recv() (*EntityWatchResponse, error) {
	m := new(EntityWatchResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// EntityStoreServer is the server API for EntityStore service.
// All implementations should embed UnimplementedEntityStoreServer
// for forward compatibility
type EntityStoreServer interface {
	Read(context.Context, *ReadEntityRequest) (*Entity, error)
	BatchRead(context.Context, *BatchReadEntityRequest) (*BatchReadEntityResponse, error)
	Create(context.Context, *CreateEntityRequest) (*CreateEntityResponse, error)
	Update(context.Context, *UpdateEntityRequest) (*UpdateEntityResponse, error)
	Delete(context.Context, *DeleteEntityRequest) (*DeleteEntityResponse, error)
	History(context.Context, *EntityHistoryRequest) (*EntityHistoryResponse, error)
	List(context.Context, *EntityListRequest) (*EntityListResponse, error)
	Watch(*EntityWatchRequest, EntityStore_WatchServer) error
}

// UnimplementedEntityStoreServer should be embedded to have forward compatible implementations.
type UnimplementedEntityStoreServer struct {
}

func (UnimplementedEntityStoreServer) Read(context.Context, *ReadEntityRequest) (*Entity, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedEntityStoreServer) BatchRead(context.Context, *BatchReadEntityRequest) (*BatchReadEntityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchRead not implemented")
}
func (UnimplementedEntityStoreServer) Create(context.Context, *CreateEntityRequest) (*CreateEntityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedEntityStoreServer) Update(context.Context, *UpdateEntityRequest) (*UpdateEntityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedEntityStoreServer) Delete(context.Context, *DeleteEntityRequest) (*DeleteEntityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedEntityStoreServer) History(context.Context, *EntityHistoryRequest) (*EntityHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method History not implemented")
}
func (UnimplementedEntityStoreServer) List(context.Context, *EntityListRequest) (*EntityListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedEntityStoreServer) Watch(*EntityWatchRequest, EntityStore_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}

// UnsafeEntityStoreServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EntityStoreServer will
// result in compilation errors.
type UnsafeEntityStoreServer interface {
	mustEmbedUnimplementedEntityStoreServer()
}

func RegisterEntityStoreServer(s grpc.ServiceRegistrar, srv EntityStoreServer) {
	s.RegisterService(&EntityStore_ServiceDesc, srv)
}

func _EntityStore_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadEntityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EntityStoreServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EntityStore_Read_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EntityStoreServer).Read(ctx, req.(*ReadEntityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EntityStore_BatchRead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchReadEntityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EntityStoreServer).BatchRead(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EntityStore_BatchRead_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EntityStoreServer).BatchRead(ctx, req.(*BatchReadEntityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EntityStore_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateEntityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EntityStoreServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EntityStore_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EntityStoreServer).Create(ctx, req.(*CreateEntityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EntityStore_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateEntityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EntityStoreServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EntityStore_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EntityStoreServer).Update(ctx, req.(*UpdateEntityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EntityStore_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteEntityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EntityStoreServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EntityStore_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EntityStoreServer).Delete(ctx, req.(*DeleteEntityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EntityStore_History_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EntityHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EntityStoreServer).History(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EntityStore_History_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EntityStoreServer).History(ctx, req.(*EntityHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EntityStore_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EntityListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EntityStoreServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EntityStore_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EntityStoreServer).List(ctx, req.(*EntityListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EntityStore_Watch_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EntityWatchRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EntityStoreServer).Watch(m, &entityStoreWatchServer{stream})
}

type EntityStore_WatchServer interface {
	Send(*EntityWatchResponse) error
	grpc.ServerStream
}

type entityStoreWatchServer struct {
	grpc.ServerStream
}

func (x *entityStoreWatchServer) Send(m *EntityWatchResponse) error {
	return x.ServerStream.SendMsg(m)
}

// EntityStore_ServiceDesc is the grpc.ServiceDesc for EntityStore service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EntityStore_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "entity.EntityStore",
	HandlerType: (*EntityStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Read",
			Handler:    _EntityStore_Read_Handler,
		},
		{
			MethodName: "BatchRead",
			Handler:    _EntityStore_BatchRead_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _EntityStore_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _EntityStore_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _EntityStore_Delete_Handler,
		},
		{
			MethodName: "History",
			Handler:    _EntityStore_History_Handler,
		},
		{
			MethodName: "List",
			Handler:    _EntityStore_List_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Watch",
			Handler:       _EntityStore_Watch_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "entity.proto",
}
