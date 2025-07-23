package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"todo/dataaccess"
	"todo/models"
	"todo/templates"

	"github.com/gin-gonic/gin"
)

var (
	tmpl           *template.Template
	todoDataAccess dataaccess.TodoDataAccess
)

func SetupHTTPHandlers(router *gin.Engine, tda dataaccess.TodoDataAccess) {
	todoDataAccess = tda

	router.LoadHTMLGlob("templates/**/*.html")

	router.Static("/css", "./css")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", nil)
	})

	// Todo Lists
	todoListEndpoint := "/todo-lists"
	router.GET(todoListEndpoint, getTodoLists)
	router.POST(todoListEndpoint, createTodoList)
	router.PUT(todoListEndpoint, updateTodoList)
	router.DELETE(todoListEndpoint, deleteTodoList)

	// Todos
	todosEndpoint := "/todos"
	router.POST(todosEndpoint, createTodo)
	router.PUT(todosEndpoint, updateTodo)
	router.DELETE(todosEndpoint, deleteTodo)
}

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

func writeInternalServerError(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusInternalServerError)
	log.Println(msg)
}
