package domain

import "time"

type CommentCreateRequest struct {
	PostId      int64
	Content     string
	Password    string
	MessageType string
	Timestamp   time.Time
}

type CommentLikeRequest struct {
	PostId      int64
	CommentId   string
	userId      string
	MessageType string
}

type CommentUpdateRequest struct {
	PostId      int64
	CommentId   string
	Content     string
	MessageType string
	Timestamp   time.Time
}
