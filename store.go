package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func NewStore(path string) (*Store, error) {
	// открыть базу
	db, err := sql.Open("sqlite", path)

	// REVIEW (blocker): пропущена проверка err от sql.Open и db.Ping().
	// sql.Open не подключается, он только валидирует DSN. Реальная проверка — Ping():
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	// И общее правило: в Go ошибка проверяется СРАЗУ после каждого вызова, который её возвращает.

	// создать стор с указателем на базу
	store := new(Store{db: db})
	return store, nil
}

func (s *Store) initSchema() error {
	// создать таблицу
	_, err := s.db.Exec(`
	CREATE TABLE todos
	(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT DEFAULT '',
		status TEXT NOT NULL DEFAULT 'new',
		author TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT datetime(),
		completed_at TEXT DEFAULT NULL
	)
	`)

	// REVIEW (по схеме, всё верно 👍):
	// - title NOT NULL, completed_at NULL — ровно то, о чём договаривались.
	// - description DEFAULT '' — ок, пустая строка как честное "нет описания".
	//
	// REVIEW (nit): created_at через datetime() отдаст текст в формате "2026-10-28 20:15:00",
	// а фронтенд и стандарт API ждут RFC3339 ("2026-10-28T20:15:00Z"). Разберём при чтении из базы —
	// тогда и решим, парсить или хранить сразу в RFC3339.
	//
	// REVIEW (nit): имя таблицы "Todos" — в SQL обычно snake_case единственное число: "todos".
	// Не ошибка, но конвенция почти всех кодовых баз.

	return err
}
