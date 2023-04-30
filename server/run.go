package server

import (
	"github.com/spf13/cobra"
	"log"
)

var (
	path string
)

var RunCmd = &cobra.Command{
	Use:     "srvrun",
	GroupID: "exec",
	Short:   "运行端口映射服务器",
	Example: "portmap srvrun -c config.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		LoadConfigFromYaml(path)
		var instances []*Instance

		for _, instanceConfig := range config.InstanceConfigs {
			instance := NewInstance(instanceConfig)
			instance.Run()
			log.Printf("[info] \n%s\n", instance.String())
			instances = append(instances, instance)
		}
		select {}
	},
}

func init() {
	RunCmd.Flags().StringVarP(&path, "conf", "c", "", "指定配置文件")
}
