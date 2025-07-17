package models

import "time"

type TodoList struct {
	Id        uint      `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	Todos     []Todo    `json:"todos"`
}
