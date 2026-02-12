package request

// SaveBookRequest is the request body for saving a book
type SaveBookRequest struct {
	Title        string `json:"title"`
	Author       string `json:"author"`
	BookISBN     string `json:"book_isbn"`
	ThumbnailURL string `json:"thumbnail_url"`
	Status       string `json:"status"`
}

// UpdateBookRequest is the request body for updating a book
type UpdateBookRequest struct {
	Title        string `json:"title"`
	Author       string `json:"author"`
	BookISBN     string `json:"book_isbn"`
	ThumbnailURL string `json:"thumbnail_url"`
	Status       string `json:"status"`
}

// SaveBookReviewRequest is the request body for saving a book review
type SaveBookReviewRequest struct {
	BookID   string `json:"book_id"`
	Content  string `json:"content"`
	Rating   int    `json:"rating"`
	IsPublic bool   `json:"is_public"`
}
