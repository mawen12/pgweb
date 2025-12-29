package history

import (
	"time"
)

// 历史记录，有查询字符串，时间戳
type Record struct {
	Query     string `json:"query"`
	Timestamp string `json:"timestamp"`
}

func New() []Record {
	return make([]Record, 0)
}

func NewRecord(query string) Record {
	return Record{
		Query:     query,
		Timestamp: time.Now().String(),
	}
}
