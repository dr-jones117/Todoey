package templates

import (
	"time"
	"todo/models"
)

type TodoListsTemplateData struct {
	TodoLists []TodoListTemplateData
}

type TodoListTemplateData struct {
	Id        uint
	Title     string
	Todos     []TodoTemplateData
	CreatedAt time.Time
}

type TodoTemplateData struct {
	Todo       models.Todo
	FocusInput bool
}
