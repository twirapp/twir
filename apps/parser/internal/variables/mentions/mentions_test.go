package mentions_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/variables/mentions"
	"github.com/twirapp/twir/libs/bus-core/twitch"
)

func TestMentionsID(t *testing.T) {
	t.Run("returns first mention ID by default", func(t *testing.T) {
		parseCtx := &types.VariableParseContext{
			ParseContext: &types.ParseContext{
				Mentions: []twitch.ChatMessageMessageFragmentMention{
					{UserId: "123", UserName: "User1", UserLogin: "user1"},
					{UserId: "456", UserName: "User2", UserLogin: "user2"},
				},
			},
		}

		result, err := mentions.ID.Handler(context.Background(), parseCtx, &types.VariableData{
			Key: "mentions.id",
		})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "123", result.Result)
	})

	t.Run("returns Nth mention ID when index provided", func(t *testing.T) {
		parseCtx := &types.VariableParseContext{
			ParseContext: &types.ParseContext{
				Mentions: []twitch.ChatMessageMessageFragmentMention{
					{UserId: "123", UserName: "User1", UserLogin: "user1"},
					{UserId: "456", UserName: "User2", UserLogin: "user2"},
				},
			},
		}

		index := "1"
		result, err := mentions.ID.Handler(context.Background(), parseCtx, &types.VariableData{
			Key:    "mentions.id",
			Params: &index,
		})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "456", result.Result)
	})

	t.Run("returns empty string when no mentions", func(t *testing.T) {
		parseCtx := &types.VariableParseContext{
			ParseContext: &types.ParseContext{
				Mentions: []twitch.ChatMessageMessageFragmentMention{},
			},
		}

		result, err := mentions.ID.Handler(context.Background(), parseCtx, &types.VariableData{
			Key: "mentions.id",
		})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "", result.Result)
	})

	t.Run("returns empty string when index out of bounds", func(t *testing.T) {
		parseCtx := &types.VariableParseContext{
			ParseContext: &types.ParseContext{
				Mentions: []twitch.ChatMessageMessageFragmentMention{
					{UserId: "123", UserName: "User1", UserLogin: "user1"},
				},
			},
		}

		index := "5"
		result, err := mentions.ID.Handler(context.Background(), parseCtx, &types.VariableData{
			Key:    "mentions.id",
			Params: &index,
		})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "", result.Result)
	})
}

func TestMentionsLogin(t *testing.T) {
	t.Run("returns first mention login by default", func(t *testing.T) {
		parseCtx := &types.VariableParseContext{
			ParseContext: &types.ParseContext{
				Mentions: []twitch.ChatMessageMessageFragmentMention{
					{UserId: "123", UserName: "User1", UserLogin: "user1"},
					{UserId: "456", UserName: "User2", UserLogin: "user2"},
				},
			},
		}

		result, err := mentions.Login.Handler(context.Background(), parseCtx, &types.VariableData{
			Key: "mentions.login",
		})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "user1", result.Result)
	})
}

func TestMentionsDisplayName(t *testing.T) {
	t.Run("returns first mention display name by default", func(t *testing.T) {
		parseCtx := &types.VariableParseContext{
			ParseContext: &types.ParseContext{
				Mentions: []twitch.ChatMessageMessageFragmentMention{
					{UserId: "123", UserName: "User1", UserLogin: "user1"},
					{UserId: "456", UserName: "User2", UserLogin: "user2"},
				},
			},
		}

		result, err := mentions.DisplayName.Handler(context.Background(), parseCtx, &types.VariableData{
			Key: "mentions.displayName",
		})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "User1", result.Result)
	})
}
