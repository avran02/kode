package repository

import (
	"database/sql"
	"fmt"

	"github.com/avran02/kode/internal/models"
	_ "github.com/lib/pq"
)

type NotesRepository interface {
	CreateNote(note models.Note) (int, error)
	GetNotesByUserID(userID int) ([]models.Note, error)
}

type notesRepository struct {
	db *sql.DB
}

func (r *notesRepository) CreateNote(note models.Note) (int, error) {
	query := "INSERT INTO notes (user_id, title, content, created_at) VALUES ($1, $2, $3, $4) RETURNING id"
	var id int
	err := r.db.QueryRow(query, note.UserID, note.Title, note.Content, note.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *notesRepository) GetNotesByUserID(userID int) ([]models.Note, error) {
	query := "SELECT id, user_id, title, content, created_at FROM notes WHERE user_id = $1"
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notes: %w", err)
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get notes: %w", err)
	}

	return notes, nil
}

func NewNotesRepository(db *sql.DB) NotesRepository {
	return &notesRepository{db: db}
}
