package asecrypto

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDecrypt(t *testing.T) {
	test := []string{
		"newConn",
		"klen666",
		"YGS123456789",
		"1234567891234567",
	}

	key := "1234567891234567"

	for _, testcase := range test {
		data, err := Encrypt([]byte(key), []byte(testcase))
		if err != nil {
			t.Fatal(err)
		}

		decrypt, err := Decrypt([]byte(key), []byte(data))
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(len(decrypt))
		if reflect.DeepEqual(testcase, decrypt) {
			t.Errorf("不一致\n want:%s\nget:%s\n", testcase, decrypt)
		}
	}
}
