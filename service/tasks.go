package service

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v4"

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

type getTasksRequest struct {
	Title  string    `form:"title"`
	Done   null.Bool `form:"done"`
	Limit  null.Int  `form:"limit"`
	Offset null.Int  `form:"offset"`
}

type getTasksResponse struct {
	Count int             `json:"count"`
	Tasks []*taskResponse `json:"tasks"`
}

func (h *Handlers) GetTasks(c *gin.Context) {
	user := middlewares.GetUser(c)

	var req getTasksRequest
	if err := c.Bind(&req); err != nil || (req.Limit.Valid && req.Limit.Int64 <= 0) || (req.Offset.Valid && req.Offset.Int64 < 0) || (!req.Limit.Valid && req.Offset.Valid) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var count int
	if err := h.db.Get(&count, "SELECT COUNT(*) FROM `tasks` WHERE `user_id` = ? AND `deleted_at` IS NULL", user.ID); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	query := "SELECT * FROM `tasks` WHERE "
	conditions := make([]string, 0, 6)
	args := make([]interface{}, 0, 6)

	conditions = append(conditions, "`user_id` = ?", "`deleted_at` IS NULL")
	args = append(args, user.ID)

	if req.Title != "" {
		conditions = append(conditions, "`title` LIKE ?")
		args = append(args, "%"+req.Title+"%")
	}
	if req.Done.Valid {
		conditions = append(conditions, "`done` = ?")
		args = append(args, req.Done.Bool)
	}

	query += strings.Join(conditions, " AND ")

	query += " ORDER BY `id`"
	if req.Limit.Valid {
		query += " LIMIT ?"
		args = append(args, req.Limit.Int64)

		if req.Offset.Valid {
			query += " OFFSET ?"
			args = append(args, req.Offset.Int64)
		}
	}

	var tasks []*db.Task
	if err := h.db.Select(&tasks, query, args...); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, getTasksResponse{
		Count: count,
		Tasks: formatTasks(tasks),
	})
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

func (h *Handlers) GetTask(c *gin.Context) {
	task := middlewares.GetTask(c)

	c.JSON(http.StatusOK, formatTask(task))
}

type editTaskRequest struct {
	Title null.String `json:"title"`
	Done  null.Bool   `json:"done"`
}

func (h *Handlers) EditTask(c *gin.Context) {
	task := middlewares.GetTask(c)

	var req editTaskRequest
	if err := c.Bind(&req); err != nil || (req.Title.Valid && req.Title.String == "") {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	query := "UPDATE `tasks` SET "
	clauses := make([]string, 0)
	args := make([]interface{}, 0)
	update := false
	if req.Title.Valid {
		clauses = append(clauses, "`title` = ?")
		args = append(args, req.Title.String)
		update = true
	}
	if req.Done.Valid {
		clauses = append(clauses, "`done` = ?")
		args = append(args, req.Done.Bool)
		update = true
	}
	query += strings.Join(clauses, ", ")
	query += " WHERE `id` = ?"
	args = append(args, task.ID)

	if update {
		if _, err := h.db.Exec(query, args...); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	c.Status(http.StatusNoContent)
}

func (h *Handlers) DeleteTask(c *gin.Context) {
	task := middlewares.GetTask(c)

	if _, err := h.db.Exec("UPDATE `tasks` SET `deleted_at` = NOW(6) WHERE `id` = ?", task.ID); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
