package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Load .env error", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Environment variable OPENAI_API_KEY is required")
	}

	client := openai.NewClient(apiKey)
	messages := make([]openai.ChatCompletionMessage, 0)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("--> ")
		text, _ := reader.ReadString('\n')
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: text,
		})

		resp, err := client.CreateChatCompletion(
			context.Background(), openai.ChatCompletionRequest{
				Model:    openai.GPT3Dot5Turbo,
				Messages: messages,
			},
		)

		if err != nil {
			fmt.Printf("Chat completion error: %v\n", err)
			continue
		}

		content := resp.Choices[0].Message.Content
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: content,
		})
		fmt.Println(content)
	}
}
