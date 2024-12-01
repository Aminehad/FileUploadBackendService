package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"upload-service/internal/app/upload-service/models"

	"upload-service/internal/app/upload-service/repository"
	"upload-service/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type handler struct {
	log        *log.Logger
	repository repository.Repository
}

func (h *handler) UploadFile(c *gin.Context) {

	// can not accept files larger than X MB
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		errors.HandleBadRequest(c, err)
		return
	}
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		errors.HandleBadRequest(c, err)
		return
	}
	defer file.Close()

	os.MkdirAll("./uploads", os.ModePerm)
	filePath := fmt.Sprintf("./uploads/%s", fileHeader.Filename)

	// save the file to disk/ local file system
	// TODO : change this to use a storage service S3 Minio/AWS
	uploadedFile, err := os.Create(filePath)
	if err != nil {
		errors.HandleInternalServerError(c, err)
		return
	}
	defer uploadedFile.Close()

	_, err = io.Copy(uploadedFile, file)
	if err != nil {
		errors.HandleInternalServerError(c, err)
		return
	}
	fileMetaData := models.File{
		Name: fileHeader.Filename,
		Url:  filePath,
	}

	var fileID *uuid.UUID
	fileID, err = h.repository.SaveFileMetaData(fileMetaData)
	if err != nil {
		if err == repository.ErrDuplicatedKeyUniqueConstraint {
			errors.HandleConflict(c, err)
			return
		}
		errors.HandleInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "file uploaded successfully", "fileID": fileID})
}

func (h *handler) ListFiles(c *gin.Context) {
	var f filters
	if err := c.ShouldBindQuery(&f); err != nil {
		errors.HandleBadRequest(c, fmt.Errorf("failed to bind query params: %w", err))
		return
	}
	var filters repository.FileFilters

	if f.FileID != "" {
		fileID, err := uuid.Parse(f.FileID)
		if err != nil {
			errors.HandleBadRequest(c, fmt.Errorf("invalid fileID: %w", err))
			return
		}
		filters.ID = fileID
	}
	files, err := h.repository.ListFiles(filters)
	if err != nil {
		if err == repository.ErrRecordNotFound {
			errors.HandleNotFound(c, fmt.Errorf("failed to list files : %w", err))
			return
		}
		errors.HandleInternalServerError(c, fmt.Errorf("failed to list files : %w", err))
		return
	}
	for i := range files {
		files[i].Url = fmt.Sprintf("/v1/files/%s", files[i].ID)
	}
	c.JSON(http.StatusOK, files)
}
func (h *handler) DownloadFile(c *gin.Context) {
}

type filters struct {
	FileID   string `form:"fileID"`
	Filename string `form:"filename"`
}

func (h *handler) UploadFileRawText(c *gin.Context) {
	rawText, err := io.ReadAll(c.Request.Body)
	if err != nil {
		errors.HandleBadRequest(c, err)
		return
	}
	defer c.Request.Body.Close()

	fileNameTmp := fmt.Sprintf("document_%s.txt", time.Now().Format("20060102_150405"))
	os.MkdirAll("./uploads", os.ModePerm)
	filePath := fmt.Sprintf("./uploads/%s", fileNameTmp)

	err = os.WriteFile(filePath, rawText, os.ModePerm)
	if err != nil {
		errors.HandleInternalServerError(c, err)
		return
	}
	fileMetaData := models.File{
		Name: fileNameTmp,
		Url:  filePath,
	}

	var fileID *uuid.UUID
	fileID, err = h.repository.SaveFileMetaData(fileMetaData)
	if err != nil {
		if err == repository.ErrDuplicatedKeyUniqueConstraint {
			errors.HandleConflict(c, err)
			return
		}
		errors.HandleInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "file uploaded successfully", "fileID": fileID})
}
