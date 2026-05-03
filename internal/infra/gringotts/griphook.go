package gringotts

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/clodoaldomarques/ledger-config/configs"
	"github.com/ollama/ollama/api"
)

type Griphook struct {
	URL   string
	Model string
}

func New() *Griphook {
	url := configs.New().GriphookIAUrl
	return &Griphook{
		URL:   url,
		Model: "griphook:latest", // ou "validador-payload"
	}
}

func (g Griphook) ValidateJSONConfig(ctx context.Context, config string) (GriphookResponse, error) {
	// Parse da URL base do serviço Ollama (ex: "http://ollama-service:11434")
	baseURL, err := url.Parse(g.URL)
	if err != nil {
		return GriphookResponse{}, err
	}

	// Cria um cliente HTTP simples
	httpClient := &http.Client{}

	// Cria o cliente do Ollama com a URL customizada
	client := api.NewClient(baseURL, httpClient)

	req := &api.GenerateRequest{
		Model:  g.Model,
		Prompt: config,
		Stream: new(bool), // false para resposta única (opcional)
	}
	*req.Stream = false

	var fullResponse string
	respFunc := func(resp api.GenerateResponse) error {
		fullResponse += resp.Response
		return nil
	}

	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		return GriphookResponse{}, err
	}

	var gr GriphookResponse
	if err := json.Unmarshal([]byte(fullResponse), &gr); err != nil {
		return GriphookResponse{}, err
	}

	return gr, nil
}
