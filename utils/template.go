package utils

import "time"

// 格式化时间
func DateFormat(date time.Time, layout string) string {
	return date.Format(layout)
}

