package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func TestList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "test 方法"})
}
