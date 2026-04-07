package daedalus

import (
	"context"
	"errors"
)

var (
	OracleConfigTempLLM    = "TEMP_LLM"
	OracleConfigDefaultLLM = "LLM"
)

var OracleErrUnsupportedLLM = errors.New("unsupported llm")

type (
	LLM string

	LLMApplication interface {
		Name() LLM

		Prompt(ctx context.Context, prompt string) (string, error)

		SetKey(ctx context.Context, key string) error
	}

	OracleApplication interface {
		Prompt(ctx context.Context, prompt string) (string, error)

		SetLLMKey(ctx context.Context, llm LLM, key string) error
	}
)
