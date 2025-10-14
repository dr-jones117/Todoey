package handlers

import (
	"fmt"
	"strconv"
	"todo/models"

	"github.com/gin-gonic/gin"
)

func MapTodoFromRequestForm(c *gin.Context) (models.Todo, error) {
	var todo models.Todo

	formId, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		return todo, fmt.Errorf("invalid id")
	}

	todoListId, err := strconv.Atoi(c.PostForm("todolistid"))
	if err != nil {
		return todo, fmt.Errorf("invalid todolistid")
	}

	formCompleted := c.PostForm("completed")
	if formCompleted == "true" {
		todo.Completed = true
	} else {
		todo.Completed = false
	}

	todo.Id = uint(formId)
	todo.TodoListId = uint(todoListId)
	todo.Value = c.PostForm("value")

	return todo, nil
}
