package gql

import (
	"context"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
	"github.com/twirapp/twir/apps/api-gql/internal/server/rate_limiter"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type rateLimitExtension struct {
	rateLimiter *rate_limiter.LeakyBucketRateLimiter
}

func (r *rateLimitExtension) ExtensionName() string {
	return "RateLimitExtension"
}

func (r *rateLimitExtension) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (r *rateLimitExtension) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	rc := graphql.GetOperationContext(ctx)

	// Skip rate limit if operation or selected schema field has @noRateLimit directive.
	if hasNoRateLimit(rc.Operation) {
		return next(ctx)
	}

	// Get client IP from gin context
	ginCtx, err := gincontext.GetGinContext(ctx)
	if err != nil {
		return next(ctx)
	}

	resp, err := r.rateLimiter.Use(ctx, &rate_limiter.LeakyOptions{
		KeyPrefix:       "api-gql:ratelimiter:graphql:" + ginCtx.ClientIP(),
		MaximumCapacity: 100,
		WindowSeconds:   60,
	}, 1)
	if err != nil {
		return next(ctx)
	}

	if !resp.Success {
		return &graphql.Response{
			Errors: gqlerror.List{{
				Message: "Rate limit exceeded",
				Extensions: map[string]interface{}{
					"code": "RATE_LIMITED",
				},
			}},
		}
	}

	return next(ctx)
}

func hasNoRateLimit(op *ast.OperationDefinition) bool {
	if op == nil {
		return false
	}
	if directivesHaveNoRateLimit(op.Directives) {
		return true
	}

	for _, sel := range op.SelectionSet {
		if selectionHasNoRateLimit(sel) {
			return true
		}
	}

	return false
}

func selectionHasNoRateLimit(sel ast.Selection) bool {
	switch selected := sel.(type) {
	case *ast.Field:
		if directivesHaveNoRateLimit(selected.Directives) {
			return true
		}

		if selected.Definition != nil && directivesHaveNoRateLimit(selected.Definition.Directives) {
			return true
		}

		for _, child := range selected.SelectionSet {
			if selectionHasNoRateLimit(child) {
				return true
			}
		}
	case *ast.InlineFragment:
		if directivesHaveNoRateLimit(selected.Directives) {
			return true
		}

		for _, child := range selected.SelectionSet {
			if selectionHasNoRateLimit(child) {
				return true
			}
		}
	}

	return false
}

func directivesHaveNoRateLimit(directives ast.DirectiveList) bool {
	for _, dir := range directives {
		if strings.EqualFold(dir.Name, "noRateLimit") {
			return true
		}
	}

	return false
}
