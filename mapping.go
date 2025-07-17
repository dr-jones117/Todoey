package main

import (
	"fmt"
	"net/http"
	"strconv"
	"todo/models"
	"todo/templates"
)

func MapTodoFromRequestForm(r *http.Request) (models.Todo, error) {
	var todo models.Todo

	formId, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		return todo, fmt.Errorf("invalid id")
	}

	todoListId, err := strconv.Atoi(r.FormValue("todolistid"))
	if err != nil {
		return todo, fmt.Errorf("invalid todolistid")
	}

	formCompleted := r.FormValue("completed")
	if formCompleted == "true" {
		todo.Completed = true
	} else {
		todo.Completed = false
	}

	todo.Id = uint(formId)
	todo.TodoListId = uint(todoListId)
	todo.Value = r.FormValue("value")

	return todo, nil
}

func MapTodoListTemplate(todoList models.TodoList) templates.TodoListTemplateData {
	todoListTemplateData := templates.TodoListTemplateData{
		Id:        todoList.Id,
		Title:     todoList.Title,
		CreatedAt: todoList.CreatedAt,
	}

	todoListTemplateData.Todos = make([]templates.TodoTemplateData, len(todoList.Todos))
	for i, todo := range todoList.Todos {
		todoListTemplateData.Todos[i] = templates.TodoTemplateData{
			Todo: todo,
		}
	}

	return todoListTemplateData
}

func MapTodoListsTemplate(todoLists []models.TodoList) templates.TodoListsTemplateData {
	todoListsTemplateData := make([]templates.TodoListTemplateData, len(todoLists))
	for i, todoList := range todoLists {
		todoListsTemplateData[i] = MapTodoListTemplate(todoList)
	}

	return templates.TodoListsTemplateData{
		TodoLists: todoListsTemplateData,
	}
}
