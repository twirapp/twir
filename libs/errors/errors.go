package errors

import (
	"errors"
	"fmt"
)

// ErrorCode представляет тип ошибки для клиента
type ErrorCode string

const (
	// ErrorCodeValidation - ошибка валидации данных.
	ErrorCodeValidation ErrorCode = "VALIDATION_ERROR"
	// ErrorCodeNotFound - запрашиваемый ресурс не найден.
	ErrorCodeNotFound ErrorCode = "NOT_FOUND"
	// ErrorCodeConflict - конфликт данных (например, дубликат).
	ErrorCodeConflict ErrorCode = "CONFLICT"
	// ErrorCodeForbidden - доступ запрещен.
	ErrorCodeForbidden ErrorCode = "FORBIDDEN"
	// ErrorCodeInternal - внутренняя ошибка сервера.
	ErrorCodeInternal ErrorCode = "INTERNAL_ERROR"
	// ErrorCodeBadRequest - неверный запрос.
	ErrorCodeBadRequest ErrorCode = "BAD_REQUEST"
	// ErrorCodeUnauthorized - требуется аутентификация.
	ErrorCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	// ErrorCodeRateLimited - превышен лимит запросов.
	ErrorCodeRateLimited ErrorCode = "RATE_LIMITED"
)

// AppError представляет кастомную ошибку приложения.
type AppError struct {
	Code    ErrorCode
	Message string
	Details map[string]any
	Err     error
}

// Фабричные функции для создания типичных ошибок

// NewValidationError создает ошибку валидации.
func NewValidationError(message string, details map[string]any) *AppError {
	return New(ErrorCodeValidation, message).WithDetails(details)
}

// NewNotFoundError создает ошибку "не найдено".
func NewNotFoundError(message string) *AppError {
	return New(ErrorCodeNotFound, message)
}

// NewConflictError создает ошибку конфликта.
func NewConflictError(message string) *AppError {
	return New(ErrorCodeConflict, message)
}

// NewForbiddenError создает ошибку запрета доступа.
func NewForbiddenError(message string) *AppError {
	return New(ErrorCodeForbidden, message)
}

// NewInternalError создает ошибку сервера.
func NewInternalError(message string, err error) *AppError {
	return Wrap(ErrorCodeInternal, message, err)
}

// NewBadRequestError создает ошибку неверного запроса.
func NewBadRequestError(message string) *AppError {
	return New(ErrorCodeBadRequest, message)
}

// NewUnauthorizedError создает ошибку аутентификации.
func NewUnauthorizedError(message string) *AppError {
	return New(ErrorCodeUnauthorized, message)
}

// NewRateLimitedError создает ошибку превышения лимита.
func NewRateLimitedError(message string) *AppError {
	return New(ErrorCodeRateLimited, message)
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// New создает новую AppError.
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap оборачивает ошибку в AppError.
func Wrap(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WithDetails добавляет детали к ошибке.
func (e *AppError) WithDetails(details map[string]any) *AppError {
	e.Details = details
	return e
}

// IsAppError проверяет, является ли ошибка AppError.
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// AsAppError приводит ошибку к типу AppError.
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
