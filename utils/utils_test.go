package utils

import (
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {
	username := "admin"
	password := "123456"
	md5_res := Md5(username + password)
	fmt.Println(md5_res)
	//t.Log(md5_res)
}