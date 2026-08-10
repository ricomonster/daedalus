package gemini

import (
	"context"
	"errors"
	"sync"

	"google.golang.org/genai"

	"github.com/ricomonster/daedalus/config"
)

var (
	name      = "gemini"
	geminiKey = "GOOGLE_API_KEY"
	model     = "gemini-3.1-flash-lite-preview"
)

var ErrKeyNotProvided = errors.New("google key not found")

type Client struct {
	config *config.Config
	client *genai.Client
	mu     *sync.Mutex
}

func New(co *config.Config) *Client {
	return &Client{config: co}
}

func (c *Client) Name() string {
	return name
}

func (c *Client) SetKey(ctx context.Context, key string) error {
	c.config.Set(geminiKey, key)
	return c.config.Save()
}

func (c *Client) Prompt(ctx context.Context, prompt string) (string, error) {
	if c.client == nil {
		err := c.instantiate(ctx)
		if err != nil {
			return "", err
		}
	}

	result, err := c.client.Models.GenerateContent(
		ctx,
		model,
		genai.Text(prompt),
		nil,
	)
	// TODO: Add validation if model already reached the limit so we can switch to another
	// model in the list
	if err != nil {
		return "", err
	}

	return result.Text(), nil
}

func (c *Client) instantiate(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.client != nil {
		return nil
	}

	k := c.config.GetString(geminiKey)
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

	c.client = client
	return nil
}
