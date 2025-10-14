package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"todo/models"
	"todo/templates"

	"github.com/gin-gonic/gin"
)

func getTodos(c *gin.Context) {
	todos, err := todoDataAccess.GetTodos(models.TodoNoApplyCompleted)

	if err != nil {
		writeInternalServerError(c.Writer, err.Error())
		return
	}

	c.HTML(http.StatusOK, "todos", gin.H{
		"todos": todos,
	})
}

func getHistoricalTodos(c *gin.Context) {
	todos, err := todoDataAccess.GetTodos(models.TodoCompleted)

	if err != nil {
		writeInternalServerError(c.Writer, err.Error())
		return
	}

	for i, todo := range todos {
		formattedTime := todo.CompletedAt.Local().Format("03:04:05 PM 01/02/2006")

		if todo.CompletedAt.IsZero() {
			formattedTime = "N/A"
		}

		todos[i].CompletedAtFormatted = formattedTime
	}

	c.HTML(http.StatusOK, "historical-todos", gin.H{
		"todos": todos,
	})
}

func deleteHistoricalTodos(c *gin.Context) {
	err := todoDataAccess.DeleteHistoricalTodos()

	if err != nil {
		writeInternalServerError(c.Writer, err.Error())
		c.Data(http.StatusBadRequest, "text/html", []byte("Unable to clear your history at this time. Try again later."))
		return
	}

	c.Data(http.StatusOK, "text/html", []byte("Successfully cleared your task history!"))
}

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

	c.HTML(http.StatusOK, "todo", todo)
}

func updateTodo(c *gin.Context) {
	todo, err := MapTodoFromRequestForm(c)
	if err != nil {
		writeInternalServerError(c.Writer, err.Error())
		return
	}

	originalTodo, err := todoDataAccess.GetTodoById(todo.Id)
	if err != nil {
		writeBadRequestError(c.Writer, "Unable to find todo with that id", err.Error())
		return
	}

	if todo.Completed && !originalTodo.Completed {
		todo.CompletedAt = time.Now().UTC()
	}

	todo, err = todoDataAccess.UpdateTodo(todo)
	if err != nil {
		writeInternalServerError(c.Writer, err.Error())
		return
	}

	log.Println(todo)

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
