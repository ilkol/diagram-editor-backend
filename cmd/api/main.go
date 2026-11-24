package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilkol/diagram-editor-backend/internal/config"
	"github.com/ilkol/diagram-editor-backend/internal/handler"
	"github.com/ilkol/diagram-editor-backend/internal/service"
	"github.com/ilkol/diagram-editor-backend/internal/storage"
)

func main() {
	cfg := config.MustLoad()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	store, err := storage.New(ctx, cfg.DSN())
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer store.Close()

	migrationCtx, cancelMigration := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelMigration()

	if err := store.RunMigrations(migrationCtx); err != nil {
		log.Fatalf("failed to do migrations: %v", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("Alive!")); err != nil {
			log.Printf("failed to write response: %v", err)
			return
		}
	})

	boardService := service.New(store)
	boardHandler := handler.NewBoardHandler(boardService)

	router.Route("/boards", func(r chi.Router) {
		r.Post("/", boardHandler.Create)
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), router); err != nil {
		log.Fatalf("failed to init server: %v", err)
	}
}
