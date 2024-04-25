package files

import (
	"github.com/gin-gonic/gin"
	"madaurus/dev/material/app/models"
)

func CreateFile(c *gin.Context) {
	var file models.Files
	err := c.ShouldBindJSON(&file)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

}
