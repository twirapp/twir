package directives

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator/v10"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var (
	validate *validator.Validate
)

func init() {
	validate = validator.New()
}

func (c *Directives) Validate(
	ctx context.Context,
	obj interface{},
	next graphql.Resolver,
	constraint string,
) (interface{}, error) {
	val, err := next(ctx)
	if err != nil {
		return nil, err
	}

	pathContext := graphql.GetPathContext(ctx)

	reqMap, ok := obj.(map[string]any)
	if ok && pathContext.Field != nil {
		v, valExists := reqMap[*pathContext.Field]
		if valExists && v == nil && strings.Contains(constraint, "omitempty") {
			return val, nil
		}
	}

	err = validate.Var(val, constraint)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			var fieldName string
			if pathContext.Field != nil {
				fieldName = *pathContext.Field
			}

			message := fmt.Sprintf(`Invalid field "%s"`, fieldName)

			extensionValidationErrors := make(map[string]any)
			for _, er := range validationErrors {
				extensionValidationErrors[er.Tag()] = er.Param()
			}

			return nil, &gqlerror.Error{
				Message:   message,
				Path:      pathContext.Path(),
				Locations: nil,
				Extensions: map[string]interface{}{
					"code":              "BAD_REQUEST",
					"validation_errors": extensionValidationErrors,
				},
			}
		}

		return nil, err
	}

	return val, nil
}
