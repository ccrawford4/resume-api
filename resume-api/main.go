package main

import (
	"context"
	"resume-api/gcloud"
	"resume-api/openai"
	"resume-api/parser"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

// ResumeRequest represents the expected request body for the /new-resume endpoint
type ResumeRequest struct {
	ResumeName     string `json:"resumeName" binding:"required"`
	JobDescription string `json:"jobDescription" binding:"required"`
}

type AllResumeRequest struct {
	JobDescription string `json:"jobDescription" binding:"required"`
}

func main() {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))

	ctx := context.Background()
	gcloudClient, err := gcloud.NewClient(ctx)
	if err != nil {
		panic("could not create GCS client: " + err.Error())
	}

	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello world")
	})

	r.POST("/new-resume", func(c *gin.Context) {
		// Extract the resume name from the request body
		var req ResumeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"success": false, "error": "Invalid request body: " + err.Error()})
			return
		}
		resumeName := req.ResumeName

		// Download the resume file from GCS
		file, err := gcloudClient.DownloadFile("user-resumes-hs-hackathon", resumeName)
		if err != nil {
			c.JSON(500, gin.H{"success": false, "error": err.Error()})
			return
		}

		// Update the resume content by parsing the PDF
		parser.UpdateResumeContent([]*parser.File{file})

		// Fetch the analysis from OpenAI
		analyses, err := openai.AnalyzeResume(req.JobDescription, []*parser.File{file})
		if err != nil {
			c.JSON(500, gin.H{"success": false, "error": err.Error()})
			return
		}

		response := gin.H{
			"success":      true,
			"data":         analyses,
			"totalResumes": len(analyses),
		}
		c.JSON(200, response)
	})

	r.POST("/", func(c *gin.Context) {
		// Resume Request for all resumes
		var req AllResumeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"success": false, "error": "Invalid request body: " + err.Error()})
			return
		}

		// Download all resume files from GCS
		fileObjects, err := gcloudClient.DownloadAllPDFs("user-resumes-hs-hackathon")
		if err != nil {
			c.JSON(500, gin.H{"success": false, "error": err.Error()})
			return
		}

		// Update the resume content
		parser.UpdateResumeContent(fileObjects)

		// Analyze the content and job description to determine a match
		analyses, err := openai.AnalyzeResume(req.JobDescription, fileObjects)
		if err != nil {
			c.JSON(500, gin.H{"success": false, "error": err.Error()})
			return
		}

		response := gin.H{
			"success":      true,
			"data":         analyses,
			"totalResumes": len(analyses),
		}
		c.JSON(200, response)
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
