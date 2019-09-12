package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/snluu/uuid"
)

func Md5(source string) string {
	md5h := md5.New()
	md5h.Write([]byte(source))
	return hex.EncodeToString(md5h.Sum(nil))
}

func UUID() string {
	return uuid.Rand().Hex()
}
