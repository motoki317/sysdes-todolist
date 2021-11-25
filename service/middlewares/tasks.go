package middlewares

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	db2 "todolist.go/db"
)

func RetrieveTask(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID, err := strconv.Atoi(c.Param("taskID"))
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		user := GetUser(c)
		var task db2.Task
		if err := db.Get(&task, "SELECT * FROM `tasks` WHERE `id` = ? AND `user_id` = ? AND `deleted_at` IS NULL", taskID, user.ID); err != nil && err != sql.ErrNoRows {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		} else if err == sql.ErrNoRows {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.Set(ctxTaskKey, &task)

		c.Next()
	}
}

func GetTask(c *gin.Context) *db2.Task {
	v, _ := c.Get(ctxTaskKey)
	return v.(*db2.Task)
}
