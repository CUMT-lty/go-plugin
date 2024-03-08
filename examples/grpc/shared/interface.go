// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package shared contains shared data between the host and plugins.
package shared

import (
	"context"
	"net/rpc"

	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-plugin/examples/grpc/proto"
)

// Handshake is a common handshake that is shared by plugin and host.
// 握手配置
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// KV is the interface that we're exposing as a plugin.
// 插件业务接口
type KV interface {
	Put(key string, value []byte) error
	Get(key string) ([]byte, error)
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
// 这个先不看
type KVPlugin struct {
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl KV
}

func (p *KVPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

func (*KVPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}

// KVGRPCPlugin This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
// 最上层使用 grpc 协议的插件接口实现类，实现类两个接口 Plugin 接口和 plugin.GRPCPlugin 接口
// 插件进程和宿主进程都会用到
type KVGRPCPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl KV
}

// GRPCServer grpc 插件进程会用到的方法
func (p *KVGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	// 注册 grpc server
	// 注册的是包装业务实现的接口的 grpc server 对象
	proto.RegisterKVServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

// GRPCClient grpc 宿主进程会用到的方法
// 返回的是一个 grpc
func (p *KVGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewKVClient(c)}, nil
}
