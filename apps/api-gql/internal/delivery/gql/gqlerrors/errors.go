package gqlerrors

import (
	"errors"

	apperrors "github.com/twirapp/twir/libs/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// ToGQLError преобразует ошибку приложения в GraphQL ошибку
func ToGQLError(err error) error {
	if err == nil {
		return nil
	}

	appErr, ok := apperrors.AsAppError(err)
	if !ok {
		// Если это не AppError, возвращаем как внутреннюю ошибку
		return &gqlerror.Error{
			Message: "Internal server error",
			Extensions: map[string]interface{}{
				"code": string(apperrors.ErrorCodeInternal),
			},
		}
	}

	gqlErr := &gqlerror.Error{
		Message: appErr.Message,
		Extensions: map[string]interface{}{
			"code": string(appErr.Code),
		},
	}

	// Добавляем детали, если они есть
	if appErr.Details != nil && len(appErr.Details) > 0 {
		gqlErr.Extensions["details"] = appErr.Details
	}

	// Для внутренних ошибок не показываем оригинальную ошибку клиенту
	if appErr.Code != apperrors.ErrorCodeInternal && appErr.Err != nil {
		gqlErr.Extensions["original_error"] = appErr.Err.Error()
	}

	return gqlErr
}

// HandleError обрабатывает ошибку и преобразует её в GraphQL формат
// Использовать в резолверах: return nil, gqlerrors.HandleError(err)
func HandleError(err error) error {
	if err == nil {
		return nil
	}

	// Если это уже gqlerror, просто возвращаем
	var gqlErr *gqlerror.Error
	if errors.As(err, &gqlErr) {
		return err
	}

	// Преобразуем AppError в GraphQL ошибку
	return ToGQLError(err)
}
