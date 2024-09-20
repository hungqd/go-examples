package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

var prompt = flag.String("p", "", "Prompt")

func main() {
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Load .env error", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Environment variable OPENAI_API_KEY is required")
	}

	client := openai.NewClient(apiKey)

	req := openai.ImageRequest{
		Prompt:         *prompt,
		Size:           openai.CreateImageSize256x256,
		ResponseFormat: openai.CreateImageResponseFormatURL,
		N:              1,
	}

	resp, err := client.CreateImage(context.Background(), req)
	if err != nil {
		fmt.Printf("Image creation error: %v\n", err)
		return
	}

	url := resp.Data[0].URL
	fmt.Println(url)
}
