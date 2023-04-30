package client

import (
	"gitee.com/klenYGS/outport/asecrypto"
	"reflect"
	"testing"
)

func TestReadConfig(t *testing.T) {
	LoadConfigFromFile("config.yaml")
	test := InstanceConfig{
		OutPort:        8080,
		ServerEndpoint: "127.0.0.1:1501",
		Password:       "klenLinux",
		Key:            "1234567891234567",
		RetryTime:      5,
	}
	if !reflect.DeepEqual(test, config) {
		t.Error("不一致")
	}
}

func TestMany(t *testing.T) {
	test := "123456789123456789"
	key := "1234567891234567"

	data1, err := asecrypto.Encrypt([]byte(key), []byte(test))
	if err != nil {
		t.Fatal(err)
	}

	data2, err := asecrypto.Encrypt([]byte(key), []byte(test))
	if err != nil {
		t.Fatal(err)
	}

	if reflect.DeepEqual(data1, data2) {
		t.Errorf("一致")
	}

}
