package eggo

import "errors"

var (
	ErrComplaintAlreadyExists = errors.New("complaint already exists")
	ErrComplaintNotFound      = errors.New("complaint not found")
)
