package domain

import "time"

type CommentCreateRequest struct {
	PostId      int64     `json:"postId"`
	Content     string    `json:"content"`
	Password    string    `json:"password"`
	MessageType string    `json:"messageType"`
	Timestamp   time.Time `json:"timestamp"`
}

type CommentLikeRequest struct {
	PostId      int64  `json:"postId"`
	CommentId   string `json:"commentId"`
	UserId      string `json:"userId"`
	MessageType string `json:"messageType"`
}

type CommentUpdateRequest struct {
	PostId      int64     `json:"postId"`
	CommentId   string    `json:"commentId"`
	Content     string    `json:"content"`
	Password    string    `json:"password"`
	MessageType string    `json:"messageType"`
	Timestamp   time.Time `json:"timestamp"`
}
