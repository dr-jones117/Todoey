package handlers

import (
	"log"
	"net/http"
	"todo/dataaccess"

	"github.com/gin-gonic/gin"
)

var (
	todoDataAccess dataaccess.TodoDataAccess
)

func writeInternalServerError(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusInternalServerError)
	log.Println(msg)
}

func SetupHTTPHandlers(router *gin.Engine, tda dataaccess.TodoDataAccess) {
	todoDataAccess = tda

	router.LoadHTMLGlob("templates/**/*.html")

	//Public Static Files
	router.Static("/css", "./css")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", nil)
	})

	// Auth
	router.POST("/register", Register)
	router.POST("/login", Login)

	authorized := router.Group("/")
	authorized.Use(AuthRequired())
	{
		// Todo Lists
		todoListEndpoint := "/todo-lists"
		authorized.GET(todoListEndpoint, getTodoLists)
		authorized.POST(todoListEndpoint, createTodoList)
		authorized.PUT(todoListEndpoint, updateTodoList)
		authorized.DELETE(todoListEndpoint, deleteTodoList)

		// Todos
		todosEndpoint := "/todos"
		authorized.POST(todosEndpoint, createTodo)
		authorized.PUT(todosEndpoint, updateTodo)
		authorized.DELETE(todosEndpoint, deleteTodo)
	}
}
