package gemini

import (
	"context"
	"errors"
	"sync"

	"google.golang.org/genai"

	"github.com/ricomonster/daedalus/config"
	"github.com/ricomonster/daedalus/daedalus"
)

var (
	name      daedalus.LLM = "gemini"
	geminiKey              = "GOOGLE_API_KEY"
	model                  = "gemini-3.1-flash-lite-preview"
)

var ErrKeyNotProvided = errors.New("google key not found")

type gemini struct {
	config *config.Config
	client *genai.Client
	mu     sync.Mutex
}

func New(co *config.Config) daedalus.LLMApplication {
	return &gemini{config: co}
}

func (g *gemini) Name() daedalus.LLM {
	return name
}

func (g *gemini) SetKey(ctx context.Context, key string) error {
	g.config.Set(geminiKey, key)
	return g.config.Save()
}

func (g *gemini) Prompt(ctx context.Context, prompt string) (string, error) {
	if g.client == nil {
		err := g.instantiate(ctx)
		if err != nil {
			return "", err
		}
	}

	result, err := g.client.Models.GenerateContent(
		ctx,
		model,
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", err
	}

	return result.Text(), nil
}

func (g *gemini) instantiate(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.client != nil {
		return nil
	}

	k := g.config.GetString(geminiKey)
	if k == "" {
		return ErrKeyNotProvided
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  k,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return err
	}

	g.client = client
	return nil
}
