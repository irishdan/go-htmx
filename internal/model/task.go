package model

import "time"

type Task struct {
	Id          int64
	Title       string
	Description string
	CompletedAt *time.Time
}
