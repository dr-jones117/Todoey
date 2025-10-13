package models

import "time"

type TodoCompletionFilter int

const (
	TodoCompleted TodoCompletionFilter = iota
	TodoIncomplete
	TodoNoApplyCompleted
)

type Todo struct {
	Id        uint   `json:"id"`
	Completed bool   `json:"completed"`
	Value     string `json:"value"`

	TodoListId    uint   `json:"todoListId"`
	TodoListTitle string `json:"todoListTitle"`

	CreatedAt            time.Time `json:"createdAt"`
	CompletedAt          time.Time `json:"completedAt"`
	CompletedAtFormatted string    `json:"completedAtFormatted"`
}
