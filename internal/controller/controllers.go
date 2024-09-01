package controller

import "github.com/avran02/kode/internal/service"

type Controller interface {
	AuthController
	NoteController
}

type controller struct {
	AuthController
	NoteController
}

func New(as service.AuthService, ns service.NotesService, yss service.YandexSpellerService) Controller {
	return &controller{
		AuthController: newAuthController(as),
		NoteController: newNoteController(ns, yss),
	}
}
