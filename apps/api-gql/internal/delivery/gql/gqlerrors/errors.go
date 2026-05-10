package gqlerrors

import (
	"errors"

	apperrors "github.com/twirapp/twir/libs/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// InternalCauseKey is a private extensions key used to pass the original error to the
// ErrorPresenter for logging. It is stripped from the response before sending to the client.
const InternalCauseKey = "_internal_cause"

// ToGQLError преобразует ошибку приложения в GraphQL ошибку
func ToGQLError(err error) error {
	if err == nil {
		return nil
	}

	appErr, ok := apperrors.AsAppError(err)
	if !ok {
		// Если это не AppError, возвращаем как внутреннюю ошибку.
		// Оригинальная ошибка сохраняется под приватным ключом для логирования в ErrorPresenter.
		return &gqlerror.Error{
			Message: "Internal server error",
			Extensions: map[string]interface{}{
				"code":           string(apperrors.ErrorCodeInternal),
				InternalCauseKey: err,
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
