package utils

import "time"

func Timestamp2ymd(tm int64) string {
	// 示例 Unix 时间戳（秒）
	unixTimestamp := tm

	// 将 Unix 时间戳转换为 time.Time 类型
	timeFromUnix := time.Unix(unixTimestamp, 0)

	// 定义格式化布局字符串
	layout := "2006-01-02 15:04:05"

	// 格式化时间为字符串
	formattedTime := timeFromUnix.Format(layout)
	return formattedTime
}

func GetTimestamp() int64 {
	// 获取当前时间
	currentTime := time.Now()

	// 将当前时间转换为 Unix 时间戳（秒）
	unixTimestamp := currentTime.Unix()

	return unixTimestamp
}
