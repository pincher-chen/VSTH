// Code generated by protoc-gen-go.
// source: github.com/GoogleCloudPlatform/gcloud-golang/bigtable/internal/table_service_proto/bigtable_table_service.proto
// DO NOT EDIT!

package google_bigtable_admin_table_v1

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_bigtable_admin_table_v11 "github.com/GoogleCloudPlatform/gcloud-golang/bigtable/internal/table_data_proto"
import google_protobuf1 "github.com/golang/protobuf/ptypes/empty"

import (
	context "github.com/golang/net/context"
	grpc "github.com/grpc/grpc-go"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion2

// Client API for BigtableTableService service

type BigtableTableServiceClient interface {
	// Creates a new table, to be served from a specified cluster.
	// The table can be created with a full set of initial column families,
	// specified in the request.
	CreateTable(ctx context.Context, in *CreateTableRequest, opts ...grpc.CallOption) (*google_bigtable_admin_table_v11.Table, error)
	// Lists the names of all tables served from a specified cluster.
	ListTables(ctx context.Context, in *ListTablesRequest, opts ...grpc.CallOption) (*ListTablesResponse, error)
	// Gets the schema of the specified table, including its column families.
	GetTable(ctx context.Context, in *GetTableRequest, opts ...grpc.CallOption) (*google_bigtable_admin_table_v11.Table, error)
	// Permanently deletes a specified table and all of its data.
	DeleteTable(ctx context.Context, in *DeleteTableRequest, opts ...grpc.CallOption) (*google_protobuf1.Empty, error)
	// Changes the name of a specified table.
	// Cannot be used to move tables between clusters, zones, or projects.
	RenameTable(ctx context.Context, in *RenameTableRequest, opts ...grpc.CallOption) (*google_protobuf1.Empty, error)
	// Creates a new column family within a specified table.
	CreateColumnFamily(ctx context.Context, in *CreateColumnFamilyRequest, opts ...grpc.CallOption) (*google_bigtable_admin_table_v11.ColumnFamily, error)
	// Changes the configuration of a specified column family.
	UpdateColumnFamily(ctx context.Context, in *google_bigtable_admin_table_v11.ColumnFamily, opts ...grpc.CallOption) (*google_bigtable_admin_table_v11.ColumnFamily, error)
	// Permanently deletes a specified column family and all of its data.
	DeleteColumnFamily(ctx context.Context, in *DeleteColumnFamilyRequest, opts ...grpc.CallOption) (*google_protobuf1.Empty, error)
	// Delete all rows in a table corresponding to a particular prefix
	BulkDeleteRows(ctx context.Context, in *BulkDeleteRowsRequest, opts ...grpc.CallOption) (*google_protobuf1.Empty, error)
}

type bigtableTableServiceClient struct {
	cc *grpc.ClientConn
}

func NewBigtableTableServiceClient(cc *grpc.ClientConn) BigtableTableServiceClient {
	return &bigtableTableServiceClient{cc}
}

func (c *bigtableTableServiceClient) CreateTable(ctx context.Context, in *CreateTableRequest, opts ...grpc.CallOption) (*google_bigtable_admin_table_v11.Table, error) {
	out := new(google_bigtable_admin_table_v11.Table)
	err := grpc.Invoke(ctx, "/google.bigtable.admin.table.v1.BigtableTableService/CreateTable", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bigtableTableServiceClient) ListTables(ctx context.Context, in *ListTablesRequest, opts ...grpc.CallOption) (*ListTablesResponse, error) {
	out := new(ListTablesResponse)
	err := grpc.Invoke(ctx, "/google.bigtable.admin.table.v1.BigtableTableService/ListTables", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bigtableTableServiceClient) GetTable(ctx context.Context, in *GetTableRequest, opts ...grpc.CallOption) (*google_bigtable_admin_table_v11.Table, error) {
	out := new(google_bigtable_admin_table_v11.Table)
	err := grpc.Invoke(ctx, "/google.bigtable.admin.table.v1.BigtableTableService/GetTable", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bigtableTableServiceClient) DeleteTable(ctx context.Context, in *DeleteTableRequest, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := grpc.Invoke(ctx, "/google.bigtable.admin.table.v1.BigtableTableService/DeleteTable", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bigtableTableServiceClient) RenameTable(ctx context.Context, in *RenameTableRequest, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := grpc.Invoke(ctx, "/google.bigtable.admin.table.v1.BigtableTableService/RenameTable", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bigtableTableServiceClient) CreateColumnFamily(ctx context.Context, in *CreateColumnFamilyRequest, opts ...grpc.CallOption) (*google_bigtable_admin_table_v11.ColumnFamily, error) {
	out := new(google_bigtable_admin_table_v11.ColumnFamily)
	err := grpc.Invoke(ctx, "/google.bigtable.admin.table.v1.BigtableTableService/CreateColumnFamily", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bigtableTableServiceClient) UpdateColumnFamily(ctx context.Context, in *google_bigtable_admin_table_v11.ColumnFamily, opts ...grpc.CallOption) (*google_bigtable_admin_table_v11.ColumnFamily, error) {
	out := new(google_bigtable_admin_table_v11.ColumnFamily)
	err := grpc.Invoke(ctx, "/google.bigtable.admin.table.v1.BigtableTableService/UpdateColumnFamily", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bigtableTableServiceClient) DeleteColumnFamily(ctx context.Context, in *DeleteColumnFamilyRequest, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := grpc.Invoke(ctx, "/google.bigtable.admin.table.v1.BigtableTableService/DeleteColumnFamily", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bigtableTableServiceClient) BulkDeleteRows(ctx context.Context, in *BulkDeleteRowsRequest, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := grpc.Invoke(ctx, "/google.bigtable.admin.table.v1.BigtableTableService/BulkDeleteRows", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for BigtableTableService service

type BigtableTableServiceServer interface {
	// Creates a new table, to be served from a specified cluster.
	// The table can be created with a full set of initial column families,
	// specified in the request.
	CreateTable(context.Context, *CreateTableRequest) (*google_bigtable_admin_table_v11.Table, error)
	// Lists the names of all tables served from a specified cluster.
	ListTables(context.Context, *ListTablesRequest) (*ListTablesResponse, error)
	// Gets the schema of the specified table, including its column families.
	GetTable(context.Context, *GetTableRequest) (*google_bigtable_admin_table_v11.Table, error)
	// Permanently deletes a specified table and all of its data.
	DeleteTable(context.Context, *DeleteTableRequest) (*google_protobuf1.Empty, error)
	// Changes the name of a specified table.
	// Cannot be used to move tables between clusters, zones, or projects.
	RenameTable(context.Context, *RenameTableRequest) (*google_protobuf1.Empty, error)
	// Creates a new column family within a specified table.
	CreateColumnFamily(context.Context, *CreateColumnFamilyRequest) (*google_bigtable_admin_table_v11.ColumnFamily, error)
	// Changes the configuration of a specified column family.
	UpdateColumnFamily(context.Context, *google_bigtable_admin_table_v11.ColumnFamily) (*google_bigtable_admin_table_v11.ColumnFamily, error)
	// Permanently deletes a specified column family and all of its data.
	DeleteColumnFamily(context.Context, *DeleteColumnFamilyRequest) (*google_protobuf1.Empty, error)
	// Delete all rows in a table corresponding to a particular prefix
	BulkDeleteRows(context.Context, *BulkDeleteRowsRequest) (*google_protobuf1.Empty, error)
}

func RegisterBigtableTableServiceServer(s *grpc.Server, srv BigtableTableServiceServer) {
	s.RegisterService(&_BigtableTableService_serviceDesc, srv)
}

func _BigtableTableService_CreateTable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTableRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BigtableTableServiceServer).CreateTable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.bigtable.admin.table.v1.BigtableTableService/CreateTable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BigtableTableServiceServer).CreateTable(ctx, req.(*CreateTableRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BigtableTableService_ListTables_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTablesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BigtableTableServiceServer).ListTables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.bigtable.admin.table.v1.BigtableTableService/ListTables",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BigtableTableServiceServer).ListTables(ctx, req.(*ListTablesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BigtableTableService_GetTable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTableRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BigtableTableServiceServer).GetTable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.bigtable.admin.table.v1.BigtableTableService/GetTable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BigtableTableServiceServer).GetTable(ctx, req.(*GetTableRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BigtableTableService_DeleteTable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTableRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BigtableTableServiceServer).DeleteTable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.bigtable.admin.table.v1.BigtableTableService/DeleteTable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BigtableTableServiceServer).DeleteTable(ctx, req.(*DeleteTableRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BigtableTableService_RenameTable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RenameTableRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BigtableTableServiceServer).RenameTable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.bigtable.admin.table.v1.BigtableTableService/RenameTable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BigtableTableServiceServer).RenameTable(ctx, req.(*RenameTableRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BigtableTableService_CreateColumnFamily_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateColumnFamilyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BigtableTableServiceServer).CreateColumnFamily(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.bigtable.admin.table.v1.BigtableTableService/CreateColumnFamily",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BigtableTableServiceServer).CreateColumnFamily(ctx, req.(*CreateColumnFamilyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BigtableTableService_UpdateColumnFamily_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_bigtable_admin_table_v11.ColumnFamily)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BigtableTableServiceServer).UpdateColumnFamily(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.bigtable.admin.table.v1.BigtableTableService/UpdateColumnFamily",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BigtableTableServiceServer).UpdateColumnFamily(ctx, req.(*google_bigtable_admin_table_v11.ColumnFamily))
	}
	return interceptor(ctx, in, info, handler)
}

func _BigtableTableService_DeleteColumnFamily_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteColumnFamilyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BigtableTableServiceServer).DeleteColumnFamily(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.bigtable.admin.table.v1.BigtableTableService/DeleteColumnFamily",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BigtableTableServiceServer).DeleteColumnFamily(ctx, req.(*DeleteColumnFamilyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BigtableTableService_BulkDeleteRows_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BulkDeleteRowsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BigtableTableServiceServer).BulkDeleteRows(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.bigtable.admin.table.v1.BigtableTableService/BulkDeleteRows",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BigtableTableServiceServer).BulkDeleteRows(ctx, req.(*BulkDeleteRowsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _BigtableTableService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.bigtable.admin.table.v1.BigtableTableService",
	HandlerType: (*BigtableTableServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateTable",
			Handler:    _BigtableTableService_CreateTable_Handler,
		},
		{
			MethodName: "ListTables",
			Handler:    _BigtableTableService_ListTables_Handler,
		},
		{
			MethodName: "GetTable",
			Handler:    _BigtableTableService_GetTable_Handler,
		},
		{
			MethodName: "DeleteTable",
			Handler:    _BigtableTableService_DeleteTable_Handler,
		},
		{
			MethodName: "RenameTable",
			Handler:    _BigtableTableService_RenameTable_Handler,
		},
		{
			MethodName: "CreateColumnFamily",
			Handler:    _BigtableTableService_CreateColumnFamily_Handler,
		},
		{
			MethodName: "UpdateColumnFamily",
			Handler:    _BigtableTableService_UpdateColumnFamily_Handler,
		},
		{
			MethodName: "DeleteColumnFamily",
			Handler:    _BigtableTableService_DeleteColumnFamily_Handler,
		},
		{
			MethodName: "BulkDeleteRows",
			Handler:    _BigtableTableService_BulkDeleteRows_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor1 = []byte{
	// 378 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xb4, 0x94, 0x3f, 0x4f, 0xeb, 0x30,
	0x14, 0xc5, 0xfb, 0x96, 0xf7, 0x9e, 0x5c, 0xe9, 0x0d, 0xd6, 0x13, 0x43, 0x90, 0x18, 0x2a, 0xb1,
	0x21, 0x47, 0x2d, 0x62, 0x60, 0x4d, 0xf9, 0xb3, 0x30, 0x54, 0xa5, 0x2c, 0x30, 0x44, 0x4e, 0x72,
	0xb1, 0x0c, 0xfe, 0x13, 0x62, 0xa7, 0xa8, 0x13, 0x5f, 0x94, 0x0f, 0x83, 0x12, 0xd7, 0x90, 0x42,
	0x85, 0x9b, 0x81, 0xa5, 0xaa, 0x7d, 0xcf, 0x39, 0xbf, 0xdc, 0x7b, 0xa3, 0xa0, 0x5b, 0xa6, 0x35,
	0x13, 0x40, 0x98, 0x16, 0x54, 0x31, 0xa2, 0x2b, 0x16, 0xe7, 0x42, 0xd7, 0x45, 0x9c, 0x71, 0x66,
	0x69, 0x26, 0x20, 0xe6, 0xca, 0x42, 0xa5, 0xa8, 0x88, 0xdb, 0x63, 0x6a, 0xa0, 0x5a, 0xf2, 0x1c,
	0xd2, 0xb2, 0xd2, 0x56, 0xbf, 0xab, 0xd2, 0x8d, 0x22, 0x69, 0x8b, 0xf8, 0x60, 0x9d, 0xed, 0x45,
	0x84, 0x16, 0x92, 0x2b, 0xe2, 0xfe, 0x2f, 0xc7, 0xd1, 0xa2, 0x2f, 0xbb, 0xa0, 0x96, 0x6e, 0x07,
	0x37, 0x15, 0x47, 0x8d, 0xf2, 0x9f, 0xe8, 0x28, 0x95, 0x60, 0x0c, 0x65, 0x60, 0xd6, 0x90, 0x7d,
	0x07, 0x89, 0xdb, 0x53, 0x56, 0xdf, 0xc7, 0x20, 0x4b, 0xbb, 0x72, 0xc5, 0xc9, 0xeb, 0x1f, 0xf4,
	0x3f, 0x59, 0xc7, 0x2c, 0x9a, 0x9f, 0x6b, 0x17, 0x82, 0x1f, 0xd0, 0x70, 0x5a, 0x01, 0xb5, 0xee,
	0x16, 0x4f, 0xc8, 0xf7, 0x03, 0x22, 0x1d, 0xf1, 0x1c, 0x9e, 0x6a, 0x30, 0x36, 0x3a, 0x0c, 0x79,
	0x5a, 0xf5, 0x68, 0x80, 0x6b, 0x84, 0xae, 0xb8, 0xb1, 0xed, 0xd1, 0xe0, 0x71, 0xc8, 0xf6, 0xa1,
	0xf5, 0xa4, 0x49, 0x1f, 0x8b, 0x29, 0xb5, 0x32, 0x0d, 0xb6, 0x40, 0x7f, 0x2f, 0xc1, 0x5d, 0xe3,
	0x38, 0x94, 0xe0, 0x95, 0xbd, 0x9b, 0xbb, 0x43, 0xc3, 0x33, 0x10, 0xb0, 0xf3, 0x20, 0x3b, 0x62,
	0xcf, 0xda, 0xf3, 0x1e, 0xbf, 0x42, 0x72, 0xde, 0xac, 0xd0, 0x85, 0xcf, 0x41, 0x51, 0xb9, 0x6b,
	0x78, 0x47, 0x1c, 0x0e, 0x7f, 0x41, 0xd8, 0x6d, 0x75, 0xaa, 0x45, 0x2d, 0xd5, 0x05, 0x95, 0x5c,
	0xac, 0xf0, 0xe9, 0x6e, 0x6f, 0x42, 0xd7, 0xe3, 0x51, 0x47, 0x41, 0x6b, 0xc7, 0x34, 0x1a, 0xe0,
	0x0a, 0xe1, 0x9b, 0xb2, 0xf8, 0xfc, 0x00, 0xbd, 0x52, 0x7a, 0x33, 0x39, 0xc2, 0x6e, 0x03, 0xfd,
	0x9a, 0xfe, 0xea, 0x09, 0xcf, 0x97, 0xa2, 0x7f, 0x49, 0x2d, 0x1e, 0x9d, 0x75, 0xae, 0x9f, 0x0d,
	0x3e, 0x09, 0x61, 0x36, 0xf5, 0x41, 0x44, 0x92, 0xa0, 0x51, 0xae, 0x65, 0x20, 0x35, 0x89, 0xb6,
	0x7d, 0x01, 0xcc, 0xac, 0x09, 0x9b, 0xfd, 0xca, 0x7e, 0xb7, 0xa9, 0xc7, 0x6f, 0x01, 0x00, 0x00,
	0xff, 0xff, 0x61, 0xcc, 0xfb, 0x30, 0x7f, 0x05, 0x00, 0x00,
}
