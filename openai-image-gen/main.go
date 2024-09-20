package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

var prompt = flag.String("p", "", "Prompt")
var outfile = flag.String("o", "", "Output file name")

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
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		N:              1,
	}

	resp, err := client.CreateImage(context.Background(), req)
	if err != nil {
		fmt.Printf("Image creation error: %v\n", err)
		return
	}

	imgBytes, err := base64.RawStdEncoding.DecodeString(resp.Data[0].B64JSON)
	if err != nil {
		fmt.Printf("Base64 decode error: %v\n", err)
		return
	}

	r := bytes.NewBuffer(imgBytes)
	imgData, err := png.Decode(r)
	if err != nil {
		fmt.Printf("PNG decode error: %v\n", err)
		return
	}

	file, err := os.Create(fmt.Sprintf("%s.png", *outfile))
	if err != nil {
		fmt.Printf("File creation error: %v\n", err)
		return
	}
	defer file.Close()

	if err := png.Encode(file, imgData); err != nil {
		fmt.Printf("PNG encode error: %v\n", err)
		return
	}
}
