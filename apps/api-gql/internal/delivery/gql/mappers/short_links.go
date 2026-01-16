package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
)

func ShortLinkViewToGQL(v shortenedurls.ViewOutput) gqlmodel.ShortLinkView {
	return gqlmodel.ShortLinkView{
		ShortLinkID: v.ShortLinkID,
		UserID:      v.UserID,
		Country:     v.Country,
		City:        v.City,
		CreatedAt:   v.CreatedAt,
	}
}

func ShortLinkViewUpdateToGQL(u shortenedurls.ShortLinkViewUpdate) gqlmodel.ShortLinkViewUpdate {
	var lastView *gqlmodel.ShortLinkView
	if u.LastView != nil {
		view := gqlmodel.ShortLinkView{
			ShortLinkID: u.LastView.ShortLinkID,
			UserID:      u.LastView.UserID,
			Country:     u.LastView.Country,
			City:        u.LastView.City,
			CreatedAt:   u.LastView.CreatedAt,
		}
		lastView = &view
	}

	return gqlmodel.ShortLinkViewUpdate{
		ShortLinkID: u.ShortLinkID,
		TotalViews:  u.TotalViews,
		LastView:    lastView,
	}
}
