package server

import (
	"github.com/spf13/cobra"
)

var (
	confPath string
)

var CreateConfCmd = &cobra.Command{
	Use:     "srvconf",
	Short:   "生成服务器默认yaml配置文件",
	GroupID: "exec",
	Example: "portmap srvconf -p config.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		CreateExampleFile(confPath)
	},
}

func init() {
	CreateConfCmd.Flags().StringVarP(&confPath, "path", "p", "", "文件生成路径")
}
