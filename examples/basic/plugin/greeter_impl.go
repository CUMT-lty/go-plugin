// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// 插件进程，仅定义插件进程需要的实现

package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-plugin/examples/basic/shared"
)

// GreeterHello Here is a real implementation of Greeter
// 插件业务接口实现类，插件进程中真正提供业务功能的类
type GreeterHello struct {
	logger hclog.Logger
}

func (g *GreeterHello) Greet() string {
	g.logger.Debug("message from GreeterHello.Greet")
	return "Hello!"
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
// 握手配置
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// 插件进程
func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	// 插件业务实现类的实例对象
	greeter := &GreeterHello{
		logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	// 插件方的插件集合
	var pluginMap = map[string]plugin.Plugin{
		// 插件名称
		// Plugin 插件接口实现类，注意 Impl 的赋值：
		// 现在是在插件进程中，Impl 赋值为插件方的插件业务实现类对象
		"greeter": &shared.GreeterPlugin{Impl: greeter},
	}

	logger.Debug("message from plugin", "foo", "bar")

	// 调用plugin.Serve() 启动监听，并提供服务
	// 插件进程调用 plugin.Serve 方法后，主线程会阻塞
	// 直到客户端调用 Dispense 方法请求插件实例时，服务器端才会实例化插件
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
