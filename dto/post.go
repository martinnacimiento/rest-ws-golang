package dto

// CreatePostRequest DTO

type CreatePostRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

// UpdatePostResponse DTO

type UpdateOnePostRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}
