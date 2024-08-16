package main

import (
	"github.com/PirateDreamer/gotab/cli"
	"github.com/PirateDreamer/gotab/logic"
)

func main() {
	// 创建默认配置
	option := &logic.Option{}
	cli.Execute(option)
	logic.HandlerOption(option)
	logic.Run(option)
}
