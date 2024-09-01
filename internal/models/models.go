package models

import "time"

type Note struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	CreatedAt time.Time
}

type User struct {
	ID       int
	Username string
	Password string
}
