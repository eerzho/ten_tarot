package failure

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrValidation       = errors.New("validation error")
	ErrInvalidDocument  = errors.New("document is nil")
	ErrAnswerIsRequired = errors.New("answer is required")
	ErrChoicesIsEmpty   = errors.New("choices is empty")
	ErrCallbackData     = errors.New("callback data is invalid")
)
