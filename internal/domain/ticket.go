package domain

type Ticket struct {
	Title       string `json:"title"`
	// Optional
	Description string `json:"description"`
}
