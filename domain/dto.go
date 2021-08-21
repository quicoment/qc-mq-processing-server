package domain

import (
	"time"
)

type CommentCreateRequest struct {
	CommentId   string
	PostId      int64     `json:"postId" validate:"required`
	Content     string    `json:"content" validate:"required`
	Password    string    `json:"password" validate:"required`
	MessageType string    `json:"messageType" validate:"required`
	Timestamp   time.Time `json:"timestamp" validate:"required`
}

type CommentLikeRequest struct {
	PostId      int64  `json:"postId" validate:"required`
	CommentId   string `json:"commentId" validate:"required`
	UserId      string `json:"userId" validate:"required`
	MessageType string `json:"messageType" validate:"required`
}

