package JsonMessage

import (
	"strconv"
)

type Email struct {
	Address string `json:"address"` // UUID
	Title   string `json:"title"`   // 제목
	Body    string `json:"body"`    // 본문
}

func IntToString(input_num int64) string {
	return strconv.FormatInt(input_num, 10)
}
