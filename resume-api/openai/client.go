package openai

import (
	context "context"
	"encoding/json"
	"fmt"
	"os"
	"resume-api/parser"
	"sync"

	openai "github.com/sashabaranov/go-openai"
)

type ResumeAnalysis struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	UploadDate      string         `json:"uploadDate"`
	FileUrl         string         `json:"fileUrl"`
	MatchPercentage int            `json:"matchPercentage"`
	Insights        map[string]any `json:"insights"`
}

func AnalyzeResume(jobDescription string, resumes []*parser.File) ([]ResumeAnalysis, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}
	client := openai.NewClient(apiKey)
	ctx := context.Background()

	results := []ResumeAnalysis{}
	id := 1

	var wg sync.WaitGroup
	for _, resume := range resumes {
		content := string(resume.Content)

		fmt.Printf("resume content: %s\n", content)
		name := resume.Name

		prompt := "Given the following job description: '" + jobDescription + "' and the following resume: '" + content + "', analyze the resume for match percentage, strengths, improvements, and missing skills. Respond in JSON with fields: matchPercentage, insights (with strengths, improvements, missingSkills as arrays of strings). NO MARKDOWN JUST PLAIN JSON"

		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{Role: openai.ChatMessageRoleSystem, Content: "You are a helpful assistant that analyzes resumes for job fit."},
					{Role: openai.ChatMessageRoleUser, Content: prompt},
				},
			})
			if err != nil || len(resp.Choices) == 0 {
				return
			}
			// Parse the response JSON
			var analysis struct {
				MatchPercentage int            `json:"matchPercentage"`
				Insights        map[string]any `json:"insights"`
			}
			err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &analysis)
			if err != nil {
				return
			}
			results = append(results, ResumeAnalysis{
				ID:              fmt.Sprintf("%d", id),
				Name:            name,
				UploadDate:      resume.CreatedAt,
				FileUrl:         resume.URL,
				MatchPercentage: analysis.MatchPercentage,
				Insights:        analysis.Insights,
			})
			id++
		}()
	}
	wg.Wait()
	return results, nil
}
