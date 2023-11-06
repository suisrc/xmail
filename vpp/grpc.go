package vpp

import "google.golang.org/grpc"

// ===================================================================
// GrpcSrvOption 修正grpc内容
type GrpcSrvOption func() []grpc.ServerOption

// NewGinEngine engine
func NewGrpcSrv(opt GrpcSrvOption) *grpc.Server {
	if opt == nil {
		return grpc.NewServer()
	}
	rpc := grpc.NewServer(opt()...)
	return rpc
}
