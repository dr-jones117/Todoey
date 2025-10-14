package dataaccess

import (
	"todo/models"
)

type TodoDataAccess interface {
	ConnectDataAccess() error
	DisconnectDataAccess() error

	// Users
	RegisterUser(userName string, email string, firstName string, lastName string, hashedPassword string, salt string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)

	// TodoList
	GetTodoLists() ([]models.TodoList, error)
	CreateTodoList(todoList models.TodoList) (models.TodoList, error)
	UpdateTodoList(todoListId uint, title string) error
	DeleteTodoList(id uint) error

	// Todo
	GetTodos(status models.TodoCompletionFilter) ([]models.Todo, error)
	GetTodoById(id uint) (models.Todo, error)
	CreateTodo(todo models.Todo) (models.Todo, error)
	UpdateTodo(todo models.Todo) (models.Todo, error)
	DeleteTodo(id uint) error

	DeleteHistoricalTodos() error
}
