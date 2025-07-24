package handlers

import (
	"net/http"
	"strconv"
	"todo/models"
	"todo/templates"

	"github.com/gin-gonic/gin"
)

func createTodo(c *gin.Context) {
	var todo models.Todo
	var err error
	var todoTemplateData templates.TodoTemplateData

	strTodoListId := c.PostForm("todolistid")
	if strTodoListId == "" {
		writeInternalServerError(c.Writer, "No todo list id was supplied")
		return
	}

	todoListId, err := strconv.Atoi(strTodoListId)
	if err != nil {
		writeInternalServerError(c.Writer, "invalid todo list id")
	}

	setFocus := c.Query("setFocus")
	if setFocus != "" {
		todoTemplateData.FocusInput = true
	}

	todo.TodoListId = uint(todoListId)
	todo, err = todoDataAccess.CreateTodo(todo)
	if err != nil {
		writeInternalServerError(c.Writer, err.Error())
		return
	}

	todoTemplateData.Todo = todo
	c.HTML(http.StatusOK, "todo", todoTemplateData)
}

func updateTodo(c *gin.Context) {
	todo, err := MapTodoFromRequestForm(c)
	if err != nil {
		writeInternalServerError(c.Writer, err.Error())
		return
	}

	todo, err = todoDataAccess.UpdateTodo(todo)
	if err != nil {
		writeInternalServerError(c.Writer, err.Error())
		return
	}

	c.HTML(http.StatusOK, "todo", todo)
}

func deleteTodo(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		writeInternalServerError(c.Writer, "Please provide an id")
		return
	}

	paramId, err := strconv.Atoi(id)
	if err != nil {
		writeInternalServerError(c.Writer, err.Error())
		return
	}

	idTofind := uint(paramId)
	if err := todoDataAccess.DeleteTodo(idTofind); err != nil {
		writeInternalServerError(c.Writer, err.Error())
		return
	}
}
