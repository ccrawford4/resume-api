package main

import (
	"resume-api/parser"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello world")
	})

	r.POST("/", func(c *gin.Context) {
		resumeContent := parser.ParseResume()
		c.JSON(200, resumeContent)
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
