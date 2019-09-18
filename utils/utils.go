package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/snluu/uuid"
	"lyanna/models"
	"math/rand"
	"time"
)

func Md5(source string) string {
	md5h := md5.New()
	md5h.Write([]byte(source))
	return hex.EncodeToString(md5h.Sum(nil))
}

func UUID() string {
	return uuid.Rand().Hex()
}

func RandomGetArray(origin []*models.Post, limit int) []*models.Post{
	rand.Shuffle(len(origin), func(i, j int) {
		origin[i],origin[j] = origin[j],origin[i]
	})
	if limit > len(origin) - 1 {
		return origin
	}
	return origin[:limit]
}

func GetCurrentTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(loc)
}