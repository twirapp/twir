package resolvers

import (
	"context"

	"github.com/twirapp/twir/libs/logger"
)

func (r *Resolver) resolveShortLinkDomain(ctx context.Context, shortLinkID string) *string {
	user, err := r.deps.Sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil
	}

	customDomain, err := r.deps.ShortLinksCustomDomainsService.GetByUserID(ctx, user.ID)
	if err != nil {
		r.deps.Logger.Warn("failed to get short links custom domain", logger.Error(err))
		return nil
	}

	if customDomain.IsNil() || !customDomain.Verified {
		return nil
	}

	domain := customDomain.Domain
	link, err := r.deps.ShortenedUrlsService.GetByShortID(ctx, &domain, shortLinkID)
	if err != nil {
		r.deps.Logger.Warn("failed to resolve short link domain", logger.Error(err))
		return nil
	}

	if link.IsNil() {
		return nil
	}

	return &domain
}
