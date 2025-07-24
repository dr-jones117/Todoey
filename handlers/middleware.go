package handlers

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userId := session.Get("user_id")

		log.Println("We have hit the auth function")
		log.Println(userId)

		if userId == nil {
			log.Println("nil hit")
			c.Header("HX-Redirect", "/login")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user_id", userId)
		c.Next()
	}
}
