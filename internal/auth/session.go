package auth

import (
	"github.com/elwafa/billion-data/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

func AdminLoginWeb(c *gin.Context, user entities.User) error {
	session, err := store.Get(c.Request, "session")
	if err != nil {
		return err
	}
	session.Values["authenticated"] = true
	session.Values["userId"] = user.Id
	return session.Save(c.Request, c.Writer)
}
