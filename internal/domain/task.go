package domain

type Task struct {
	Title       string `json:"title"`
	// Optional
	Description string `json:"description"`
	CreatedAt 	string `json:"created_at"`
}
