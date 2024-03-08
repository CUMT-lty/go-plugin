// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shared

import (
	"context"

	"github.com/hashicorp/go-plugin/examples/grpc/proto"
)

// GRPCClient is an implementation of KV that talks over RPC.
// 属于宿主进程，包装的 grpc 客户端，内部方法是发起对应的 rpc 调用
// 充当 stub 的角色
type GRPCClient struct{ client proto.KVClient }

func (m *GRPCClient) Put(key string, value []byte) error {
	_, err := m.client.Put(context.Background(), &proto.PutRequest{
		Key:   key,
		Value: value,
	})
	return err
}

func (m *GRPCClient) Get(key string) ([]byte, error) {
	resp, err := m.client.Get(context.Background(), &proto.GetRequest{
		Key: key,
	})
	if err != nil {
		return nil, err
	}

	return resp.Value, nil
}

// GRPCServer Here is the gRPC server that GRPCClient talks to.
// 属于插件进程，包装插件业务接口实现，真正去调用业务接口对象的方法
type GRPCServer struct {
	// This is the real implementation
	Impl KV
}

func (m *GRPCServer) Put(ctx context.Context, req *proto.PutRequest) (*proto.Empty, error) {
	return &proto.Empty{}, m.Impl.Put(req.Key, req.Value)
}

func (m *GRPCServer) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	v, err := m.Impl.Get(req.Key)
	return &proto.GetResponse{Value: v}, err
}
