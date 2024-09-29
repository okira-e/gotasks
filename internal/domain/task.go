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
	// Column      string `json:"column"`
	CreatedAt   string `json:"created_at"`
}

func NewTask(title string, description string) *Task {
	ret := new(Task)
	
	id := uuid.New()

	ret.Id = id.String()
	ret.Title = title
	ret.Description = description
	ret.CreatedAt = time.Now().UTC().String()
	
	return ret
}
