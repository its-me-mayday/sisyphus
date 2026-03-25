package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/its-me-mayday/sisyphus/bff/internal/client"
	"github.com/its-me-mayday/sisyphus/bff/internal/handler"
)

func main() {
	port := getEnv("HTTP_PORT", "8090")
	tasksAddr := getEnv("TASKS_ADDR", "localhost:50052")
	aclAddr := getEnv("ACL_ADDR", "localhost:50053")

	// gRPC clients
	tasksClient, tasksConn, err := client.NewTasksClient(tasksAddr)
	if err != nil {
		log.Fatalf("tasks client: %v", err)
	}
	defer tasksConn.Close()

	aclClient, aclConn, err := client.NewAclClient(aclAddr)
	if err != nil {
		log.Fatalf("acl client: %v", err)
	}
	defer aclConn.Close()

	// Handler
	h := handler.New(tasksClient, aclClient)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}))

	// Health
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Lists
		r.Get("/lists", h.ListLists)
		r.Post("/lists", h.CreateList)
		r.Get("/lists/{id}", h.GetList)
		r.Put("/lists/{id}", h.UpdateList)
		r.Delete("/lists/{id}", h.DeleteList)

		// Todos
		r.Post("/lists/{listId}/todos", h.CreateTodo)
		r.Put("/todos/{id}", h.UpdateTodo)
		r.Delete("/todos/{id}", h.DeleteTodo)
	})

	log.Printf("✓ bff listening on :%s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
