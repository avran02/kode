package mapper

import (
	"github.com/avran02/kode/internal/dto"
	"github.com/avran02/kode/internal/models"
)

func ToCreateNoteResponse(note models.Note) dto.CreateNoteResponse {
	return dto.CreateNoteResponse{
		ID:      note.ID,
		Title:   note.Title,
		Content: note.Content,
	}
}

func ToGetNotesResponse(notes []models.Note) dto.GetNotesResponse {
	notesResponse := make([]dto.NoteResponse, 0, len(notes))

	for _, note := range notes {
		notesResponse = append(notesResponse, dto.NoteResponse{
			ID:      note.ID,
			Title:   note.Title,
			Content: note.Content,
		})
	}

	return notesResponse
}
