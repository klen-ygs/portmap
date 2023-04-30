package client

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type InstanceConfig struct {
	OutPort        int    `mapstructure:"outPort"`
	ServerEndpoint string `mapstructure:"serverEndpoint"`
	Key            string `mapstructure:"key"`
	Password       string `mapstructure:"password"`
	RetryTime      int    `mapstructure:"retry"`
}

var config InstanceConfig

func LoadConfigFromFile(path string) {
	if path == "" {
		panic("请配置文件路径")
	}

	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("读取配置文件失败:%w", err))
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("获取配置信息失败:%w", err))
	}

}

var exampleFile = `
# outPort是要映射的端口
outPort: 8080
# serverEndpoint是连接服务器的地址，其中端口对应服务器的outPort
serverEndpoint: 127.0.0.1:1501
# password是主机与服务器连接验证的密码
password: klenLinux
# key是主机与服务器通信时加密使用的key，长度必须为16
key: 1234567891234567
# retry是连接服务器失败时，重试时间，单位为秒
retry: 5`

func CreateConfigFile(path string) {
	file, err := os.Create(path)
	if err != nil {
		panic(fmt.Errorf("创建文件失败:%w", err).Error())
	}
	defer file.Close()

	_, err = file.Write([]byte(exampleFile))
	if err != nil {
		panic(fmt.Errorf("写入文件失败:%w", err))
	}
}
