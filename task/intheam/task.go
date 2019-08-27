package intheam

import "time"

type Task struct {
	ID          string        `json:"id"`
	UUID        string        `json:"uuid"`
	ShortID     int           `json:"short_id"`
	Status      string        `json:"status"`
	Urgency     float64       `json:"urgency"`
	Description string        `json:"description"`
	Project     string        `json:"project"`
	Entry       time.Time     `json:"entry"`
	Modified    time.Time     `json:"modified"`
	Depends     []string      `json:"depends"`
	Blocks      []interface{} `json:"blocks"`
	Tags        []string      `json:"tags"`
	Udas        struct{}      `json:"udas"`
}
