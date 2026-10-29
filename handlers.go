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

	t.Status = StatusNew
	t.CreatedAt = time.Now().UTC()

	res, err := s.db.Exec(`
	INSERT INTO todos (title, description, author, status, created_at) VALUES (?, ?, ?, ?, ?)
	`, t.Title, t.Description, t.Author, t.Status, t.CreatedAt.Format(time.RFC3339))

	if err != nil {
		return err
	}

	id, errInsert := res.LastInsertId()

	if errInsert != nil {
		return errInsert
	}

	t.ID = id

	return nil
}

// handleCreateTodo — хендлер POST /todos.
func (s *Store) handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is missing", http.StatusBadRequest)
		return
	}

	todo := Todo{
		Title:       req.Title,
		Description: req.Description,
		Author:      req.Author,
	}

	errFromDB := s.CreateTodo(&todo)

	if errFromDB != nil {
		http.Error(w, "Error from DB", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)

}

// TODO 5 (я потом): подключим это в main.go:
//   mux.HandleFunc("POST /todos", store.handleCreateTodo)
// и обсудим, почему handleCreateTodo — метод Store, а не голая функция.
