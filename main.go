package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"todolist.go/db"
	"todolist.go/service"
	"todolist.go/service/middlewares"
)

const (
	defaultPort = 8000
)

func main() {
	// initialize DB connection
	dsn := db.DefaultDSN(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	dbConn, err := db.Connect(dsn)
	if err != nil {
		log.Fatal(err)
	}

	// port to listen on
	var port int
	if strPort := os.Getenv("PORT"); strPort != "" {
		_port, err := strconv.Atoi(strPort)
		if err != nil {
			log.Fatal(err)
		}
		port = _port
	} else {
		port = defaultPort
	}

	// initialize Gin engine
	engine := gin.Default()

	// session store
	sessSecret := os.Getenv("SESSION_SECRET")
	if sessSecret == "" {
		log.Fatal("SESSION_SECRET is empty")
	}
	store := cookie.NewStore([]byte(sessSecret))
	engine.Use(sessions.Sessions("cookie-session", store))
	engine.Use(gin.ErrorLogger())

	// routing
	h := service.NewHandlers(dbConn, store)
	api := engine.Group("/api", middlewares.IsLoggedIn(dbConn, store))
	{
		apiUsers := api.Group("/users")
		{
			apiUsers.GET("/me", h.GetMe)
			apiUsers.PATCH("/me", h.EditMe)
			apiUsers.DELETE("/me", h.DeleteMe)
			apiUsers.PUT("/me/password", h.EditMyPassword)
		}
		apiTasks := api.Group("/tasks")
		{
			apiTasks.GET("", h.GetTasks)
			apiTasks.POST("", h.CreateTask)
			apiTaskID := apiTasks.Group("/:taskID", middlewares.RetrieveTask(dbConn))
			{
				apiTaskID.GET("", h.GetTask)
				apiTaskID.PATCH("", h.EditTask)
				apiTaskID.DELETE("", h.DeleteTask)
			}
		}
	}
	apiNoAuth := engine.Group("/api")
	{
		apiNoAuth.POST("/signup", h.SignUp)
		apiNoAuth.POST("/login", h.Login)
		apiNoAuth.POST("/logout", h.Logout)
	}

	// start server
	if err := engine.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal(err)
	}
}
