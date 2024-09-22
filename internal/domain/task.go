package domain

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	// Optional
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

func NewTask(title string, description string) Task {
	id := uuid.New()

	return Task{
		Id:          id.String(),
		Title:       title,
		Description: description,
		CreatedAt:   time.Now().UTC().String(),
	}
}
