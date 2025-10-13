package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Service struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

func NewService(c *http.Client, base, key string) *Service {
	if c == nil {
		c = &http.Client{ Timeout: 10 * time.Second }
	}
	if base == "" {
		base = "https://api.weatherapi.com/v1"
	}
	return &Service{client: c, baseURL: base, apiKey: key}
}

type currentResp struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func (s *Service) CurrentTempC(ctx context.Context, query string) (float64, error) {
	url := fmt.Sprintf("%s/current.json?key=%s&q=%s", s.baseURL, s.apiKey, query)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := s.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weather status %d", resp.StatusCode)
	}

	var data currentResp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}
	return data.Current.TempC, nil
}
