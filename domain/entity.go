package domain

import (
	"time"
)

type Comment struct {
	ID        string    `json:"id"`
	PostId    int64     `json:"postId"`
	Content   string    `json:"content"`
	Password  string    `json:"password"`
	Timestamp time.Time `json:"timestamp"`
}

type QueueName struct {
	name string
}

func NewComment(request CommentCreateRequest) *Comment {
	comment := new(Comment)
	comment.ID = request.CommentId
	comment.PostId = request.PostId
	comment.Password = request.Password
	comment.Timestamp = request.Timestamp
	comment.Content = request.Content
	return comment
}
