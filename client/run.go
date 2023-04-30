package client

import "github.com/spf13/cobra"

var (
	confPath string
)

var RunCmd = &cobra.Command{
	Use:     "clirun",
	GroupID: "exec",
	Short:   "执行主机端端口映射程序",
	Run: func(cmd *cobra.Command, args []string) {
		LoadConfigFromFile(confPath)
		instance := NewInstance(config)
		instance.Run()
	},

	Example: "portmap clirun -c config.yaml",
}

func init() {
	RunCmd.Flags().StringVarP(&confPath, "conf", "c", "", "指定执行的配置文件")
}
