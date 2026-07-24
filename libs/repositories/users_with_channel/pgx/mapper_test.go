package pgx

import (
	"testing"

	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
)

func TestMapBindingToEntityMapsModelNilToEntityNil(t *testing.T) {
	t.Parallel()

	if entity := mapBindingToEntity(channelplatformsmodel.Nil); !entity.IsNil() {
		t.Fatal("model.Nil must map to entity.Nil")
	}
}
