package entity

import "time"

// Example is a placeholder entity that demonstrates the shape expected by
// every other layer (repository, service, HTTP handler). It is intentionally
// generic: rename/replace it with your real domain entity (Invoice, Client,
// Campaign, ...) when you fork this repo. Keep the same conventions:
//   - `json` tags for the HTTP layer.
//   - `firestore` tags only if you plug in the firestore repository.
//   - A dedicated *Response struct for paginated list endpoints.
type Example struct {
	CreatedAt   time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" firestore:"updatedAt"`
	ID          string    `json:"id" firestore:"-"`
	Name        string    `json:"name" binding:"required" firestore:"name"`
	Description string    `json:"description" firestore:"description"`
}

// ExamplesResponse is the response for the paginated list endpoint.
type ExamplesResponse struct {
	Items      []*Example `json:"items"`
	TotalItems int        `json:"totalItems"`
}
