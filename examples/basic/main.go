// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// 宿主进程（也就是插件调用放，仅定义插件调用放需要的实现）

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-plugin/examples/basic/shared"
)

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

//
//// pluginMap is the map of plugins we can dispense.
//// 宿主机进程的插件集
//var pluginMap = map[string]plugin.Plugin{
//	"greeter": &shared.GreeterPlugin{},
//}

// 插件调用者进程（宿主进程）
func main() {
	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	// pluginMap is the map of plugins we can dispense.
	// 宿主机进程的插件集
	var pluginMap = map[string]plugin.Plugin{
		// 插件名字
		// 穿一个插件 Plugin 接口实现类对象
		// 注意，在插件调用方没有给 GreeterPlugin 的 Impl 字段赋值（因为这个是要通过 rpc 从插件进程拿的）
		"greeter": &shared.GreeterPlugin{},
	}

	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{ // 每一个插件进程对应一个 Client 客户端实例，用来管理插件进程的生命周期
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command("./plugin/greeter"), // 命令行启动插件进程
		Logger:          logger,
	})
	defer client.Kill()

	// Connect via RPC
	// 获得插件服务的客户端
	rpcClient, err := client.Client() // 这个 rpcClient 是对标准的 rpc 或 grpc 客户端的包装
	if err != nil {
		log.Fatal(err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("greeter") // 返回的是宿主进程定义的插件业务接口实现类对象
	if err != nil {
		log.Fatal(err)
	}

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	greeter := raw.(shared.Greeter) // 类型断言，转换为实际接口类型，这个 greeter 是 GreeterMain
	fmt.Println(greeter.Greet())    // 内部其实发起了 rpc 调用
}

// TODO：思考一个问题
// 如果两方协作开发，A只提供插件功能，B只想调用插件功能
// 那么哪些代码需要 A 写，哪些代码需要 B 写
