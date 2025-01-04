// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v4.25.2
// source: api/v1/transaction.proto

package ebank

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
	TransactionService_Deposit_FullMethodName               = "/proto.TransactionService/Deposit"
	TransactionService_Withdraw_FullMethodName              = "/proto.TransactionService/Withdraw"
	TransactionService_GetTransactionHistory_FullMethodName = "/proto.TransactionService/GetTransactionHistory"
)

// TransactionServiceClient is the client API for TransactionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TransactionServiceClient interface {
	// 입금/출금
	Deposit(ctx context.Context, in *DepositRequest, opts ...grpc.CallOption) (*TransactionResponse, error)
	Withdraw(ctx context.Context, in *WithdrawRequest, opts ...grpc.CallOption) (*TransactionResponse, error)
	// 거래 내역 조회
	GetTransactionHistory(ctx context.Context, in *GetTransactionHistoryRequest, opts ...grpc.CallOption) (*GetTransactionHistoryResponse, error)
}

type transactionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTransactionServiceClient(cc grpc.ClientConnInterface) TransactionServiceClient {
	return &transactionServiceClient{cc}
}

func (c *transactionServiceClient) Deposit(ctx context.Context, in *DepositRequest, opts ...grpc.CallOption) (*TransactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TransactionResponse)
	err := c.cc.Invoke(ctx, TransactionService_Deposit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionServiceClient) Withdraw(ctx context.Context, in *WithdrawRequest, opts ...grpc.CallOption) (*TransactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TransactionResponse)
	err := c.cc.Invoke(ctx, TransactionService_Withdraw_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionServiceClient) GetTransactionHistory(ctx context.Context, in *GetTransactionHistoryRequest, opts ...grpc.CallOption) (*GetTransactionHistoryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTransactionHistoryResponse)
	err := c.cc.Invoke(ctx, TransactionService_GetTransactionHistory_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TransactionServiceServer is the server API for TransactionService service.
// All implementations must embed UnimplementedTransactionServiceServer
// for forward compatibility.
type TransactionServiceServer interface {
	// 입금/출금
	Deposit(context.Context, *DepositRequest) (*TransactionResponse, error)
	Withdraw(context.Context, *WithdrawRequest) (*TransactionResponse, error)
	// 거래 내역 조회
	GetTransactionHistory(context.Context, *GetTransactionHistoryRequest) (*GetTransactionHistoryResponse, error)
	mustEmbedUnimplementedTransactionServiceServer()
}

// UnimplementedTransactionServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTransactionServiceServer struct{}

func (UnimplementedTransactionServiceServer) Deposit(context.Context, *DepositRequest) (*TransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deposit not implemented")
}
func (UnimplementedTransactionServiceServer) Withdraw(context.Context, *WithdrawRequest) (*TransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Withdraw not implemented")
}
func (UnimplementedTransactionServiceServer) GetTransactionHistory(context.Context, *GetTransactionHistoryRequest) (*GetTransactionHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransactionHistory not implemented")
}
func (UnimplementedTransactionServiceServer) mustEmbedUnimplementedTransactionServiceServer() {}
func (UnimplementedTransactionServiceServer) testEmbeddedByValue()                            {}

// UnsafeTransactionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TransactionServiceServer will
// result in compilation errors.
type UnsafeTransactionServiceServer interface {
	mustEmbedUnimplementedTransactionServiceServer()
}

func RegisterTransactionServiceServer(s grpc.ServiceRegistrar, srv TransactionServiceServer) {
	// If the following call pancis, it indicates UnimplementedTransactionServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&TransactionService_ServiceDesc, srv)
}

func _TransactionService_Deposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DepositRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServiceServer).Deposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TransactionService_Deposit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServiceServer).Deposit(ctx, req.(*DepositRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TransactionService_Withdraw_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WithdrawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServiceServer).Withdraw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TransactionService_Withdraw_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServiceServer).Withdraw(ctx, req.(*WithdrawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TransactionService_GetTransactionHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTransactionHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServiceServer).GetTransactionHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TransactionService_GetTransactionHistory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServiceServer).GetTransactionHistory(ctx, req.(*GetTransactionHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TransactionService_ServiceDesc is the grpc.ServiceDesc for TransactionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TransactionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.TransactionService",
	HandlerType: (*TransactionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Deposit",
			Handler:    _TransactionService_Deposit_Handler,
		},
		{
			MethodName: "Withdraw",
			Handler:    _TransactionService_Withdraw_Handler,
		},
		{
			MethodName: "GetTransactionHistory",
			Handler:    _TransactionService_GetTransactionHistory_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v1/transaction.proto",
}
