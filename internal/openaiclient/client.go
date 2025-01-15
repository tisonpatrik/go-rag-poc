package openaiclient

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// Service interface
type Service interface {
	ChatCompletion(ctx context.Context, params openai.ChatCompletionNewParams) (*openai.ChatCompletion, error)
}

// service struct
type service struct {
	client *openai.Client
}

// Implementation of ChatCompletion method
func (s *service) ChatCompletion(ctx context.Context, params openai.ChatCompletionNewParams) (*openai.ChatCompletion, error) {
	return s.client.Chat.Completions.New(ctx, params)
}

var (
	apiKey         = os.Getenv("OPENAI_API_KEY")
	clientOnce     sync.Once
	clientInstance *service
)

func New() Service {
	clientOnce.Do(func() {
		if apiKey == "" {
			log.Fatalf("OPENAI_API_KEY environment variable is not set")
		}

		client := openai.NewClient(
			option.WithAPIKey(apiKey),
		)

		clientInstance = &service{
			client: client,
		}

		log.Printf("OpenAI client initialized successfully")
	})

	return clientInstance
}
