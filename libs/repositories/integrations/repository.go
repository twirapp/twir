package integrations

import (
    "context"
    
    "github.com/twirapp/twir/libs/repositories/integrations/model"
)

type Repository interface {
    GetByService(ctx context.Context, service model.Service) (model.Integration, error)
}