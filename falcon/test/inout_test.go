package test

import (
	"falconService/falcon"
	"fmt"
	"strings"
	"testing"
)

func TestInput (t *testing.T) {
	var AKSK string
	var myFalcon falcon.FALCON
	fmt.Println("请输入ak++sk：")
	fmt.Scanf("%s", &AKSK)

	tmpAKSK := strings.Split(AKSK, "++")
	if len(tmpAKSK) == 2 {
		myFalcon.AccessKeyId = tmpAKSK[0]
		myFalcon.SecretAccessKeyId = tmpAKSK[1]

		status := myFalcon.GetApiKey()
		if !status {
			fmt.Println("初始化 myFalcon 失败")
			return
		}

	} else {
		fmt.Println("输入格式 ak++sk")
	}

}


func TestSlice (t *testing.T) {
	commands := []string{"asdas", "aaa"}

	for _, command := range commands[1:] {
		fmt.Println(command)
	}
}
