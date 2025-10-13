package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/WENDELLDELIMA/go_cep_weater/internal/server"
)

func TestWeather_Success(t *testing.T) {
	// Mock ViaCEP
	via := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"localidade":"São Paulo","uf":"SP"}`))
	}))
	defer via.Close()

	// Mock WeatherAPI
	weather := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"current":{"temp_c": 25.0}}`))
	}))
	defer weather.Close()

	os.Setenv("VIA_CEP_BASE", via.URL)
	os.Setenv("WEATHERAPI_BASE", weather.URL)
	os.Setenv("WEATHERAPI_KEY", "dummy")

	s := server.New()
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=01001000", nil)
	w := httptest.NewRecorder()

	s.HandleWeather(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var body map[string]float64
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body["temp_C"] != 25.0 {
		t.Fatalf("temp_C expected 25.0, got %v", body["temp_C"])
	}
	if body["temp_F"] != 77.0 {
		t.Fatalf("temp_F expected 77.0, got %v", body["temp_F"])
	}
	if body["temp_K"] != 298.0 {
		t.Fatalf("temp_K expected 298.0, got %v", body["temp_K"])
	}
}

func TestWeather_InvalidCEP(t *testing.T) {
	s := server.New()
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=ABC", nil)
	w := httptest.NewRecorder()

	s.HandleWeather(w, req)
	if w.Result().StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Result().StatusCode)
	}
}

func TestWeather_CEPNotFound(t *testing.T) {
	// Mock ViaCEP retorna erro true
	via := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"erro": true}`))
	}))
	defer via.Close()

	os.Setenv("VIA_CEP_BASE", via.URL)
	os.Setenv("WEATHERAPI_BASE", "http://invalid") // não deve chamar
	os.Setenv("WEATHERAPI_KEY", "dummy")

	s := server.New()
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=01001000", nil)
	w := httptest.NewRecorder()

	s.HandleWeather(w, req)
	if w.Result().StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Result().StatusCode)
	}
}
