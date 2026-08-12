package handler

import "github.com/gin-gonic/gin"

func AdminLogin(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Admin login successful"})
}
