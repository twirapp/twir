package kappagen

import (
	"reflect"
	"testing"

	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/libs/bus-core/api"
)

func TestBuildTriggerEmotes(t *testing.T) {
	got := BuildTriggerEmotes([]*types.ParseContextEmote{
		{
			ID:  "4148074",
			URL: "https://files.kick.com/emotes/4148074/fullsize",
			Positions: []*types.ParseContextEmotePosition{
				{Start: 10, End: 15},
				{Start: 17, End: 22},
			},
		},
		{
			ID:  "4148074",
			URL: "https://files.kick.com/emotes/4148074/fullsize",
			Positions: []*types.ParseContextEmotePosition{
				{Start: 24, End: 29},
			},
		},
	})

	want := []api.TriggerKappagenEmote{
		{Id: "4148074", Url: "https://files.kick.com/emotes/4148074/fullsize", Positions: []string{"10-15", "17-22", "24-29"}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("BuildTriggerEmotes() = %#v, want %#v", got, want)
	}
}
