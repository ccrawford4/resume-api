package openai

import (
	context "context"
	"encoding/json"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type ResumeAnalysis struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	UploadDate      string                 `json:"uploadDate"`
	FileUrl         string                 `json:"fileUrl"`
	MatchPercentage int                    `json:"matchPercentage"`
	Insights        map[string]interface{} `json:"insights"`
}

func AnalyzeResume(jobDescription string, resumes map[string]string) ([]ResumeAnalysis, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}
	client := openai.NewClient(apiKey)
	ctx := context.Background()

	results := []ResumeAnalysis{}
	id := 1
	fmt.Printf("resumes %s\n", resumes)
	for name, content := range resumes {
		fmt.Printf("name: %s, content: %s\n", name, content)
		prompt := "Given the following job description: '" + jobDescription + "' and the following resume: '" + content + "', analyze the resume for match percentage, strengths, improvements, and missing skills. Respond in JSON with fields: matchPercentage, insights (with strengths, improvements, missingSkills as arrays of strings). NO MARKDOWN JUST PLAIN JSON"
		resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: "You are a helpful assistant that analyzes resumes for job fit."},
				{Role: openai.ChatMessageRoleUser, Content: prompt},
			},
		})
		fmt.Printf("resp: %v\n", resp.Choices)
		if err != nil || len(resp.Choices) == 0 {
			continue
		}
		// Parse the response JSON
		var analysis struct {
			MatchPercentage int            `json:"matchPercentage"`
			Insights        map[string]any `json:"insights"`
		}
		err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &analysis)
		if err != nil {
			continue
		}
		results = append(results, ResumeAnalysis{
			ID:              fmt.Sprintf("%d", id),
			Name:            name,
			UploadDate:      "2024-01-15", // Placeholder
			FileUrl:         "/api/resume-analysis/resumes/" + fmt.Sprintf("%d", id) + "/preview",
			MatchPercentage: analysis.MatchPercentage,
			Insights:        analysis.Insights,
		})
		id++
	}
	return results, nil
}
