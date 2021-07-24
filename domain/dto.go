package domain

type CommentRequest struct {
	Content  string
	Password string
}

type CommentLikeRequest struct {
	commentId int64
}

type QueueCreateRequest struct {
	queueName string
}
