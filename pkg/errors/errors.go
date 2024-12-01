package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleBadRequest(c *gin.Context, err error) {
	c.Error(err)
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
func HandleInternalServerError(c *gin.Context, err error) {
	c.Error(err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
func HandleNotFound(c *gin.Context, err error) {
	c.Error(err)
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
}
func HandleUnauthorized(c *gin.Context, err error) {
	c.Error(err)
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
}
func HandleForbidden(c *gin.Context, err error) {
	c.Error(err)
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
}
func HandleConflict(c *gin.Context, err error) {
	c.Error(err)
	c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
}
