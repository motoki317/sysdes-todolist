package service

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"

	"todolist.go/db"
	"todolist.go/service/middlewares"
)

const (
	bcryptDefaultCost        = 10
	mysqlErrorDuplicateEntry = 1062
)

var (
	bcryptNonExistPassword []byte
)

func init() {
	hash, err := bcrypt.GenerateFromPassword([]byte(""), bcryptDefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	bcryptNonExistPassword = hash
}

type signUpRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h *Handlers) SignUp(c *gin.Context) {
	var req signUpRequest
	if err := c.Bind(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Password == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptDefaultCost)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if _, err := h.db.Exec("INSERT INTO `users` (`name`, `password`) VALUES (?, ?)", req.Name, hashedPass); err != nil {
		// Check name conflict
		if merr, ok := err.(*mysql.MySQLError); ok && merr.Number == mysqlErrorDuplicateEntry {
			c.AbortWithStatus(http.StatusConflict)
			return
		}
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusCreated)
}

type loginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h *Handlers) Login(c *gin.Context) {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Password == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var (
		user   db.User
		hash   []byte
		exists bool
	)
	err := h.db.Get(&user, "SELECT * FROM `users` WHERE `name` = ? AND `deleted_at` IS NULL", req.Name)
	if err != nil && err != sql.ErrNoRows {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	} else if err == sql.ErrNoRows {
		hash = bcryptNonExistPassword // ユーザーが存在しない場合も bcrypt を通す
	} else {
		hash = []byte(user.Password)
		exists = true
	}

	if err := bcrypt.CompareHashAndPassword(hash, []byte(req.Password)); err != nil || !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// login
	if err := middlewares.SetLoginSession(c, h.store, user.ID); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handlers) Logout(c *gin.Context) {
	if err := middlewares.RevokeSession(c, h.store); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
