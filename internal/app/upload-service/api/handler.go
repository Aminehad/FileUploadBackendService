package api

import (
	"log"
	"upload-service/internal/app/upload-service/repository"

	"github.com/gin-gonic/gin"
)

type handler struct {
	log        *log.Logger
	repository repository.Repository
}

func (h *handler) UploadFile(c *gin.Context) {

}

func (h *handler) ListFiles(c *gin.Context) {
}
func (h *handler) DownloadFile(c *gin.Context) {
}
