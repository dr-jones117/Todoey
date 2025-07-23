package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userId := session.Get("user_id")

		if userId == nil {
			c.Header("HX-Redirect", "/login")
			c.Status(http.StatusNoContent)
			return
		}

		c.Set("user_id", userId)
		c.Next()
	}
}
