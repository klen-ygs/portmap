package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestErrJoin(t *testing.T) {
	var tar error = errors.New("55")
	err := errors.Join(nil, tar)
	if !reflect.DeepEqual(err.Error(), tar.Error()) {
		t.Fatalf("错误值不相同")
	}
}

func TestChannelClose(t *testing.T) {
	c := make(chan struct{}, 1)
	select {
	default:
		t.Error("不能接收close消息")
	case c <- struct{}{}:
	}

	<-c
}

func TestReadConfig(t *testing.T) {
	LoadConfigFromYaml("config.yaml")
	test := Config{
		InstanceConfigs: []InstanceConfig{
			{
				WebPort:  1500,
				OutPort:  1501,
				Password: "klenLinux",
				Name:     "测试",
				Key:      "1234567891234567",
			},
		},
	}
	fmt.Printf("%#v", config.InstanceConfigs[0].Key)
	if !reflect.DeepEqual(test, config) {
		t.Error("不一致")
	}
}

func TestBinary(t *testing.T) {
	buf := make([]byte, 100)

	testLen := 50
	binary.PutVarint(buf, int64(testLen))
	fmt.Printf("%v", buf[:3])

	varint, i := binary.Varint(buf[:3])
	fmt.Printf("tar: %d n: %d", varint, i)
}
