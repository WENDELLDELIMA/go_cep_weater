package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/WENDELLDELIMA/go_cep_weater/internal/server"
)

func main() {
	addr := ":" + getEnv("PORT", "8080")
	s := server.New()

	r := chi.NewRouter()
	r.Get("/weather", s.HandleWeather)

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
