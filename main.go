package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	store, err := NewStore("todos.db")

	if err != nil {
		log.Fatalf("Create store err: %v", err)
	}

	defer store.Close()

	initErr := store.initSchema()

	if initErr != nil {
		log.Fatalf("Init schema failed: %v", initErr)
	}

	mux.HandleFunc("POST /todos", store.handleCreateTodo)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
