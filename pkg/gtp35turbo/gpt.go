package gtp35turbo

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"io"
)

func Response(client *openai.Client, prompt string) (string, error) {
	// open ai
	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo0301,
		MaxTokens: 2000,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Stream: true,
	}
	stream, err := client.CreateChatCompletionStream(ctx, req)

	if err != nil {
		return "", err
	}

	return getStream(stream), nil
}

func getStream(stream *openai.ChatCompletionStream) string {
	var result string
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			log.Info().Msg("\nStream finished")
			stream.Close()
			break
		}

		if err != nil {
			log.Warn().Msgf("\nStream error: %v\n", err)
			stream.Close()
			break
		}
		log.Printf(response.Choices[0].Delta.Content)
		result += response.Choices[0].Delta.Content
	}
	return result
}
