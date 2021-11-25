package service

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"

	"todolist.go/service/middlewares"
)

type getMeResponse struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *Handlers) GetMe(c *gin.Context) {
	user := middlewares.GetUser(c)

	c.JSON(http.StatusOK, getMeResponse{
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

type editMeRequest struct {
	Name null.String `json:"name"`
}

func (h *Handlers) EditMe(c *gin.Context) {
	user := middlewares.GetUser(c)

	var req editMeRequest
	if err := c.Bind(&req); err != nil || (req.Name.Valid && req.Name.String == "") {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if req.Name.Valid {
		if _, err := h.db.Exec("UPDATE `users` SET `name` = ? WHERE `id` = ?", req.Name.String, user.ID); err != nil {
			// Check name conflict
			if merr, ok := err.(*mysql.MySQLError); ok && merr.Number == mysqlErrorDuplicateEntry {
				c.AbortWithStatus(http.StatusConflict)
				return
			}
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	c.Status(http.StatusNoContent)
}

func (h *Handlers) DeleteMe(c *gin.Context) {
	user := middlewares.GetUser(c)

	if _, err := h.db.Exec("UPDATE `users` SET `deleted_at` = NOW(6) WHERE `id` = ?", user.ID); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := middlewares.RevokeSession(c, h.store); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

type editMyPasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h *Handlers) EditMyPassword(c *gin.Context) {
	user := middlewares.GetUser(c)

	var req editMyPasswordRequest
	if err := c.Bind(&req); err != nil || req.OldPassword == "" || req.NewPassword == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Check old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcryptDefaultCost)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if _, err := h.db.Exec("UPDATE `users` SET `password` = ? WHERE `id` = ?", newHash, user.ID); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
