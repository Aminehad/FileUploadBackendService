package api

import (
	"context"
	"log"
	"net/http"
	"upload-service/internal/app/upload-service/repository"

	"github.com/gin-gonic/gin"
)

type app struct {
	Router *gin.Engine
	port   string
}

func (a app) Start(context.Context) error {
	return a.Router.Run(":" + a.port)
}

func (a app) Stop() {}

type Configuration struct {
	// AppName    string
	// AppVersion string
	Log        *log.Logger
	Repository repository.Repository
}

// WE COULD VEW TOTALY DONE GIN.DEFAULT() BUT WE WANTED TO SHOW YOU HOW TO CREATE A NEW CUSTUM GIN ENGINE
func NewServer() *gin.Engine {
	// Create a new Gin engine without default middleware
	r := gin.New()

	// Add middleware
	r.Use(gin.Logger())   // Logs HTTP requests
	r.Use(gin.Recovery()) // Recovers from panics to avoid crashing

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	return r
}

// New instanciates app server
func New(c Configuration) *app {
	router := NewServer()

	h := &handler{
		log:        c.Log,
		repository: c.Repository,
	}

	files := router.Group("/v1/files")
	{
		files.POST("/auth", h.Login)

		files.Use(MiddlewareCheckLoginJWT())
		files.GET("/", h.ListFiles)
		files.POST("/upload", h.UploadFile)
		files.GET("/:fileID", h.DownloadFile)

	}

	return &app{router, "5051"}
}
