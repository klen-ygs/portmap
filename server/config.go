package server

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	InstanceConfigs []InstanceConfig `mapstructure:"instances"`
}

// InstanceConfig 是一个端口映射在服务器端的配置
type InstanceConfig struct {

	// WebPort 是用户连接时的端口
	WebPort int `mapstruture:"webPort"`

	// LinuxPort 是用户映射主机连接服务器的端口
	OutPort int `mapstructure:"outPort"`

	// Password 是用户主机与服务器连接时的口令，以防外人连接
	Password string `mapstructure:"password"`

	// Name 是该映射的名字，用于在映射组中区分其他映射
	Name string `mapstructure:"name"`

	// Key 是服务器与目标主机通信时的加密的Key，长度应为16字节
	Key string `mapstructure:"key"`
}

// config 是全局配置
var config Config

// LoadConfigFromYaml 是从指定路径中读取配置并加载到全局配置config中
// 若路径为空，则注入默认配置
func LoadConfigFromYaml(path string) {
	if path == "" {
		defaultConfig()
		return
	}

	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("读取配置文件失败: %w", err))
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("加载配置文件失败: %w", err))
	}
}

// defaultConfig 是用户未指定配置文件时，调用此函数默认配置
func defaultConfig() {
	log.Printf("[info] 使用默认配置")
	config.InstanceConfigs = append(config.InstanceConfigs, InstanceConfig{
		WebPort:  1500,
		OutPort:  1501,
		Password: "klenLinux",
		Key:      "1234567891234567",
	})
}

var configExample = `instances:
    # name是该端口映射实例的名字
  - name: example
    # webPort是用户连接的端口
    webPort: 1500
    # outPort是要映射端口的主机连接服务器的端口
    outPort: 1501
    # password是主机连接服务器验证密码
    password: klenLinux
    # key是主机与服务器通信时的加密生成key，长度必须为16
    key: 1234567891234567
`

func CreateExampleFile(path string) {
	file, err := os.Create(path)
	if err != nil {
		panic(fmt.Errorf("创建文件失败:%w", err).Error())
	}
	defer file.Close()

	_, err = file.Write([]byte(configExample))
	if err != nil {
		panic(fmt.Errorf("写入文件失败:%w", err))
	}
}
