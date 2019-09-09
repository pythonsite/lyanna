package utils

import "testing"

func TestMd5(t *testing.T) {
	username := "admin"
	password := "123456"
	md5_res := Md5(username + password)
	t.Log(md5_res)
}