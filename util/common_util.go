package util

import "encoding/json"

// 切片转字符串
func SliceToString(data interface{}) string {
	b, err := json.Marshal(data)
	CheckError(err, "sliceToString")
	return string(b)
}
