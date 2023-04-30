package client

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"gitee.com/klenYGS/portmap/asecrypto"
)

type Instance struct {
	config          InstanceConfig
	CmdConn         net.Conn
	newComputerPass string
	newConnPass     string
	newConnCmd      string
}

func NewInstance(config InstanceConfig) *Instance {
	if len(config.Key) != 16 {
		panic(fmt.Errorf("预期key值长度为16"))
	}

	newComputer, err := asecrypto.Encrypt([]byte(config.Key), []byte(fmt.Sprintf("newConn for %s", config.Password)))
	if err != nil {
		panic(fmt.Errorf("生成通信密文失败: %w", err))
	}

	newConn, err := asecrypto.Encrypt([]byte(config.Key), []byte(config.Password))
	if err != nil {
		panic(fmt.Errorf("生成通信密文失败: %w", err))
	}

	newConnCmd := fmt.Sprintf("newConn for %s", config.Password)

	return &Instance{
		config:          config,
		newConnPass:     string(newConn),
		newComputerPass: string(newComputer),
		newConnCmd:      string(newConnCmd),
	}
}

// ConnectToServer 向服务器发起连接，并发送新主机密文
// 如果未能连接，将会不断重试，重试时间由配置文件指明
func (i *Instance) ConnectToServer() {
	log.Printf("[info] 连接服务器。。。")
	var errReport bool
	var count int
	for {
		dial, err := net.Dial("tcp", i.config.ServerEndpoint)
		if err != nil {
			if !errReport {
				log.Printf("连接服务器失败: %s\n", err.Error())
				log.Printf("重试中。。。")
			}
			log.Printf("[info] 重试次数: %d", count)
			count++
			errReport = true
			time.Sleep(time.Second * time.Duration(i.config.RetryTime))
			continue
		}

		i.CmdConn = dial
		err = i.sendNewComputer()
		if err != nil {
			dial.Close()
			continue
		}
		break
	}
	log.Printf("[info] 连接成功")
}

func (i *Instance) Run() {
	i.ConnectToServer()
	i.waitCmd()
}

// read100Byte 时从i.connCmd 中读取100个字节的数据
func (i *Instance) read100Byte() ([]byte, error) {
	var count int
	buf := make([]byte, 100)

	for {
		if count == 100 {
			break
		}
		n, err := i.CmdConn.Read(buf[count:])
		if err != nil {
			return nil, err
		}
		count += n
	}
	return buf, nil
}

// sendNewComputer 是向服务器发送主机新连接密文
func (i *Instance) sendNewComputer() error {
	buf := make([]byte, 100)
	binary.PutVarint(buf[:binary.MaxVarintLen16], int64(len(i.newComputerPass)))
	copy(buf[binary.MaxVarintLen16:], i.newComputerPass)
	_, err := i.CmdConn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

// waitCmd 是与服务器建立连接后，等待服务器下指令
// 如果发生异常情况，与服务器断开连接，将调用ConnectToServer重新连接服务器
func (i *Instance) waitCmd() {
	for {
		data, err := i.read100Byte()
		if err != nil {
			log.Printf("[warn] 与服务器连接异常: %s", err.Error())
			log.Printf("[info] 5s后发起新连接")
			time.Sleep(time.Second * 5)
			i.ConnectToServer()
			continue
		}

		dataLen, n := binary.Varint(data[:binary.MaxVarintLen16])
		if n <= 0 {
			log.Printf("[warn] 密文头部错误")
			continue
		}
		data = data[binary.MaxVarintLen16 : dataLen+binary.MaxVarintLen16]
		data, err = asecrypto.Decrypt([]byte(i.config.Key), data)
		if err != nil {
			log.Printf("[warn] 密文解密错误")
			continue
		}
		cmd := string(data)
		if cmd == i.newConnCmd {
			log.Printf("[info] 新连接命令\n")
			go i.newConnToServer()
		} else {
			log.Printf("[warn] 收到异常命令")
		}
	}
}

// newConnToServer 是发起新连接到服务器，同时发起新连接到本机映射的端口，
// 之后转发两个连接发送的信息
func (i *Instance) newConnToServer() {
	dial, err := net.Dial("tcp", i.config.ServerEndpoint)
	if err != nil {
		return
	}
	defer dial.Close()

	conn, err := net.Dial("tcp", ":"+strconv.Itoa(i.config.OutPort))
	if err != nil {
		return
	}
	defer conn.Close()
	err = i.sendNewConn(dial)
	if err != nil {
		dial.Close()
		conn.Close()
		log.Printf("[Error] 发起新连接验证失败")
		return
	}

	endChan := make(chan struct{}, 2)

	go i.tcpToTcp(endChan, conn, dial)
	go i.tcpToTcp(endChan, dial, conn)
	<-endChan
}

// tcpToTcp 是一个tcp连接向另一个tcp连接发送消息，如果连接出错，则向endChan发送信息
func (i *Instance) tcpToTcp(endChan chan struct{}, dst, src net.Conn) {
	_, err := io.Copy(dst, src)
	if err != nil {
		log.Printf("[Info] %s -> %s 连接被关闭", src.RemoteAddr(), dst.RemoteAddr())
	}
	endChan <- struct{}{}
}

func (i *Instance) sendNewConn(conn net.Conn) error {
	buf := make([]byte, 100)
	binary.PutVarint(buf[:binary.MaxVarintLen16], int64(len(i.newConnPass)))
	copy(buf[binary.MaxVarintLen16:], i.newConnPass)
	_, err := conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}
