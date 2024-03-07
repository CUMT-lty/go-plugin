// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// 插件进程 & 宿主进程（插件调用方）共同需要的接口及其实现类

package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Greeter is the interface that we're exposing as a plugin.
// 插件的业务接口，其实就是规定了插件的功能
// 这个业务接口，插件进程和宿主进程都需要实现一遍
type Greeter interface {
	Greet() string
}

// GreeterMain Here is an implementation that talks over RPC
// 宿主进程的插件业务接口实现类
// 宿主进程中插件业务实现类的每一个方法，都是要去调插件进程中插件业务实现类的对应的方法。
type GreeterMain struct{ client *rpc.Client }

func (g *GreeterMain) Greet() string {
	var resp string
	err := g.client.Call("Plugin.Greet", new(interface{}), &resp) // rpc 调用插件进程中的对应方法
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// GreeterRPCServer Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
// 被插件进程使用，承载 rpc/grpc 通信，也就是说，是这个结构体真正调用插件进程中的业务方法
type GreeterRPCServer struct {
	// This is the real implementation
	Impl Greeter // 这里应该会被赋值为插件进程的业务实现类对象
}

func (s *GreeterRPCServer) Greet(args interface{}, resp *string) error {
	*resp = s.Impl.Greet() // 调用插件进程的业务实现类方法
	return nil
}

// GreeterPlugin This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a GreeterRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return GreeterRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
// 插件 Plugin 接口实现类，被插件进程 和 宿主进程共用
type GreeterPlugin struct {
	// Impl Injection
	Impl Greeter
}

// Server 方法由插件进程调用，返回的是承载 rpc/grpc 通信结构体对象
func (p *GreeterPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &GreeterRPCServer{Impl: p.Impl}, nil
}

// Client 方法由宿主进程调用，返回的是宿主进程的插件业务实现类对象
func (GreeterPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &GreeterMain{client: c}, nil
}
