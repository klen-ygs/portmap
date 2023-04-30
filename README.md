# portmap

#### 介绍
端口映射工具，需要一个处在公网环境的服务器做流量转发


#### 安装教程

1.  安装Go环境，请参考https://go.dev
2.  在终端输入 **go install gitee.com/klenYGS/portmap**

## 使用说明

#### 一. 生成服务端配置文件

 1. 打开终端输入 **portmap srvconf -p srvconfig.yaml**

    ```yaml
    instances:
        # name是该端口映射实例的名字
      - name: 服务器
        # webPort是用户连接的端口
        webPort: 8080
        # outPort是要映射端口的主机连接服务器的端口
        outPort: 1501
        # password是主机连接服务器验证密码
        password: klenLinux
        # key是主机与服务器通信时的加密生成key，长度必须为16
        key: 1234567891234567
    ```

    服务端可以转发多个端口映射实例，只需在instances之后添加多个实例配置即可

2) 进入srvconfig.yaml修改配置文件

#### 二. 生成主机端配置文件

1. 打开终端输入 **portmap cliconf -p cliconfig.yaml**

   ```yaml
   
   # outPort是要映射的端口
   outPort: 8000
   # serverEndpoint是连接服务器的地址，其中端口对应服务器的outPort
   serverEndpoint: 127.0.0.1:1501
   # password是主机与服务器连接验证的密码
   password: klenLinux
   # key是主机与服务器通信时加密使用的key，长度必须为16
   key: 1234567891234567
   # retry是连接服务器失败时，重试时间，单位为秒
   retry: 5
   ```

   

2) 进入cliconfig.yaml修改配置文件

#### 三. 打开服务端软件 

```shell
portmap srvrun -c srvconfig.yaml
```



#### 四. 打开主机端软件 

```shell
portmap clirun -c cliconfig.yaml
```

