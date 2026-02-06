package dto

// CursorPagination represents cursor-based pagination metadata.
// Cursor is the "after" cursor used to fetch the current page.
// NextCursor is the cursor to be used to fetch the next page.
type CursorPagination struct {
	Limit       int     `json:"limit"`
	Cursor      *string `json:"cursor,omitempty"`
	NextCursor  *string `json:"next_cursor,omitempty"`
	HasNextPage bool    `json:"has_next_page"`
}

