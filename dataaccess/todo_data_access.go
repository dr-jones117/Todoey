package dataaccess

import (
	"todo/models"
)

type TodoDataAccess interface {
	ConnectDataAccess() error
	DisconnectDataAccess() error

	// Todo
	CreateTodo(todo models.Todo) (models.Todo, error)
	UpdateTodo(todo models.Todo) (models.Todo, error)
	DeleteTodo(id uint) error

	// TodoList
	GetTodoLists() ([]models.TodoList, error)
	CreateTodoList(todoList models.TodoList) (models.TodoList, error)
	UpdateTodoList(todoListId uint, title string) error
	DeleteTodoList(id uint) error
}
