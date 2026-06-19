package domain

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found in DB")
	ErrUserAlreadyExist = errors.New("user already exist in DB")
)
