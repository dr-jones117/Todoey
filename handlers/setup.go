package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
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

func writeBadRequestError(w http.ResponseWriter, msg string, errorMessage string) {
	http.Error(w, msg, http.StatusBadRequest)
	log.Println(errorMessage)
}

func getTemplateFileList(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".html" {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func SetupHTTPHandlers(router *gin.Engine, tda dataaccess.TodoDataAccess) {
	todoDataAccess = tda

	// router.LoadHTMLGlob("templates/**/*.html")
	templateFilePaths, err := getTemplateFileList("templates")
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}

	router.LoadHTMLFiles(templateFilePaths...)

	//Public Static Files
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.File("favicon.ico")
	})

	router.Static("/css", "./css")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", gin.H{
			"showLogoutButton": true,
		})
	})

	router.GET("/history", func(c *gin.Context) {
		c.HTML(http.StatusOK, "history", gin.H{
			"showLogoutButton": true,
		})
	})

	// Auth
	// router.POST("/register", Register)
	router.POST("/login", Login)
	router.POST("/logout", Logout)
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login", gin.H{
			"showLogoutButton": false,
		})
	})

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
		authorized.GET(todosEndpoint, getTodos)
		authorized.POST(todosEndpoint, createTodo)
		authorized.PUT(todosEndpoint, updateTodo)
		authorized.DELETE(todosEndpoint, deleteTodo)

		authorized.GET("/historical-todos", getHistoricalTodos)
	}
}
