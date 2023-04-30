package client

import "github.com/spf13/cobra"

var (
	path string
)

var CreateConfCmd = &cobra.Command{
	Use:     "cliconf",
	GroupID: "exec",
	Short:   "生成主机端的默认yaml配置文件",
	Run: func(cmd *cobra.Command, args []string) {
		CreateConfigFile(path)
	},

	Example: "portmap cliconf -p config.yaml",
}

func init() {
	CreateConfCmd.Flags().StringVarP(&path, "path", "p", "", "生成的配置文件路径")
}
