package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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

// Generate random salt
func GenerateSalt(length int) (string, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return hex.EncodeToString(salt), nil
}

// Hash password with salt
func HashPasswordWithSalt(password, salt string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password + salt))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Verify password with salt
func VerifyPasswordWithSalt(password, salt, hashedPassword string) bool {
	return HashPasswordWithSalt(password, salt) == hashedPassword
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
	router.POST("/register", func(c *gin.Context) {
		userName := c.PostForm("userName")
		if userName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a username"})
			return
		}

		email := c.PostForm("email")
		if email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a email"})
			return
		}

		firstName := c.PostForm("firstName")
		if firstName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a first name"})
			return
		}

		lastName := c.PostForm("lastName")
		if lastName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a last name"})
			return
		}

		password := c.PostForm("password")
		if password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a password"})
			return
		}

		confirmPassword := c.PostForm("passwordConfirm")
		if confirmPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a password confirmation"})
			return
		}

		if password != confirmPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords dont match"})
			return
		}

		salt, err := GenerateSalt(16)
		if err != nil {
			log.Println(err.Error())
			writeInternalServerError(c.Writer, "Error registering user")
			return
		}

		hashedPassword := HashPasswordWithSalt(password, salt)
		_, err = tda.RegisterUser(userName, email, firstName, lastName, hashedPassword, salt)
		if err != nil {
			log.Println(err.Error())
			writeInternalServerError(c.Writer, "Error registering user")
			return
		}

		c.Header("HX-Redirect", "/login")
		c.Status(http.StatusNoContent)
	})

	router.POST("/login", func(c *gin.Context) {
		userName := c.PostForm("userName")
		if userName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a username"})
			return
		}
		password := c.PostForm("password")
		if password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a password"})
			return
		}

		user, err := tda.GetUserByUsername(userName)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "username or password was incorrect"})
			return
		}

		isValid := VerifyPasswordWithSalt(password, user.PasswordSalt, user.PasswordHash)
		if !isValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username or password was incorrect"})
			return
		}

		successMsg := fmt.Sprintf("The user %s successfully signed in", user.Username)
		log.Println(successMsg)
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
		authorized.POST(todosEndpoint, createTodo)
		authorized.PUT(todosEndpoint, updateTodo)
		authorized.DELETE(todosEndpoint, deleteTodo)
	}

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
