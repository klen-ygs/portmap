package server

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"gitee.com/klenYGS/portmap/asecrypto"
)

type Instance struct {
	config          InstanceConfig
	computer        net.Conn
	listenWeb       net.Listener
	listenTarget    net.Listener
	connChannel     chan net.Conn
	closeChannel    chan struct{}
	closeErr        error
	closed          atomic.Bool
	newComputerPass string
	newConnPass     string
}

func NewInstance(config InstanceConfig) *Instance {
	if len(config.Key) != 16 {
		panic(fmt.Errorf("实例:%s key长度预期为16，实际为%d", config.Name, len(config.Key)).Error())
	}
	newComputer := fmt.Sprintf("newConn for %s", config.Password)

	newConnPass := config.Password

	return &Instance{
		config:          config,
		computer:        nil,
		connChannel:     make(chan net.Conn, 1),
		closeChannel:    make(chan struct{}, 1),
		newComputerPass: newComputer,
		newConnPass:     newConnPass,
	}
}

func (i *Instance) Run() {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go i.listenFromWeb(wg)
	go i.listenForTarget(wg)
	go i.waitForClose(wg)

}

func (i *Instance) listenFromWeb(wg *sync.WaitGroup) {
	listen, err := net.Listen("tcp", ":"+strconv.Itoa(i.config.WebPort))
	if err != nil {
		i.closeErr = errors.Join(i.closeErr, fmt.Errorf("映射:%s 打开tcp监听失败:%w", i.config.Name, err))
		i.notifyClose()
		return
	}
	defer listen.Close()
	i.listenWeb = listen
	wg.Done()

	for {
		userConn, err := listen.Accept()
		if err != nil {
			if !i.closed.Load() {
				i.closeErr = errors.Join(fmt.Errorf("映射:%s 获取用户连接失败:%w", i.config.Name, err))
			}
			i.notifyClose()
			return
		}
		log.Printf("[info] 新连接接入 地址:%s", userConn.RemoteAddr())
		go i.newConn(userConn)
	}
}

func (i *Instance) sendNewConn() {
	if i.computer == nil {
		log.Printf("[warn] %s的主机未连接", i.config.Name)
		return
	}

	buf := make([]byte, 100)
	send := fmt.Sprintf("newConn for %s", i.config.Password)
	data, err := asecrypto.Encrypt([]byte(i.config.Key), []byte(send))
	if err != nil {
		if !i.closed.Load() {
			i.closeErr = errors.Join(i.closeErr, fmt.Errorf("newConn加密错误：%w", err))
		}
		i.notifyClose()
		return
	}

	copy(buf[binary.MaxVarintLen16:], data)
	binary.PutVarint(buf[:binary.MaxVarintLen16], int64(len(data)))
	_, err = i.computer.Write(buf)
	if err != nil {
		if !i.closed.Load() {
			i.closeErr = errors.Join(i.closeErr, fmt.Errorf("[Error] 实例: %s 向主机发送新连接命令失败", i.config.Name))
		}
		i.notifyClose()
	}

}

func (i *Instance) tcpToTcp(endChan chan struct{}, dst, src net.Conn) {
	_, err := io.Copy(dst, src)
	if err != nil {
		log.Printf("[Info] 实例:%s | %s -> %s 连接被关闭", i.config.Name, src.RemoteAddr(), dst.RemoteAddr())
	}
	endChan <- struct{}{}
}

func (i *Instance) newConn(userConn net.Conn) {
	i.sendNewConn()
	newComputerConn := <-i.connChannel

	endChan := make(chan struct{}, 2)
	go i.tcpToTcp(endChan, newComputerConn, userConn)
	go i.tcpToTcp(endChan, userConn, newComputerConn)

	<-endChan
	userConn.Close()
	newComputerConn.Close()
}

func (i *Instance) notifyClose() {
	i.closed.Store(true)
	select {
	case i.closeChannel <- struct{}{}:
	default:
	}
}

func (i *Instance) listenForTarget(wg *sync.WaitGroup) {
	listen, err := net.Listen("tcp", ":"+strconv.Itoa(i.config.OutPort))
	if err != nil {
		if !i.closed.Load() {
			i.closeErr = errors.Join(i.closeErr, fmt.Errorf("打开监听目标主机失败"))
		}
		i.notifyClose()
		return
	}
	defer listen.Close()

	i.listenTarget = listen
	wg.Done()

	for {
		tarConn, err := listen.Accept()
		if err != nil {
			if !i.closed.Load() {
				i.closeErr = errors.Join(i.closeErr, fmt.Errorf("获取目标主机连接失败: %w", err))
			}
			i.notifyClose()
			return
		}
		if isOur, isComputer, err := i.isOurComputerAndConnType(tarConn); err == nil {
			if !isOur {
				log.Printf("[warn] 非目标主机连入 地址%s", tarConn.RemoteAddr())
				tarConn.Close()
				continue
			}
			if isComputer {
				log.Printf("[info] 目标主机连入 地址:%s", tarConn.RemoteAddr())
				i.computer = tarConn
				continue
			} else {
				log.Printf("[info] 主机新连接 地址:%s", tarConn.RemoteAddr())
				i.connChannel <- tarConn
			}

		} else {
			log.Printf("[Error] 验证连接身份失败 %s", err.Error())
			tarConn.Close()
		}
	}

}

func (i *Instance) read100Byte(ctx context.Context, conn net.Conn) ([]byte, error) {
	buf := make([]byte, 100)
	var count int

	for {
		if count == 100 {
			break
		}
		select {
		default:
		case <-ctx.Done():
			return nil, errors.New("主机未在规定时间内认证")
		}
		conn.SetReadDeadline(time.Now().Add(time.Second * 5))
		n, err := conn.Read(buf[count:])
		if err != nil {
			return nil, err
		}

		count += n
	}
	conn.SetReadDeadline(time.Time{})
	return buf, nil
}

func (i *Instance) isOurComputerAndConnType(conn net.Conn) (bool, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	data, err := i.read100Byte(ctx, conn)
	if err != nil {
		return false, false, err
	}

	dataLen, n := binary.Varint(data[:binary.MaxVarintLen16])
	if n <= 0 {
		return false, false, fmt.Errorf("[warn] 数据加密头错误")
	}
	data = data[binary.MaxVarintLen16 : dataLen+binary.MaxVarintLen16]
	data, err = asecrypto.Decrypt([]byte(i.config.Key), data)
	if err != nil {
		return false, false, fmt.Errorf("密文解密失败:%w", err)
	}
	dataStr := string(data)
	if dataStr == i.newConnPass {
		return true, false, nil
	} else if dataStr == i.newComputerPass {
		return true, true, nil
	}

	return false, false, nil
}

func (i *Instance) waitForClose(wg *sync.WaitGroup) {
	wg.Wait()
	<-i.closeChannel

	if i.listenTarget != nil {
		i.listenTarget.Close()
	}
	if i.listenWeb != nil {
		i.listenWeb.Close()
	}

}

func (i *Instance) String() string {
	return fmt.Sprintf("映射实例:%s\nweb端口:%d\n主机连接端口:%d\n", i.config.Name, i.config.WebPort, i.config.OutPort)
}
