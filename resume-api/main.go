package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello world")
	})

	r.POST("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "POST received"})
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
