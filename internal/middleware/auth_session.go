package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"net/http"
)

// Initialize a session store
var store = sessions.NewCookieStore([]byte("super-secret-key"))

// Middleware to check if user is authenticated
func WebAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, _ := store.Get(c.Request, "session")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}
		c.Set("userId", session.Values["userId"])
		c.Next()
	}
}
