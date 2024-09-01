package controller

import (
	"encoding/json"
	"log/slog"
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
	noteService service.NotesService
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

	slog.Info("create note", "title", req.Title, "content", req.Content, "user_id", userID)

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

func newNoteController(s service.NotesService) NoteController {
	return &noteController{
		noteService: s,
	}
}
