package keywords

import (
	"context"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl/deps"
	"github.com/satont/tsuwari/libs/grpc/generated/api/keywords"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Keywords struct {
	*deps.Deps
}

func (c *Keywords) KeywordsGetAll(ctx context.Context, empty *emptypb.Empty) (*keywords.GetAllResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Keywords) KeywordsGetById(ctx context.Context, request *keywords.GetByIdRequest) (*keywords.Keyword, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Keywords) KeywordsCreate(ctx context.Context, request *keywords.CreateRequest) (*keywords.Keyword, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Keywords) KeywordsDelete(ctx context.Context, request *keywords.DeleteRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Keywords) KeywordsPut(ctx context.Context, request *keywords.PutRequest) (*keywords.Keyword, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Keywords) KeywordsPatch(ctx context.Context, request *keywords.PatchRequest) (*keywords.Keyword, error) {
	//TODO implement me
	panic("implement me")
}
