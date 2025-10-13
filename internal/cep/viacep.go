package cep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrNotFound = errors.New("cep not found")

type Service struct {
	client  *http.Client
	baseURL string
}

func NewService(c *http.Client, base string) *Service {
	if c == nil {
		c = &http.Client{ Timeout: 10 * time.Second }
	}
	if base == "" {
		base = "https://viacep.com.br"
	}
	return &Service{client: c, baseURL: base}
}

type viaCepResp struct {
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
	Erro       bool   `json:"erro"`
}

type Location struct {
	City  string
	State string
}

func (s *Service) Lookup(ctx context.Context, cep string) (Location, error) {
	url := fmt.Sprintf("%s/ws/%s/json/", s.baseURL, cep)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := s.client.Do(req)
	if err != nil {
		return Location{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return Location{}, ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return Location{}, fmt.Errorf("viacep status %d", resp.StatusCode)
	}

	var data viaCepResp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return Location{}, err
	}
	if data.Erro || data.Localidade == "" || data.UF == "" {
		return Location{}, ErrNotFound
	}
	return Location{City: data.Localidade, State: data.UF}, nil
}
