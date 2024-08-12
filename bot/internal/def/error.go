package def

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidDocument = errors.New("document is nil")
	ErrChoicesIsEmpty  = errors.New("choices is empty")
	ErrCallbackData    = errors.New("callback data is invalid")
	ErrContextData     = errors.New("context data is invalid")
	ErrInvalidType     = errors.New("invalid type")
)
