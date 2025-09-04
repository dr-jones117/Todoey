package handlers

import (
	"fmt"
	"log"
	"net/http"
	"todo/handlers/auth"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
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

	salt, err := auth.GenerateSalt(16)
	if err != nil {
		log.Println(err.Error())
		writeInternalServerError(c.Writer, "Error registering user")
		return
	}

	hashedPassword := auth.HashPasswordWithSalt(password, salt)
	_, err = todoDataAccess.RegisterUser(userName, email, firstName, lastName, hashedPassword, salt)
	if err != nil {
		log.Println(err.Error())
		writeInternalServerError(c.Writer, "Error registering user")
		return
	}

	c.Header("HX-Redirect", "/login")
	c.Status(http.StatusNoContent)
}

func Login(c *gin.Context) {
	userName := c.PostForm("userName")
	if userName == "" {
		c.Data(http.StatusBadRequest, "text/html", []byte("Please provide a username"))
		return
	}
	password := c.PostForm("password")
	if password == "" {
		c.Data(http.StatusBadRequest, "text/html", []byte("Please provide a password"))
		return
	}

	user, err := todoDataAccess.GetUserByUsername(userName)
	if err != nil {
		log.Println(err.Error())
		c.Data(http.StatusBadRequest, "text/html", []byte("Username or password was incorrect"))
		return
	}

	isValid := auth.VerifyPasswordWithSalt(password, user.PasswordSalt, user.PasswordHash)
	if !isValid {
		log.Println("Couldn't verify the password for user", userName)
		c.Data(http.StatusBadRequest, "text/html", []byte("Username or password was incorrect"))
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.Id)
	session.Save()

	successMsg := fmt.Sprintf("The user %s successfully signed in", user.Username)
	log.Println(successMsg)

	c.Header("HX-Redirect", "/")
	c.Status(http.StatusOK)
}
