package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/WENDELLDELIMA/go_cep_weater/internal/cep"
	"github.com/WENDELLDELIMA/go_cep_weater/internal/weather"
)

type Server struct {
	cepSvc     *cep.Service
	weatherSvc *weather.Service
}

func New() *Server {
	cepBase := getenv("VIA_CEP_BASE", "https://viacep.com.br")
	weatherBase := getenv("WEATHERAPI_BASE", "https://api.weatherapi.com/v1")
	apiKey := os.Getenv("WEATHERAPI_KEY")
	httpClient := &http.Client{ Timeout: 10 * time.Second }

	return &Server{
		cepSvc:     cep.NewService(httpClient, cepBase),
		weatherSvc: weather.NewService(httpClient, weatherBase, apiKey),
	}
}

var reCEP = regexp.MustCompile(`^\d{5}-?\d{3}$`)

func (s *Server) HandleWeather(w http.ResponseWriter, r *http.Request) {
	cepParam := r.URL.Query().Get("cep")
	if !reCEP.MatchString(cepParam) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	// Normalize CEP by removing hyphen
	normalizedCEP := strings.ReplaceAll(cepParam, "-", "")

	loc, err := s.cepSvc.Lookup(r.Context(), normalizedCEP)
	if err != nil {
		if errors.Is(err, cep.ErrNotFound) {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		}
		http.Error(w, "via cep error", http.StatusBadGateway)
		return
	}

	// Build query for WeatherAPI
	q := fmt.Sprintf("%s,%s,Brazil", loc.City, loc.State)
	
	cTemp, err := s.weatherSvc.CurrentTempC(r.Context(), q)
	if err != nil {
		http.Error(w, fmt.Sprintf("weather provider error: %v", err), http.StatusBadGateway)
		return
	}

	// Convert per spec
	tempC := round1(cTemp)
	tempF := round1(cTemp*1.8 + 32.0) // F = C * 1.8 + 32
	tempK := round1(cTemp + 273.0)    // K = C + 273 (as per requirement)
	resp := map[string]float64{
		"temp_C": tempC,
		"temp_F": tempF,
		"temp_K": tempK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func round1(v float64) float64 {
	return math.Round(v*10) / 10
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
