package controller

import "errors"

var (
	ErrEmptyContentOrTitle = errors.New("empty content or title")
	ErrTitleTooLong        = errors.New("title too long")
	ErrInvalidContent      = errors.New("invalid content")
)
