package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TestList(c *gin.Context) {
	fmt.Println(c.Get("UUID"))
	fmt.Println(c.Get("UID"))
	c.JSON(http.StatusOK, gin.H{"message": "test 方法"})
}
