package application

import (
	"context"

	"github.com/ricomonster/daedalus/config"
	"github.com/ricomonster/daedalus/daedalus"
)

type oracleApplication struct {
	config *config.Config
	gemini daedalus.LLMApplication
}

func NewOracleApplication(
	co *config.Config,
	gemini daedalus.LLMApplication,
) daedalus.OracleApplication {
	// Everytime this is instantiated, remove the temp values
	co.Set(daedalus.OracleConfigTempLLM, "")

	return &oracleApplication{co, gemini}
}

func (oa *oracleApplication) Prompt(ctx context.Context, prompt string) (string, error) {
	// Get the default llm
	llm := oa.config.GetString(daedalus.OracleConfigDefaultLLM)
	switch daedalus.LLM(llm) {
	case oa.gemini.Name():
	default:
		return oa.gemini.Prompt(ctx, prompt)
	}

	return "", daedalus.OracleErrUnsupportedLLM
}

func (oa *oracleApplication) SetLLMKey(
	ctx context.Context,
	llm daedalus.LLM,
	key string,
) error {
	// Get default llm?

	switch llm {
	case oa.gemini.Name():
	default:
		return oa.gemini.SetKey(ctx, key)
	}

	return nil
}
