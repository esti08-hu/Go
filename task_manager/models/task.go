package models
import (
	"time"
)

// Task represents a task in the system
type Task struct {
	ID          string    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string `json:"status"`
}
