package service

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"todolist.go/db"
	"todolist.go/service/middlewares"
)

type taskResponse struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func formatTask(task *db.Task) *taskResponse {
	return &taskResponse{
		ID:        task.ID,
		Title:     task.Title,
		Done:      task.Done,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
}

func formatTasks(tasks []*db.Task) []*taskResponse {
	ret := make([]*taskResponse, len(tasks))
	for i := 0; i < len(tasks); i++ {
		ret[i] = formatTask(tasks[i])
	}
	return ret
}

func (h *Handlers) GetTasks(c *gin.Context) {
	user := middlewares.GetUser(c)

	var tasks []*db.Task
	if err := h.db.Select(&tasks, "SELECT * FROM `tasks` WHERE `user_id` = ? AND `deleted_at` IS NULL ORDER BY `id`", user.ID); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, formatTasks(tasks))
}

type createTaskRequest struct {
	Title string `json:"title"`
}

func (h *Handlers) CreateTask(c *gin.Context) {
	user := middlewares.GetUser(c)

	var req createTaskRequest
	if err := c.Bind(&req); err != nil || req.Title == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res, err := h.db.Exec("INSERT INTO `tasks` (`user_id`, `title`) VALUES (?, ?)", user.ID, req.Title)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var task db.Task
	if err := h.db.Get(&task, "SELECT * FROM `tasks` WHERE `id` = ?", lastInsertID); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, formatTask(&task))
}
