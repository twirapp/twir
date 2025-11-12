package user_creator

import (
	"context"
	"errors"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/repositories/users"
	usermodel "github.com/twirapp/twir/libs/repositories/users/model"
	usersstats "github.com/twirapp/twir/libs/repositories/users_stats"
	usersstatsmodel "github.com/twirapp/twir/libs/repositories/users_stats/model"
	userswithstatsrepository "github.com/twirapp/twir/libs/repositories/userswithstats"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	UsersStatsRepo     usersstats.Repository
	UsersRepo          users.Repository
	UsersWithStatsRepo userswithstatsrepository.Repository
	TrManager          trm.Manager
}

func New(opts Opts) *UserCreatorService {
	return &UserCreatorService{
		usersStatsRepo:     opts.UsersStatsRepo,
		usersRepo:          opts.UsersRepo,
		usersWithStatsRepo: opts.UsersWithStatsRepo,
		trManager:          opts.TrManager,
	}
}

type UserCreatorService struct {
	usersStatsRepo     usersstats.Repository
	usersRepo          users.Repository
	usersWithStatsRepo userswithstatsrepository.Repository
	trManager          trm.Manager
}

type CreateUserInput struct {
	UserID                   string
	ChannelID                *string
	Badges                   []twitch.ChatMessageBadge
	UsedEmotesWithThirdParty *int
	ShouldUpdateStats        bool
}

func (c *UserCreatorService) UnsureUser(ctx context.Context, input CreateUserInput) (
	*usermodel.User,
	*usersstatsmodel.UserStat,
	error,
) {
	if input.UserID == "" {
		return nil, nil, fmt.Errorf("UnsureUser: user_id is required")
	}

	isMod, isVip, isSubscriber := getUserRoles(input.Badges)

	// If there's no channel context, we only need to ensure the base user exists.
	if input.ChannelID == nil {
		user, err := c.ensureUserExists(ctx, input.UserID)
		if err != nil {
			return nil, nil, fmt.Errorf("UnsureUser: failed to ensure user without channelID: %w", err)
		}
		return user, nil, nil
	}

	// With a channel context, we handle both user and their channel-specific stats.
	return c.handleUserWithChannel(ctx, input, isMod, isVip, isSubscriber)
}

// handleUserWithChannel orchestrates the logic for a user within a specific channel.
func (c *UserCreatorService) handleUserWithChannel(
	ctx context.Context,
	input CreateUserInput,
	isMod, isVip, isSubscriber bool,
) (*usermodel.User, *usersstatsmodel.UserStat, error) {
	userWithStats, err := c.usersWithStatsRepo.GetByUserAndChannelID(
		ctx, userswithstatsrepository.GetByUserAndChannelIDInput{
			UserID:    input.UserID,
			ChannelID: *input.ChannelID,
		},
	)

	// Handle case where the user OR stats are not found.
	if err != nil {
		if errors.Is(err, userswithstatsrepository.ErrNotFound) {
			// This means the user might exist, but their stats for this channel do not.
			// Or the user doesn't exist at all. This function handles both cases.
			return c.findUserAndCreateStats(ctx, input, isMod, isVip, isSubscriber)
		}
		return nil, nil, fmt.Errorf("UnsureUser: failed to get user with stats: %w", err)
	}

	// User and stats record were found, proceed with updating stats.
	ensuredUser := &userWithStats.User
	ensuredStats, err := c.updateUserStats(ctx, input, isMod, isVip, isSubscriber)
	if err != nil {
		return nil, nil, fmt.Errorf("UnsureUser: failed to update user stats: %w", err)
	}

	return ensuredUser, ensuredStats, nil
}

// findUserAndCreateStats handles the scenario where stats for a user in a channel do not exist.
// It first ensures the base user record exists, then creates the stats record.
func (c *UserCreatorService) findUserAndCreateStats(
	ctx context.Context,
	input CreateUserInput,
	isMod, isVip, isSubscriber bool,
) (*usermodel.User, *usersstatsmodel.UserStat, error) {
	// Using a transaction to ensure both user creation (if needed) and stats creation succeed together.
	var ensuredUser *usermodel.User
	var ensuredStats *usersstatsmodel.UserStat

	err := c.trManager.Do(
		ctx, func(trCtx context.Context) error {
			var err error
			// This will find the user or create them if they don't exist.
			ensuredUser, err = c.ensureUserExists(trCtx, input.UserID)
			if err != nil {
				return fmt.Errorf("failed to ensure user exists in transaction: %w", err)
			}

			// Now that we're sure the user exists, create their stats for the channel.
			ensuredStats, err = c.createUserStats(trCtx, input, isMod, isVip, isSubscriber)
			if err != nil {
				return fmt.Errorf("failed to create stats in transaction: %w", err)
			}
			return nil
		},
	)

	if err != nil {
		return nil, nil, fmt.Errorf("UnsureUser transaction failed: %w", err)
	}

	return ensuredUser, ensuredStats, nil
}

// ensureUserExists gets a user by ID or creates them if they are not found.
func (c *UserCreatorService) ensureUserExists(ctx context.Context, userID string) (
	*usermodel.User,
	error,
) {
	user, err := c.usersRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, usermodel.ErrNotFound) {
			// User not found, create a new one.
			return c.createUser(
				ctx, users.CreateInput{
					ID:     userID,
					ApiKey: lo.ToPtr(uuid.NewString()),
				},
			)
		}
		// Another error occurred.
		return nil, err
	}
	return &user, nil
}

// createUserStats creates a new user statistics record.
func (c *UserCreatorService) createUserStats(
	ctx context.Context,
	input CreateUserInput,
	isMod, isVip, isSubscriber bool,
) (*usersstatsmodel.UserStat, error) {
	// On creation, set initial values. If the user sent a message, count starts at 1.
	initialMessages := 0
	if input.ShouldUpdateStats {
		initialMessages = 1
	}

	initialEmotes := 0
	if input.UsedEmotesWithThirdParty != nil {
		initialEmotes = *input.UsedEmotesWithThirdParty
	}

	return c.usersStatsRepo.Create(
		ctx, usersstats.CreateInput{
			UserID:            input.UserID,
			ChannelID:         *input.ChannelID,
			Messages:          int32(initialMessages),
			Emotes:            initialEmotes,
			IsMod:             isMod,
			IsVip:             isVip,
			IsSubscriber:      isSubscriber,
			UsedChannelPoints: 0,
			Watched:           0,
			Reputation:        0,
		},
	)
}

// updateUserStats updates an existing user statistics record.
func (c *UserCreatorService) updateUserStats(
	ctx context.Context,
	input CreateUserInput,
	isMod, isVip, isSubscriber bool,
) (*usersstatsmodel.UserStat, error) {
	// Prepare fields for incremental updates.
	numberFields := make(map[usersstats.IncrementInputFieldName]usersstats.NumberFieldUpdate)

	if input.ShouldUpdateStats {
		numberFields[usersstats.IncrementInputFieldMessages] = usersstats.NumberFieldUpdate{
			Count:       1,
			IsIncrement: true,
		}
	}
	if input.UsedEmotesWithThirdParty != nil {
		numberFields[usersstats.IncrementInputFieldEmotes] = usersstats.NumberFieldUpdate{
			Count:       *input.UsedEmotesWithThirdParty,
			IsIncrement: true,
		}
	}

	// Call the repository to update the record.
	return c.usersStatsRepo.CreateOrUpdate(
		ctx, input.UserID, *input.ChannelID, usersstats.UpdateInput{
			NumberFields: numberFields,
			IsMod:        &isMod,
			IsVip:        &isVip,
			IsSubscriber: &isSubscriber,
		},
	)
}

// createUser is a helper to create a new user record.
func (c *UserCreatorService) createUser(ctx context.Context, input users.CreateInput) (
	*usermodel.User,
	error,
) {
	newUser, err := c.usersRepo.Create(ctx, input)
	if err != nil {
		return nil, err
	}
	return &newUser, nil
}

// getUserRoles extracts role information from Twitch badges.
func getUserRoles(badges []twitch.ChatMessageBadge) (isMod, isVip, isSubscriber bool) {
	for _, b := range badges {
		switch b.SetId {
		case "moderator":
			isMod = true
		case "vip":
			isVip = true
		case "subscriber", "founder":
			isSubscriber = true
		}
	}
	return
}
