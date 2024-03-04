package models


import "time"

type Task struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description,omitempty"`
    Priority    int       `json:"priority"`
    DueDate     time.Time `json:"due_date"`
}
