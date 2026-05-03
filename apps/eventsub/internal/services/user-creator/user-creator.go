package user_creator

import (
	"context"
	"errors"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/bus-core/generic"
	"github.com/twirapp/twir/libs/entities/platform"
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
	PlatformID               string
	Platform                 platform.Platform
	ChannelID                *string
	Badges                   []generic.ChatMessageBadge
	UsedEmotesWithThirdParty *int
	ShouldUpdateStats        bool
	IsBroadcaster            bool
	IsModerator              bool
	IsVip                    bool
	IsSubscriber             bool
}

func (c *UserCreatorService) UnsureUser(ctx context.Context, input CreateUserInput) (
	*usermodel.User,
	*usersstatsmodel.UserStat,
	error,
) {
	if input.UserID == "" {
		return nil, nil, fmt.Errorf("UnsureUser: user_id is required")
	}

	// If there's no channel context, we only need to ensure the base user exists.
	if input.ChannelID == nil {
		user, err := c.ensureUserExists(ctx, input)
		if err != nil {
			return nil, nil, fmt.Errorf("UnsureUser: failed to ensure user without channelID: %w", err)
		}
		return user, nil, nil
	}

	// With a channel context, we handle both user and their channel-specific stats.
	return c.handleUserWithChannel(ctx, input)
}

// handleUserWithChannel orchestrates the logic for a user within a specific channel.
func (c *UserCreatorService) handleUserWithChannel(
	ctx context.Context,
	input CreateUserInput,
) (*usermodel.User, *usersstatsmodel.UserStat, error) {
	ensuredUser, err := c.ensureUserExists(ctx, input)
	if err != nil {
		return nil, nil, fmt.Errorf("UnsureUser: failed to ensure user exists: %w", err)
	}

	parsedChannelID, err := uuid.Parse(*input.ChannelID)
	if err != nil {
		return nil, nil, fmt.Errorf("UnsureUser: failed to parse channel id: %w", err)
	}

	_, err = c.usersWithStatsRepo.GetByUserAndChannelID(
		ctx, userswithstatsrepository.GetByUserAndChannelIDInput{
			UserID:    ensuredUser.ID,
			ChannelID: parsedChannelID,
		},
	)
	// Handle case where the user OR stats are not found.
	if err != nil {
		if errors.Is(err, userswithstatsrepository.ErrNotFound) {
			// This means the user might exist, but their stats for this channel do not.
			// Or the user doesn't exist at all. This function handles both cases.
			return c.createStatsForExistingUser(
				ctx,
				ensuredUser,
				parsedChannelID,
				input,
				input.IsModerator,
				input.IsVip,
				input.IsSubscriber,
			)
		}
		return nil, nil, fmt.Errorf("UnsureUser: failed to get user with stats: %w", err)
	}

	// User and stats record were found, proceed with updating stats.
	ensuredStats, err := c.updateUserStats(
		ctx,
		ensuredUser.ID,
		parsedChannelID,
		input,
		input.IsModerator,
		input.IsVip,
		input.IsSubscriber,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("UnsureUser: failed to update user stats: %w", err)
	}

	return ensuredUser, ensuredStats, nil
}

// findUserAndCreateStats handles the scenario where stats for a user in a channel do not exist.
// It first ensures the base user record exists, then creates the stats record.
func (c *UserCreatorService) createStatsForExistingUser(
	ctx context.Context,
	ensuredUser *usermodel.User,
	channelID uuid.UUID,
	input CreateUserInput,
	isMod, isVip, isSubscriber bool,
) (*usermodel.User, *usersstatsmodel.UserStat, error) {
	// Using a transaction to ensure stats creation succeeds consistently.
	var ensuredStats *usersstatsmodel.UserStat

	err := c.trManager.Do(
		ctx, func(trCtx context.Context) error {
			var err error
			ensuredStats, err = c.createUserStats(trCtx, ensuredUser.ID, channelID, input, isMod, isVip, isSubscriber)
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
func (c *UserCreatorService) ensureUserExists(ctx context.Context, input CreateUserInput) (
	*usermodel.User,
	error,
) {
	if parsedUserID, parseErr := uuid.Parse(input.UserID); parseErr == nil {
		user, err := c.usersRepo.GetByID(ctx, parsedUserID)
		if err == nil {
			return &user, nil
		}
		if !errors.Is(err, usermodel.ErrNotFound) {
			return nil, err
		}
	}

	plat := input.Platform
	if plat == "" {
		plat = platform.PlatformTwitch
	}

	platformID := input.PlatformID
	if platformID == "" {
		platformID = input.UserID
	}

	if platformID != "" {
		user, err := c.usersRepo.GetByPlatformID(ctx, plat, platformID)
		if err == nil {
			return &user, nil
		}
		if !errors.Is(err, usermodel.ErrNotFound) {
			return nil, err
		}
	}

	return c.createUser(
		ctx, users.CreateInput{
			Platform:   plat,
			PlatformID: platformID,
		},
	)
}

func (c *UserCreatorService) createUserStats(
	ctx context.Context,
	userID uuid.UUID,
	channelID uuid.UUID,
	input CreateUserInput,
	isMod, isVip, isSubscriber bool,
) (*usersstatsmodel.UserStat, error) {
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
			UserID:            userID,
			ChannelID:         channelID,
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

// createUserStats creates a new user statistics record.
// updateUserStats updates an existing user statistics record.
func (c *UserCreatorService) updateUserStats(
	ctx context.Context,
	userID uuid.UUID,
	channelID uuid.UUID,
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
		ctx, userID, channelID, usersstats.UpdateInput{
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
