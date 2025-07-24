package handlers

import (
	"net/http"
	"strconv"
	"todo/models"

	"github.com/gin-gonic/gin"
)

func getTodoLists(c *gin.Context) {
	todoLists, err := todoDataAccess.GetTodoLists()
	if err != nil {
		writeInternalServerError(c.Writer, err.Error())
		return
	}

	c.HTML(http.StatusOK, "todoLists", MapTodoListsTemplate(todoLists))
}

func createTodoList(c *gin.Context) {
	var todoList models.TodoList
	todoList, err := todoDataAccess.CreateTodoList(todoList)
	if err != nil {
		writeInternalServerError(c.Writer, err.Error())
		return
	}

	c.HTML(http.StatusOK, "todoList", MapTodoListTemplate(todoList))
}

func updateTodoList(c *gin.Context) {
	strId := c.PostForm("id")
	title := c.PostForm("title")

	if strId == "" {
		writeInternalServerError(c.Writer, "please provide a todolist id")
		return
	}

	id, err := strconv.Atoi(strId)
	if err != nil {
		writeInternalServerError(c.Writer, err.Error())
		return
	}

	todoDataAccess.UpdateTodoList(uint(id), title)
}

func deleteTodoList(c *gin.Context) {
	id := c.Query("todolistid")
	if id == "" {
		writeInternalServerError(c.Writer, "please provide a todolist id")
		return
	}

	idStr, err := strconv.Atoi(id)
	if err != nil {
		writeInternalServerError(c.Writer, "invalid id")
		return

	}
	err = todoDataAccess.DeleteTodoList(uint(idStr))
	if err != nil {
		writeInternalServerError(c.Writer, err.Error())
		return
	}
}
