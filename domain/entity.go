package domain

import (
	"time"
)

type Comment struct {
	ID        int64     `json:"id"`
	PostId    int64     `json:"postId"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Content   string    `json:"content"`
	Password  string    `json:"password"`
}

type QueueName struct {
	name string
}

