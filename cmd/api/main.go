package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
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

	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatalf("failed to init server: %v", err)
	}
}
