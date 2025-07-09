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
	ResumeName string `json:"resumeName" binding:"required"`
}

const jobDescription = `
Job Title: Cloud Infrastructure Engineer (Entry-Level)
Location: San Francisco, CA (Hybrid or Remote Eligible)
Employment Type: Full-Time
Start Date: Immediate / Upon Graduation (May 2025)

About the Role:
We are seeking a Cloud Infrastructure Engineer to join our high-impact infrastructure team. This role is ideal for a recent graduate or early-career engineer who is passionate about cloud automation, infrastructure as code, and scalable distributed systems. You’ll play a key role in designing, building, and maintaining the infrastructure that powers our production environments. You'll work cross-functionally with software engineers and DevOps to streamline deployments, improve performance, and enhance observability.

Key Responsibilities:
Design, implement, and manage cloud-native infrastructure using Terraform, Kubernetes, and ArgoCD.

Automate service deployment and CI/CD pipelines using Github Actions, Docker, and Argo Rollouts.

Monitor system performance and error alerts using AWS CloudWatch, SNS, and log metric filters.

Support service migrations across AWS and GCP, using tools such as ECS, EKS, GKE, and ALB.

Maintain and optimize databases and data pipelines using PostgreSQL, ElasticSearch, and RDS.

Collaborate with backend teams to develop internal tools and dashboards that improve system visibility.

Drive performance optimizations and scalability improvements (e.g., homepage load times, query efficiency).

Troubleshoot production issues and resolve infrastructure-related incidents.

Required Skills & Qualifications:
B.S. in Computer Science or related technical field (Graduating May 2025 or recent graduate).

Strong programming skills in Go, Python, TypeScript, or JavaScript.

Hands-on experience with cloud platforms (AWS, GCP) and container orchestration tools (Kubernetes, Docker).

Familiarity with CI/CD pipelines, infrastructure as code (Terraform), and GitHub Actions.

Experience with monitoring and observability tools like CloudWatch, SNS, and ElasticSearch.

Strong understanding of DevOps principles, cloud security, and system reliability practices.

Nice to Have:
Previous internship experience at companies working with large-scale cloud infrastructure.

Exposure to multi-environment deployments (e.g., blue/green, canary) and tools like Argo Rollouts.

Familiarity with frontend frameworks like React.js or Next.js is a plus.

Contributions to open-source or academic hackathon awards (e.g., Most Innovative Project).

Why Join Us?
Work on cutting-edge infrastructure projects with a strong focus on automation and scalability.

Join a collaborative, forward-thinking engineering culture that values innovation and ownership.

Grow your cloud engineering skills through hands-on experience and mentorship.

Competitive salary, equity options, and professional development opportunities.`

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
		analyses, err := openai.AnalyzeResume(jobDescription, []*parser.File{file})
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
		fileObjects, err := gcloudClient.DownloadAllPDFs("user-resumes-hs-hackathon")
		if err != nil {
			c.JSON(500, gin.H{"success": false, "error": err.Error()})
			return
		}

		// Update the resume content
		parser.UpdateResumeContent(fileObjects)

		// Analyze the content and job description to determine a match
		analyses, err := openai.AnalyzeResume(jobDescription, fileObjects)
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
