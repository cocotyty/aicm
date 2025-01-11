package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cocotyty/aicm/internal/config"
	"github.com/cocotyty/aicm/internal/git"
)

const (
	systemPrompt = `You are a helpful AI assistant that helps generate commit messages.
For code changes, first describe what changed in each file, then generate a concise commit message.
The commit message should follow the conventional commit format: "type: description"`
)

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model    string          `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func GenerateCommitMessage(cfg *config.Config, changes []git.FileChange) (string, error) {
	// 第一步：生成代码变更描述
	descriptions, err := generateFileDescriptions(cfg, changes)
	if err != nil {
		return "", err
	}

	// 第二步：生成commit message
	return generateFinalCommitMessage(cfg, descriptions)
}

func generateFileDescriptions(cfg *config.Config, changes []git.FileChange) (map[string]string, error) {
	log.Println("Starting to generate file descriptions")
	var prompt string
	for _, change := range changes {
		prompt += fmt.Sprintf("File: %s\nDiff:\n%s\n\n", change.FileName, change.Diff)
	}
	prompt += "Please describe the changes in each file."
	log.Printf("Generated prompt for %d files", len(changes))

	messages := []OpenAIMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: prompt},
	}

	response, err := callOpenAI(cfg, messages)
	if err != nil {
		log.Printf("Error generating file descriptions: %v", err)
		return nil, err
	}

	descriptions := make(map[string]string)
	for _, change := range changes {
		descriptions[change.FileName] = response
	}
	log.Println("Successfully generated file descriptions")
	return descriptions, nil
}

func generateFinalCommitMessage(cfg *config.Config, descriptions map[string]string) (string, error) {
	log.Println("Starting to generate final commit message")
	var prompt string
	for file, desc := range descriptions {
		prompt += fmt.Sprintf("%s: %s\n", file, desc)
	}
	prompt += "Please generate a commit message following the conventional commit format."
	log.Printf("Generated commit message prompt with %d file descriptions", len(descriptions))

	messages := []OpenAIMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: prompt},
	}

	msg, err := callOpenAI(cfg, messages)
	if err != nil {
		log.Printf("Error generating commit message: %v", err)
		return "", err
	}
	log.Printf("Successfully generated commit message: %s", msg)
	return msg, nil
}

func callOpenAI(cfg *config.Config, messages []OpenAIMessage) (string, error) {
	log.Println("Making OpenAI API call")
	request := OpenAIRequest{
		Model:    cfg.LLMModel,
		Messages: messages,
	}
	for _, msg := range messages {
		log.Printf("%s: %s \n", msg.Role, msg.Content)
	}
	requestBody, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error marshaling request: %v", err)
		return "", err
	}

	log.Printf("Sending request to %s with model %s", cfg.LLMAPIURL, cfg.LLMModel)
	req, err := http.NewRequest("POST", cfg.LLMAPIURL+"/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.LLMAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making HTTP request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("API returned non-200 status: %d", resp.StatusCode)
		return "", fmt.Errorf("OpenAI API returned status: %d", resp.StatusCode)
	}

	var response OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("Error decoding response: %v", err)
		return "", err
	}

	if len(response.Choices) == 0 {
		log.Println("No choices in API response")
		return "", fmt.Errorf("no response from OpenAI")
	}

	log.Println("Successfully received response from OpenAI")
	return response.Choices[0].Message.Content, nil
}
