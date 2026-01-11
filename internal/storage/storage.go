package storage

import "errors"

var (
	ErrUrlAlreadyExists = errors.New("url already exists")
	ErrUrlNotFound      = errors.New("url not found")
)
