# portmap

#### 介绍
端口映射工具

#### 软件架构
软件架构说明


#### 安装教程

1.  安装Go环境，请参考https://go.dev
2.  在终端输入 go install gitee.com/klenYGS/portmap

#### 使用说明

1.  生成服务端配置文件
    1）打开终端输入 portmap srvconf -p srvconfig.yaml
    2) 进入srvconfig.yaml修改配置文件
2.  生成主机端配置文件
    1）打开终端输入 portmap cliconf -p cliconfig.yaml
    2) 进入cliconfig.yaml修改配置文件
3.  打开服务端软件 portmap srvrun -c srvconfig.yaml
4.  打开主机端软件 portmap clirun -c cliconfig.yaml
