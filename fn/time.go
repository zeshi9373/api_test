package fn

import "time"

// 时间戳(秒)
func Time() int64 {
	return time.Now().Unix()
}

// 时间戳(毫秒)
func TimeMillis() int64 {
	return time.Now().UnixMilli()
}

// 时间戳(微秒)
func TimeMicros() int64 {
	return time.Now().UnixMicro()
}

// 时间戳(纳秒)
func TimeNanos() int64 {
	return time.Now().UnixNano()
}

// 日期
func Date() string {
	return time.Now().Format("2006-01-02")
}

// 时间
func DateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
