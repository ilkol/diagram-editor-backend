package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilkol/diagram-editor-backend/internal/config"
)

func main() {
	cfg := config.MustLoad()

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("Alive!")); err != nil {
			log.Printf("failed to write respone: %v", err)
			return
		}
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), router); err != nil {
		log.Fatalf("failed to init server: %v", err)
	}
}
