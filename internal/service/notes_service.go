package service

import (
	"time"

	"github.com/avran02/kode/internal/models"
	"github.com/avran02/kode/internal/repository"
)

type NotesService interface {
	CreateNote(userID int, title, content string) (models.Note, error)
	GetNotes(userID int) ([]models.Note, error)
}

type notesService struct {
	repo repository.NotesRepository
}

func (s *notesService) CreateNote(userID int, title, content string) (models.Note, error) {
	note := models.Note{
		UserID:    userID,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
	}

	id, err := s.repo.CreateNote(note)
	if err != nil {
		return models.Note{}, err
	}

	note.ID = id
	return note, nil
}

func (s *notesService) GetNotes(userID int) ([]models.Note, error) {
	return s.repo.GetNotesByUserID(userID)
}

func NewNotesService(repo repository.NotesRepository) NotesService {
	return &notesService{repo: repo}
}
