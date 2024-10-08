package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func WebGuest() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, _ := store.Get(c.Request, "session")
		if auth, ok := session.Values["authenticated"].(bool); ok && auth {
			c.Redirect(http.StatusFound, "/dashboard")
			c.Abort()
			return
		}
		c.Next()
	}
}
