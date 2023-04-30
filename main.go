package main

import (
	"gitee.com/klenYGS/portmap/client"
	"gitee.com/klenYGS/portmap/server"

	"github.com/spf13/cobra"
	"log"
)

var root = cobra.Command{}

func main() {
	defer RunErr()

	err := root.Execute()
	if err != nil {
		panic(err)
	}
}

func init() {
	root.AddGroup(&cobra.Group{
		ID:    "exec",
		Title: "运行命令",
	})

	root.AddCommand(server.CreateConfCmd)
	root.AddCommand(server.RunCmd)
	root.AddCommand(client.RunCmd)
	root.AddCommand(client.CreateConfCmd)
}

func RunErr() {
	p := recover()
	if p != nil {
		log.Printf("执行发生错误: %v", p)
	}
}
