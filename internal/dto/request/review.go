package request

// CreateReviewRequest is the request body for creating a review
type CreateReviewRequest struct {
	Content  string `json:"content"`
	Rating   int    `json:"rating"`
	IsPublic bool   `json:"is_public"`
}

// UpdateReviewRequest is the request body for updating a review
type UpdateReviewRequest struct {
	Content  *string `json:"content,omitempty"`
	Rating   *int    `json:"rating,omitempty"`
	IsPublic *bool   `json:"is_public,omitempty"`
}
