package domain

import "time"

type Task struct {
	Title       string `json:"title"`
	// Optional
	Description string `json:"description"`
	CreatedAt 	string `json:"created_at"`
}

func NewTask(title string, description string) Task {
	return Task {
		Title: title,
		Description: description,
		CreatedAt: time.Now().UTC().String(),
	}
}