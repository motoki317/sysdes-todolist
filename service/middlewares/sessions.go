package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const sessionName = "session"

const (
	userIDKey = "userID"
)

const (
	ContextUserIDKey = "userID"
)

func IsLoggedIn(store sessions.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess, err := store.Get(c.Request, sessionName)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		iUserID, ok := sess.Values[userIDKey]
		userID, castOK := iUserID.(uint64)
		if !ok || !castOK {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(ContextUserIDKey, userID)
		c.Next()
	}
}

func SetLoginSession(c *gin.Context, store sessions.Store, userID uint64) error {
	sess, err := store.Get(c.Request, sessionName)
	if err != nil {
		sess, err = store.New(c.Request, sessionName)
		if err != nil {
			return err
		}
	}
	sess.Values[userIDKey] = userID
	return store.Save(c.Request, c.Writer, sess)
}
