package util

import "time"

// const
const (
	Millisecond = 1
	Second      = Millisecond * 1000
	Minute      = Second * 60
	Hour        = Minute * 60
	Day         = Hour * 24
	Week        = Day * 7
)

// NowTs 毫秒级时间戳
func NowTs() int64 {
	return time.Now().UTC().UnixNano() / 1e6
}

// FindInt64 查找切片中是否存在某个数字并返回所在的索引值，若不存在返回-1
func FindInt64(aim []int64, value int64) int {
	for i, v := range aim {
		if v == value {
			return i
		}
	}
	return -1
}