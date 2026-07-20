package services

import (
	"strings"
	"testing"

	"github.com/lib/pq"
	model "github.com/twirapp/twir/libs/gomodels"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestHasCommandConflict(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		defaultName    string
		defaultAliases []string
		existing       []model.ChannelsCommands
		want           bool
	}{
		{
			name:        "no conflict",
			defaultName: "mmr",
			existing: []model.ChannelsCommands{
				{Name: "uptime", Aliases: pq.StringArray{"up"}},
			},
		},
		{
			name:           "default alias collides with existing name",
			defaultName:    "dota",
			defaultAliases: []string{"MMR"},
			existing: []model.ChannelsCommands{
				{Name: "mmr"},
			},
			want: true,
		},
		{
			name:           "default alias collides with existing alias",
			defaultName:    "dota",
			defaultAliases: []string{"winrate"},
			existing: []model.ChannelsCommands{
				{Name: "stats", Aliases: pq.StringArray{"WINRATE"}},
			},
			want: true,
		},
		{
			name:           "unrelated default aliases do not collide",
			defaultName:    "dota",
			defaultAliases: []string{"mmr", "rank"},
			existing: []model.ChannelsCommands{
				{Name: "uptime", Aliases: pq.StringArray{"up"}},
			},
		},
		{
			name:        "exact name conflict",
			defaultName: "mmr",
			existing: []model.ChannelsCommands{
				{Name: "mmr"},
			},
			want: true,
		},
		{
			name:        "case insensitive name conflict",
			defaultName: "mmr",
			existing: []model.ChannelsCommands{
				{Name: "MMR"},
			},
			want: true,
		},
		{
			name:        "alias conflict",
			defaultName: "mmr",
			existing: []model.ChannelsCommands{
				{Name: "rank", Aliases: pq.StringArray{"mmr"}},
			},
			want: true,
		},
		{
			name:        "case insensitive alias conflict",
			defaultName: "mmr",
			existing: []model.ChannelsCommands{
				{Name: "rank", Aliases: pq.StringArray{"MMR"}},
			},
			want: true,
		},
		{
			name:        "unrelated alias does not conflict",
			defaultName: "mmr",
			existing: []model.ChannelsCommands{
				{Name: "rank", Aliases: pq.StringArray{"rating"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := hasCommandConflict(tt.defaultName, tt.defaultAliases, tt.existing); got != tt.want {
				t.Fatalf("hasCommandConflict(%q, %v, %v) = %v, want %v", tt.defaultName, tt.defaultAliases, tt.existing, got, tt.want)
			}
		})
	}
}

func TestDefaultCommandsOnConflictUsesChannelNameDoNothing(t *testing.T) {
	db, err := gorm.Open(
		postgres.New(postgres.Config{}),
		&gorm.Config{
			DisableAutomaticPing:   true,
			DryRun:                 true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		t.Fatalf("open dry-run database: %v", err)
	}

	result := db.Clauses(defaultCommandsOnConflict()).Create(&model.ChannelsCommands{
		ID:        "command-id",
		ChannelID: "channel-id",
		Name:      "mmr",
	})
	if result.Error != nil {
		t.Fatalf("create default command: %v", result.Error)
	}

	query := strings.Join(strings.Fields(result.Statement.SQL.String()), " ")
	if !strings.Contains(query, `ON CONFLICT ("channelId","name") DO NOTHING`) {
		t.Fatalf("generated query %q does not use the channel/name DoNothing conflict clause", query)
	}
}
