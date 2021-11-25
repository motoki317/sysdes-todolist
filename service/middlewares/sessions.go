package middlewares

import (
	"database/sql"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	db2 "todolist.go/db"
)

const (
	sessionName      = "session"
	sessionUserIDKey = "userID"
)

func IsLoggedIn(db *sqlx.DB, store sessions.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess, err := store.Get(c.Request, sessionName)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		iUserID, ok := sess.Values[sessionUserIDKey]
		userID, castOK := iUserID.(uint64)
		if !ok || !castOK {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user db2.User
		if err := db.Get(&user, "SELECT * FROM `users` WHERE `id` = ? AND `deleted_at` IS NULL", userID); err != nil && err != sql.ErrNoRows {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		} else if err == sql.ErrNoRows {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(ctxUserKey, &user)
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
	sess.Values[sessionUserIDKey] = userID
	return store.Save(c.Request, c.Writer, sess)
}

func RevokeSession(c *gin.Context, store sessions.Store) error {
	sess, err := store.Get(c.Request, sessionName)
	if err != nil {
		sess, err = store.New(c.Request, sessionName)
		if err != nil {
			return err
		}
	}
	sess.Values = make(map[interface{}]interface{})
	return store.Save(c.Request, c.Writer, sess)
}

func GetUser(c *gin.Context) *db2.User {
	v, _ := c.Get(ctxUserKey)
	return v.(*db2.User)
}
