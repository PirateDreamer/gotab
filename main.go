package main

import (
	"gotab/cli"
	"gotab/internal"
)

func main() {
	// 创建默认配置
	option := &internal.Option{}
	cli.Execute(option)
	internal.HandlerOption(option)
	internal.Run(option)
}
