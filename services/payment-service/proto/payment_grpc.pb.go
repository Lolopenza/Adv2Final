// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.2
// source: proto/payment.proto

package proto

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
	PaymentService_CreatePayment_FullMethodName       = "/payment.PaymentService/CreatePayment"
	PaymentService_GetPayment_FullMethodName          = "/payment.PaymentService/GetPayment"
	PaymentService_ListPayments_FullMethodName        = "/payment.PaymentService/ListPayments"
	PaymentService_ConfirmPayment_FullMethodName      = "/payment.PaymentService/ConfirmPayment"
	PaymentService_RefundPayment_FullMethodName       = "/payment.PaymentService/RefundPayment"
	PaymentService_CancelPayment_FullMethodName       = "/payment.PaymentService/CancelPayment"
	PaymentService_GenerateInvoice_FullMethodName     = "/payment.PaymentService/GenerateInvoice"
	PaymentService_SendPaymentReminder_FullMethodName = "/payment.PaymentService/SendPaymentReminder"
	PaymentService_DeletePayment_FullMethodName       = "/payment.PaymentService/DeletePayment"
	PaymentService_UpdatePayment_FullMethodName       = "/payment.PaymentService/UpdatePayment"
)

// PaymentServiceClient is the client API for PaymentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PaymentServiceClient interface {
	CreatePayment(ctx context.Context, in *CreatePaymentRequest, opts ...grpc.CallOption) (*CreatePaymentResponse, error)
	GetPayment(ctx context.Context, in *GetPaymentRequest, opts ...grpc.CallOption) (*GetPaymentResponse, error)
	ListPayments(ctx context.Context, in *ListPaymentsRequest, opts ...grpc.CallOption) (*ListPaymentsResponse, error)
	ConfirmPayment(ctx context.Context, in *ConfirmPaymentRequest, opts ...grpc.CallOption) (*ConfirmPaymentResponse, error)
	RefundPayment(ctx context.Context, in *RefundPaymentRequest, opts ...grpc.CallOption) (*RefundPaymentResponse, error)
	CancelPayment(ctx context.Context, in *CancelPaymentRequest, opts ...grpc.CallOption) (*CancelPaymentResponse, error)
	GenerateInvoice(ctx context.Context, in *GenerateInvoiceRequest, opts ...grpc.CallOption) (*GenerateInvoiceResponse, error)
	SendPaymentReminder(ctx context.Context, in *SendPaymentReminderRequest, opts ...grpc.CallOption) (*SendPaymentReminderResponse, error)
	DeletePayment(ctx context.Context, in *DeletePaymentRequest, opts ...grpc.CallOption) (*DeletePaymentResponse, error)
	UpdatePayment(ctx context.Context, in *UpdatePaymentRequest, opts ...grpc.CallOption) (*UpdatePaymentResponse, error)
}

type paymentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPaymentServiceClient(cc grpc.ClientConnInterface) PaymentServiceClient {
	return &paymentServiceClient{cc}
}

func (c *paymentServiceClient) CreatePayment(ctx context.Context, in *CreatePaymentRequest, opts ...grpc.CallOption) (*CreatePaymentResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreatePaymentResponse)
	err := c.cc.Invoke(ctx, PaymentService_CreatePayment_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentServiceClient) GetPayment(ctx context.Context, in *GetPaymentRequest, opts ...grpc.CallOption) (*GetPaymentResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPaymentResponse)
	err := c.cc.Invoke(ctx, PaymentService_GetPayment_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentServiceClient) ListPayments(ctx context.Context, in *ListPaymentsRequest, opts ...grpc.CallOption) (*ListPaymentsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListPaymentsResponse)
	err := c.cc.Invoke(ctx, PaymentService_ListPayments_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentServiceClient) ConfirmPayment(ctx context.Context, in *ConfirmPaymentRequest, opts ...grpc.CallOption) (*ConfirmPaymentResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ConfirmPaymentResponse)
	err := c.cc.Invoke(ctx, PaymentService_ConfirmPayment_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentServiceClient) RefundPayment(ctx context.Context, in *RefundPaymentRequest, opts ...grpc.CallOption) (*RefundPaymentResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RefundPaymentResponse)
	err := c.cc.Invoke(ctx, PaymentService_RefundPayment_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentServiceClient) CancelPayment(ctx context.Context, in *CancelPaymentRequest, opts ...grpc.CallOption) (*CancelPaymentResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CancelPaymentResponse)
	err := c.cc.Invoke(ctx, PaymentService_CancelPayment_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentServiceClient) GenerateInvoice(ctx context.Context, in *GenerateInvoiceRequest, opts ...grpc.CallOption) (*GenerateInvoiceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GenerateInvoiceResponse)
	err := c.cc.Invoke(ctx, PaymentService_GenerateInvoice_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentServiceClient) SendPaymentReminder(ctx context.Context, in *SendPaymentReminderRequest, opts ...grpc.CallOption) (*SendPaymentReminderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SendPaymentReminderResponse)
	err := c.cc.Invoke(ctx, PaymentService_SendPaymentReminder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentServiceClient) DeletePayment(ctx context.Context, in *DeletePaymentRequest, opts ...grpc.CallOption) (*DeletePaymentResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeletePaymentResponse)
	err := c.cc.Invoke(ctx, PaymentService_DeletePayment_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentServiceClient) UpdatePayment(ctx context.Context, in *UpdatePaymentRequest, opts ...grpc.CallOption) (*UpdatePaymentResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdatePaymentResponse)
	err := c.cc.Invoke(ctx, PaymentService_UpdatePayment_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PaymentServiceServer is the server API for PaymentService service.
// All implementations must embed UnimplementedPaymentServiceServer
// for forward compatibility.
type PaymentServiceServer interface {
	CreatePayment(context.Context, *CreatePaymentRequest) (*CreatePaymentResponse, error)
	GetPayment(context.Context, *GetPaymentRequest) (*GetPaymentResponse, error)
	ListPayments(context.Context, *ListPaymentsRequest) (*ListPaymentsResponse, error)
	ConfirmPayment(context.Context, *ConfirmPaymentRequest) (*ConfirmPaymentResponse, error)
	RefundPayment(context.Context, *RefundPaymentRequest) (*RefundPaymentResponse, error)
	CancelPayment(context.Context, *CancelPaymentRequest) (*CancelPaymentResponse, error)
	GenerateInvoice(context.Context, *GenerateInvoiceRequest) (*GenerateInvoiceResponse, error)
	SendPaymentReminder(context.Context, *SendPaymentReminderRequest) (*SendPaymentReminderResponse, error)
	DeletePayment(context.Context, *DeletePaymentRequest) (*DeletePaymentResponse, error)
	UpdatePayment(context.Context, *UpdatePaymentRequest) (*UpdatePaymentResponse, error)
	mustEmbedUnimplementedPaymentServiceServer()
}

// UnimplementedPaymentServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPaymentServiceServer struct{}

func (UnimplementedPaymentServiceServer) CreatePayment(context.Context, *CreatePaymentRequest) (*CreatePaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePayment not implemented")
}
func (UnimplementedPaymentServiceServer) GetPayment(context.Context, *GetPaymentRequest) (*GetPaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPayment not implemented")
}
func (UnimplementedPaymentServiceServer) ListPayments(context.Context, *ListPaymentsRequest) (*ListPaymentsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPayments not implemented")
}
func (UnimplementedPaymentServiceServer) ConfirmPayment(context.Context, *ConfirmPaymentRequest) (*ConfirmPaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfirmPayment not implemented")
}
func (UnimplementedPaymentServiceServer) RefundPayment(context.Context, *RefundPaymentRequest) (*RefundPaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefundPayment not implemented")
}
func (UnimplementedPaymentServiceServer) CancelPayment(context.Context, *CancelPaymentRequest) (*CancelPaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelPayment not implemented")
}
func (UnimplementedPaymentServiceServer) GenerateInvoice(context.Context, *GenerateInvoiceRequest) (*GenerateInvoiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateInvoice not implemented")
}
func (UnimplementedPaymentServiceServer) SendPaymentReminder(context.Context, *SendPaymentReminderRequest) (*SendPaymentReminderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendPaymentReminder not implemented")
}
func (UnimplementedPaymentServiceServer) DeletePayment(context.Context, *DeletePaymentRequest) (*DeletePaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePayment not implemented")
}
func (UnimplementedPaymentServiceServer) UpdatePayment(context.Context, *UpdatePaymentRequest) (*UpdatePaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePayment not implemented")
}
func (UnimplementedPaymentServiceServer) mustEmbedUnimplementedPaymentServiceServer() {}
func (UnimplementedPaymentServiceServer) testEmbeddedByValue()                        {}

// UnsafePaymentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PaymentServiceServer will
// result in compilation errors.
type UnsafePaymentServiceServer interface {
	mustEmbedUnimplementedPaymentServiceServer()
}

func RegisterPaymentServiceServer(s grpc.ServiceRegistrar, srv PaymentServiceServer) {
	// If the following call pancis, it indicates UnimplementedPaymentServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PaymentService_ServiceDesc, srv)
}

func _PaymentService_CreatePayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePaymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).CreatePayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentService_CreatePayment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServiceServer).CreatePayment(ctx, req.(*CreatePaymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentService_GetPayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPaymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).GetPayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentService_GetPayment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServiceServer).GetPayment(ctx, req.(*GetPaymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentService_ListPayments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPaymentsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).ListPayments(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentService_ListPayments_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServiceServer).ListPayments(ctx, req.(*ListPaymentsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentService_ConfirmPayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfirmPaymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).ConfirmPayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentService_ConfirmPayment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServiceServer).ConfirmPayment(ctx, req.(*ConfirmPaymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentService_RefundPayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefundPaymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).RefundPayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentService_RefundPayment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServiceServer).RefundPayment(ctx, req.(*RefundPaymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentService_CancelPayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelPaymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).CancelPayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentService_CancelPayment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServiceServer).CancelPayment(ctx, req.(*CancelPaymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentService_GenerateInvoice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateInvoiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).GenerateInvoice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentService_GenerateInvoice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServiceServer).GenerateInvoice(ctx, req.(*GenerateInvoiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentService_SendPaymentReminder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendPaymentReminderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).SendPaymentReminder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentService_SendPaymentReminder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServiceServer).SendPaymentReminder(ctx, req.(*SendPaymentReminderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentService_DeletePayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePaymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).DeletePayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentService_DeletePayment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServiceServer).DeletePayment(ctx, req.(*DeletePaymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentService_UpdatePayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePaymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).UpdatePayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentService_UpdatePayment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServiceServer).UpdatePayment(ctx, req.(*UpdatePaymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PaymentService_ServiceDesc is the grpc.ServiceDesc for PaymentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PaymentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "payment.PaymentService",
	HandlerType: (*PaymentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePayment",
			Handler:    _PaymentService_CreatePayment_Handler,
		},
		{
			MethodName: "GetPayment",
			Handler:    _PaymentService_GetPayment_Handler,
		},
		{
			MethodName: "ListPayments",
			Handler:    _PaymentService_ListPayments_Handler,
		},
		{
			MethodName: "ConfirmPayment",
			Handler:    _PaymentService_ConfirmPayment_Handler,
		},
		{
			MethodName: "RefundPayment",
			Handler:    _PaymentService_RefundPayment_Handler,
		},
		{
			MethodName: "CancelPayment",
			Handler:    _PaymentService_CancelPayment_Handler,
		},
		{
			MethodName: "GenerateInvoice",
			Handler:    _PaymentService_GenerateInvoice_Handler,
		},
		{
			MethodName: "SendPaymentReminder",
			Handler:    _PaymentService_SendPaymentReminder_Handler,
		},
		{
			MethodName: "DeletePayment",
			Handler:    _PaymentService_DeletePayment_Handler,
		},
		{
			MethodName: "UpdatePayment",
			Handler:    _PaymentService_UpdatePayment_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/payment.proto",
}
