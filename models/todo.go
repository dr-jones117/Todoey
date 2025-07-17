package models

import "time"

type Todo struct {
	Id         uint      `json:"id"`
	Completed  bool      `json:"completed"`
	Value      string    `json:"value"`
	CreatedAt  time.Time `json:"createdAt"`
	TodoListId uint      `json:"todoListId"`
}
