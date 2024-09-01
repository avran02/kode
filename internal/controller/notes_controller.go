package controller

import (
	"encoding/json"
	"net/http"

	contextkeys "github.com/avran02/kode/internal/context_keys"
	"github.com/avran02/kode/internal/dto"
	"github.com/avran02/kode/internal/mapper"
	"github.com/avran02/kode/internal/service"
)

type NoteController interface {
	CreateNote(w http.ResponseWriter, r *http.Request)
	GetNotes(w http.ResponseWriter, r *http.Request)
}

type noteController struct {
	noteService          service.NotesService
	yandexSpellerService service.YandexSpellerService
}

func (c *noteController) CreateNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(contextkeys.UserID).(int)
	if !ok {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	var req dto.CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.validateNote(w, req); err != nil {
		// Ошибка будет отправлена на клиент в методе validateNote
		return
	}

	noteModel, err := c.noteService.CreateNote(userID, req.Title, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(mapper.ToCreateNoteResponse(noteModel)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *noteController) GetNotes(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(contextkeys.UserID).(int)
	if !ok {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	notes, err := c.noteService.GetNotes(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(mapper.ToGetNotesResponse(notes)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *noteController) validateNote(w http.ResponseWriter, req dto.CreateNoteRequest) error {
	if req.Title == "" || req.Content == "" {
		http.Error(w, "empty title or content", http.StatusBadRequest)
		return ErrEmptyContentOrTitle
	}

	if len(req.Title) > 255 {
		http.Error(w, "title is too long", http.StatusBadRequest)
		return ErrTitleTooLong
	}

	spellErrors, err := c.yandexSpellerService.CheckText(req.Title + " " + req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	if len(spellErrors) > 0 {
		spellErrorsJSON, err := json.Marshal(spellErrors)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}

		http.Error(w, string(spellErrorsJSON), http.StatusBadRequest)
		return ErrInvalidContent
	}

	return nil
}

func newNoteController(s service.NotesService, yss service.YandexSpellerService) NoteController {
	return &noteController{
		noteService:          s,
		yandexSpellerService: yss,
	}
}
