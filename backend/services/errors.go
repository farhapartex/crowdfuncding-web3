package services

type ValidationError struct{ Message string }

func (e *ValidationError) Error() string { return e.Message }

func NewValidationError(message string) error {
	return &ValidationError{Message: message}
}

type NotFoundError struct{ Message string }

func (e *NotFoundError) Error() string { return e.Message }

func NewNotFoundError(message string) error {
	return &NotFoundError{Message: message}
}

type ForbiddenError struct{ Message string }

func (e *ForbiddenError) Error() string { return e.Message }

func NewForbiddenError(message string) error {
	return &ForbiddenError{Message: message}
}

type UnauthorizedError struct{ Message string }

func (e *UnauthorizedError) Error() string { return e.Message }

func NewUnauthorizedError(message string) error {
	return &UnauthorizedError{Message: message}
}

type ConflictError struct{ Message string }

func (e *ConflictError) Error() string { return e.Message }

func NewConflictError(message string) error {
	return &ConflictError{Message: message}
}

type UnavailableError struct{ Message string }

func (e *UnavailableError) Error() string { return e.Message }

func NewUnavailableError(message string) error {
	return &UnavailableError{Message: message}
}
