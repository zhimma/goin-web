package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginController struct {
}

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"name": "马雄飞"})
}
