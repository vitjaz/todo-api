package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// Todo — то, что ходит по HTTP (JSON).
// Упоминания Store здесь быть не должно — это чистая модель данных.
type Todo struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	Author      string     `json:"author"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

// createTodoRequest — что присылает клиент на POST /todos.
// Отдельный тип (не Todo!), потому что клиент НЕ может задавать id/created_at/completed_at.
type createTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
}

// TODO 1 (ты): допиши Store.CreateTodo(t *Todo) error в store.go:
//   - INSERT с плейсхолдерами (?), НИКОГДА не склеивай SQL со строками пользователя — SQL-инъекции
//   - после INSERT достань сгенерированный id: result.LastInsertId()
//   - поставь t.CreatedAt = time.Now().UTC() и t.Status = StatusNew перед INSERT
//     (или оставь это серверу через DEFAULT — но тогда верни id и перечитай)
func (s *Store) CreateTodo(t *Todo) error {
	return nil // заглушка
}

// handleCreateTodo — хендлер POST /todos.
func (s *Store) handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest

	// TODO 2 (ты): декодируй тело запроса в req через json.NewDecoder(r.Body).Decode(&req)
	//   - если ошибка декодирования → http.Error(w, "...", http.StatusBadRequest) и return
	//   - проверь req.Title != "" — иначе тоже 400 (NOT NULL из схемы здесь, на входе)

	// TODO 3 (ты): собери Todo из req, вызови s.CreateTodo(&todo)
	//   - ошибка базы → 500, НЕ свети текст ошибки базы клиенту (http.Error с общим текстом)

	// TODO 4 (ты): верни созданный объект:
	//   w.Header().Set("Content-Type", "application/json")
	//   w.WriteHeader(http.StatusCreated) // 201, а не 200 — ресурс создан
	//   json.NewEncoder(w).Encode(todo)
}

// TODO 5 (я потом): подключим это в main.go:
//   mux.HandleFunc("POST /todos", store.handleCreateTodo)
// и обсудим, почему handleCreateTodo — метод Store, а не голая функция.
