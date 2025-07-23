package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("hello from auth!")
		log.Println(c.Cookie("sessionId"))
		log.Println(c.Request.URL.Path)
		c.Next()
	}
}
