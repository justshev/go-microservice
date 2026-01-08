package task

import "time"

type Task struct {
	ID        int64       `json:"id"`
	Name string   `json:"name"`
	Completed bool `json:"completed"`
	CreatedAt time.Time   `json:"created_at"`

}

