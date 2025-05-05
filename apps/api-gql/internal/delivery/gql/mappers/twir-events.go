package mappers

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/twitch"
)

func TwirEventBaseInfoToGql(channelID, channelName string) *gqlmodel.EventBaseInfo {
	return &gqlmodel.EventBaseInfo{
		ChannelID:   channelID,
		ChannelName: channelName,
	}
}

func TwirEventFollowToGql(event events.FollowMessage) gqlmodel.EventFollowMessage {
	return gqlmodel.EventFollowMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserName:        event.UserName,
		UserDisplayName: event.UserDisplayName,
		UserID:          event.UserID,
	}
}

func TwirEventSubscribeToGql(event events.SubscribeMessage) gqlmodel.EventSubscribeMessage {
	return gqlmodel.EventSubscribeMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserName:        event.UserName,
		UserDisplayName: event.UserDisplayName,
		Level:           event.Level,
		UserID:          event.UserID,
	}
}

func TwirEventSubGiftToGql(event events.SubGiftMessage) gqlmodel.EventSubGiftMessage {
	return gqlmodel.EventSubGiftMessage{
		BaseInfo:          TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		SenderUserName:    event.SenderUserName,
		SenderDisplayName: event.SenderDisplayName,
		TargetUserName:    event.TargetUserName,
		TargetDisplayName: event.TargetDisplayName,
		Level:             event.Level,
		SenderUserID:      event.SenderUserID,
	}
}

func TwirEventReSubscribeToGql(event events.ReSubscribeMessage) gqlmodel.EventReSubscribeMessage {
	return gqlmodel.EventReSubscribeMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserName:        event.UserName,
		UserDisplayName: event.UserDisplayName,
		Months:          int(event.Months),
		Streak:          int(event.Streak),
		IsPrime:         event.IsPrime,
		Message:         event.Message,
		Level:           event.Level,
		UserID:          event.UserID,
	}
}

func TwirEventRedemptionCreatedToGql(event events.RedemptionCreatedMessage) gqlmodel.EventRedemptionCreatedMessage {
	return gqlmodel.EventRedemptionCreatedMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserName:        event.UserName,
		UserDisplayName: event.UserDisplayName,
		ID:              event.ID,
		RewardName:      event.RewardName,
		RewardCost:      event.RewardCost,
		Input:           event.Input,
		UserID:          event.UserID,
	}
}

func TwirEventCommandUsedToGql(event events.CommandUsedMessage) gqlmodel.EventCommandUsedMessage {
	return gqlmodel.EventCommandUsedMessage{
		BaseInfo: TwirEventBaseInfoToGql(
			event.BaseInfo.ChannelID,
			event.BaseInfo.ChannelName,
		),
		CommandID:          event.CommandID,
		CommandName:        event.CommandName,
		UserName:           event.UserName,
		UserDisplayName:    event.UserDisplayName,
		CommandInput:       event.CommandInput,
		UserID:             event.UserID,
		IsDefault:          event.IsDefault,
		DefaultCommandName: event.DefaultCommandName,
		MessageID:          event.MessageID,
	}
}

func TwirEventFirstUserMessageToGql(event events.FirstUserMessageMessage) gqlmodel.EventFirstUserMessageMessage {
	return gqlmodel.EventFirstUserMessageMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserID:          event.UserID,
		UserName:        event.UserName,
		UserDisplayName: event.UserDisplayName,
		MessageID:       event.MessageID,
	}
}

func TwirEventRaidedToGql(event events.RaidedMessage) gqlmodel.EventRaidedMessage {
	return gqlmodel.EventRaidedMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserName:        event.UserName,
		UserDisplayName: event.UserDisplayName,
		Viewers:         int(event.Viewers),
		UserID:          event.UserID,
	}
}

func TwirEventTitleOrCategoryChangedToGql(event events.TitleOrCategoryChangedMessage) gqlmodel.EventTitleOrCategoryChangedMessage {
	return gqlmodel.EventTitleOrCategoryChangedMessage{
		BaseInfo:    TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		OldTitle:    event.OldTitle,
		NewTitle:    event.NewTitle,
		OldCategory: event.OldCategory,
		NewCategory: event.NewCategory,
	}
}

func TwirEventChatClearToGql(event events.ChatClearMessage) gqlmodel.EventChatClearMessage {
	return gqlmodel.EventChatClearMessage{
		BaseInfo: TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
	}
}

func TwirEventDonateToGql(event events.DonateMessage) gqlmodel.EventDonateMessage {
	return gqlmodel.EventDonateMessage{
		BaseInfo: TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserName: event.UserName,
		Amount:   event.Amount,
		Currency: event.Currency,
		Message:  event.Message,
	}
}

func TwirEventKeywordMatchedToGql(event events.KeywordMatchedMessage) gqlmodel.EventKeywordMatchedMessage {
	return gqlmodel.EventKeywordMatchedMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		KeywordID:       event.KeywordID,
		KeywordName:     event.KeywordName,
		KeywordResponse: event.KeywordResponse,
		UserID:          event.UserID,
		UserName:        event.UserName,
		UserDisplayName: event.UserDisplayName,
	}
}

func TwirEventGreetingSendedToGql(event events.GreetingSendedMessage) gqlmodel.EventGreetingSendedMessage {
	return gqlmodel.EventGreetingSendedMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserID:          event.UserID,
		UserName:        event.UserName,
		UserDisplayName: event.UserDisplayName,
		GreetingText:    event.GreetingText,
	}
}

func TwirEventPollChoiceToGql(choice events.PollChoice) gqlmodel.EventPollChoice {
	return gqlmodel.EventPollChoice{
		ID:                  choice.ID,
		Title:               choice.Title,
		BitsVotes:           int(choice.BitsVotes),
		ChannelsPointsVotes: int(choice.ChannelsPointsVotes),
		Votes:               int(choice.Votes),
	}
}

func TwirEventPollBitsVotesToGql(bitsVoting events.PollBitsVotes) *gqlmodel.EventPollBitsVotes {
	return &gqlmodel.EventPollBitsVotes{
		Enabled:       bitsVoting.Enabled,
		AmountPerVote: int(bitsVoting.AmountPerVote),
	}
}

func TwirEventPollChannelPointsVotesToGql(channelPointsVoting events.PollChannelPointsVotes) *gqlmodel.EventPollChannelPointsVotes {
	return &gqlmodel.EventPollChannelPointsVotes{
		Enabled:       channelPointsVoting.Enabled,
		AmountPerVote: int(channelPointsVoting.AmountPerVote),
	}
}

func TwirEventPollInfoToGql(info events.PollInfo) *gqlmodel.EventPollInfo {
	choices := make([]gqlmodel.EventPollChoice, 0, len(info.Choices))
	for _, choice := range info.Choices {
		choices = append(choices, TwirEventPollChoiceToGql(choice))
	}

	return &gqlmodel.EventPollInfo{
		Title:                info.Title,
		Choices:              choices,
		BitsVoting:           TwirEventPollBitsVotesToGql(info.BitsVoting),
		ChannelsPointsVoting: TwirEventPollChannelPointsVotesToGql(info.ChannelsPointsVoting),
	}
}

func TwirEventPollBeginToGql(event events.PollBeginMessage) gqlmodel.EventPollBeginMessage {
	return gqlmodel.EventPollBeginMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserName:        event.UserName,
		UserDisplayName: event.UserDisplayName,
		Info:            TwirEventPollInfoToGql(event.Info),
	}
}

func TwirEventPollProgressToGql(event events.PollProgressMessage) gqlmodel.EventPollProgressMessage {
	return gqlmodel.EventPollProgressMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserName:        event.UserName,
		UserDisplayName: event.UserDisplayName,
		Info:            TwirEventPollInfoToGql(event.Info),
	}
}

func TwirEventPollEndToGql(event events.PollEndMessage) gqlmodel.EventPollEndMessage {
	return gqlmodel.EventPollEndMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserName:        event.UserName,
		UserDisplayName: event.UserDisplayName,
		Info:            TwirEventPollInfoToGql(event.Info),
	}
}

func TwirEventPredictionTopPredictorToGql(predictor events.PredictionTopPredictor) gqlmodel.EventPredictionTopPredictor {
	var pointsWin *int
	if predictor.PointsWin != nil {
		pointsWin = lo.ToPtr(int(*predictor.PointsWin))
	}

	return gqlmodel.EventPredictionTopPredictor{
		UserName:        predictor.UserName,
		UserDisplayName: predictor.UserDisplayName,
		UserID:          predictor.UserID,
		PointsUsed:      int(predictor.PointsUsed),
		PointsWin:       pointsWin,
	}
}

func TwirEventPredictionOutcomeToGql(outcome events.PredictionOutcome) gqlmodel.EventPredictionOutcome {
	topPredictors := make([]gqlmodel.EventPredictionTopPredictor, 0, len(outcome.TopPredictors))
	for _, predictor := range outcome.TopPredictors {
		topPredictors = append(topPredictors, TwirEventPredictionTopPredictorToGql(predictor))
	}

	return gqlmodel.EventPredictionOutcome{
		ID:            outcome.ID,
		Title:         outcome.Title,
		Color:         outcome.Color,
		Users:         int(outcome.Users),
		ChannelPoints: int(outcome.ChannelPoints),
		TopPredictors: topPredictors,
	}
}

func TwirEventPredictionInfoToGql(info events.PredictionInfo) *gqlmodel.EventPredictionInfo {
	outcomes := make([]gqlmodel.EventPredictionOutcome, 0, len(info.Outcomes))
	for _, outcome := range info.Outcomes {
		outcomes = append(outcomes, TwirEventPredictionOutcomeToGql(outcome))
	}

	return &gqlmodel.EventPredictionInfo{
		Title:    info.Title,
		Outcomes: outcomes,
	}
}

func TwirEventPredictionBeginToGql(event events.PredictionBeginMessage) gqlmodel.EventPredictionBeginMessage {
	return gqlmodel.EventPredictionBeginMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserName:        event.UserName,
		UserDisplayName: event.UserDisplayName,
		Info:            TwirEventPredictionInfoToGql(event.Info),
	}
}

func TwirEventPredictionProgressToGql(event events.PredictionProgressMessage) gqlmodel.EventPredictionProgressMessage {
	return gqlmodel.EventPredictionProgressMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserName:        event.UserName,
		UserDisplayName: event.UserDisplayName,
		Info:            TwirEventPredictionInfoToGql(event.Info),
	}
}

func TwirEventPredictionLockToGql(event events.PredictionLockMessage) gqlmodel.EventPredictionLockMessage {
	return gqlmodel.EventPredictionLockMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserName:        event.UserName,
		UserDisplayName: event.UserDisplayName,
		Info:            TwirEventPredictionInfoToGql(event.Info),
	}
}

func TwirEventPredictionEndToGql(event events.PredictionEndMessage) gqlmodel.EventPredictionEndMessage {
	return gqlmodel.EventPredictionEndMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserName:        event.UserName,
		UserDisplayName: event.UserDisplayName,
		Info:            TwirEventPredictionInfoToGql(event.Info),
	}
}

func TwirEventStreamOnlineToGql(event twitch.StreamOnlineMessage) gqlmodel.EventStreamOnlineMessage {
	return gqlmodel.EventStreamOnlineMessage{
		BaseInfo:  TwirEventBaseInfoToGql(event.ChannelID, ""),
		StartedAt: event.StartedAt,
	}
}

func TwirEventStreamOfflineToGql(event twitch.StreamOfflineMessage) gqlmodel.EventStreamOfflineMessage {
	return gqlmodel.EventStreamOfflineMessage{
		BaseInfo: TwirEventBaseInfoToGql(event.ChannelID, ""),
		EndedAt:  time.Now(),
	}
}

func TwirEventStreamFirstUserJoinToGql(event events.StreamFirstUserJoinMessage) gqlmodel.EventStreamFirstUserJoinMessage {
	return gqlmodel.EventStreamFirstUserJoinMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserID:          event.UserID,
		UserName:        event.UserLogin,
		UserDisplayName: event.UserLogin,
	}
}

func TwirEventChannelBanToGql(event events.ChannelBanMessage) gqlmodel.EventChannelBanMessage {
	var expiresAt *string
	if event.EndsAt != "" {
		expiresAt = lo.ToPtr(event.EndsAt)
	}

	return gqlmodel.EventChannelBanMessage{
		BaseInfo: TwirEventBaseInfoToGql(
			event.BaseInfo.ChannelID,
			event.BaseInfo.ChannelName,
		),
		UserID:               event.UserID,
		UserName:             event.UserLogin,
		UserDisplayName:      event.UserName,
		ModeratorID:          event.ModeratorUserID,
		ModeratorName:        event.ModeratorUserLogin,
		ModeratorDisplayName: event.ModeratorUserName,
		Reason:               &event.Reason,
		IsPermanent:          event.IsPermanent,
		ExpiresAt:            expiresAt,
	}
}

func TwirEventChannelUnbanRequestCreateToGql(event events.ChannelUnbanRequestCreateMessage) gqlmodel.EventChannelUnbanRequestCreateMessage {
	return gqlmodel.EventChannelUnbanRequestCreateMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserID:          "",
		UserName:        event.UserLogin,
		UserDisplayName: event.UserName,
		Message:         event.Text,
	}
}

func TwirEventChannelUnbanRequestResolveToGql(event events.ChannelUnbanRequestResolveMessage) gqlmodel.EventChannelUnbanRequestResolveMessage {
	return gqlmodel.EventChannelUnbanRequestResolveMessage{
		BaseInfo: TwirEventBaseInfoToGql(
			event.BaseInfo.ChannelID,
			event.BaseInfo.ChannelName,
		),
		UserID:               event.UserID,
		UserName:             event.UserLogin,
		UserDisplayName:      event.UserName,
		ModeratorID:          event.ModeratorUserID,
		ModeratorName:        event.ModeratorUserLogin,
		ModeratorDisplayName: event.ModeratorUserName,
		Approved:             !event.Declined,
	}
}

func TwirEventChannelMessageDeleteToGql(event events.ChannelMessageDeleteMessage) gqlmodel.EventChannelMessageDeleteMessage {
	return gqlmodel.EventChannelMessageDeleteMessage{
		BaseInfo: TwirEventBaseInfoToGql(
			event.BaseInfo.ChannelID,
			event.BaseInfo.ChannelName,
		),
		MessageID:            event.MessageId,
		UserID:               event.UserId,
		UserName:             event.UserLogin,
		UserDisplayName:      event.UserName,
		ModeratorID:          event.BaseInfo.ChannelID,
		ModeratorName:        event.BroadcasterUserLogin,
		ModeratorDisplayName: event.BroadcasterUserName,
	}
}

func TwirEventVipAddedToGql(event events.VipAddedMessage) gqlmodel.EventVipAddedMessage {
	return gqlmodel.EventVipAddedMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserID:          event.UserID,
		UserName:        event.UserName,
		UserDisplayName: event.UserName,
	}
}

func TwirEventVipRemovedToGql(event events.VipRemovedMessage) gqlmodel.EventVipRemovedMessage {
	return gqlmodel.EventVipRemovedMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserID:          event.UserID,
		UserName:        event.UserName,
		UserDisplayName: event.UserName,
	}
}

func TwirEventModeratorAddedToGql(event events.ModeratorAddedMessage) gqlmodel.EventModeratorAddedMessage {
	return gqlmodel.EventModeratorAddedMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserID:          event.UserID,
		UserName:        event.UserName,
		UserDisplayName: event.UserName,
	}
}

func TwirEventModeratorRemovedToGql(event events.ModeratorRemovedMessage) gqlmodel.EventModeratorRemovedMessage {
	return gqlmodel.EventModeratorRemovedMessage{
		BaseInfo:        TwirEventBaseInfoToGql(event.BaseInfo.ChannelID, event.BaseInfo.ChannelName),
		UserID:          event.UserID,
		UserName:        event.UserName,
		UserDisplayName: event.UserName,
	}
}

func MapEventToGqlType(eventName string, data []byte) (gqlmodel.EventMessage, error) {
	switch eventName {
	case events.FollowSubject:
		followMessage := events.FollowMessage{}
		err := json.Unmarshal(data, &followMessage)
		if err != nil {
			return nil, err
		}

		return TwirEventFollowToGql(followMessage), nil
	case events.SubscribeSubject:
		subscribeMessage := events.SubscribeMessage{}
		err := json.Unmarshal(data, &subscribeMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventSubscribeToGql(subscribeMessage), nil
	case events.ReSubscribeSubject:
		reSubscribeMessage := events.ReSubscribeMessage{}
		err := json.Unmarshal(data, &reSubscribeMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventReSubscribeToGql(reSubscribeMessage), nil
	case events.TitleOrCategoryChangedSubject:
		titleOrCategoryChangedMessage := events.TitleOrCategoryChangedMessage{}
		err := json.Unmarshal(data, &titleOrCategoryChangedMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventTitleOrCategoryChangedToGql(titleOrCategoryChangedMessage), nil
	case events.ChatClearSubject:
		chatClearMessage := events.ChatClearMessage{}
		err := json.Unmarshal(data, &chatClearMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventChatClearToGql(chatClearMessage), nil
	case events.DonateSubject:
		donateMessage := events.DonateMessage{}
		err := json.Unmarshal(data, &donateMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventDonateToGql(donateMessage), nil
	case events.GreetingSendedSubject:
		greetingSendedMessage := events.GreetingSendedMessage{}
		err := json.Unmarshal(data, &greetingSendedMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventGreetingSendedToGql(greetingSendedMessage), nil
	case events.SubGiftSubject:
		subGiftMessage := events.SubGiftMessage{}
		err := json.Unmarshal(data, &subGiftMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventSubGiftToGql(subGiftMessage), nil
	case events.RedemptionCreatedSubject:
		redemptionCreatedMessage := events.RedemptionCreatedMessage{}
		err := json.Unmarshal(data, &redemptionCreatedMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventRedemptionCreatedToGql(redemptionCreatedMessage), nil
	case events.CommandUsedSubject:
		commandUsedMessage := events.CommandUsedMessage{}
		err := json.Unmarshal(data, &commandUsedMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventCommandUsedToGql(commandUsedMessage), nil
	case events.FirstUserMessageSubject:
		firstUserMessageMessage := events.FirstUserMessageMessage{}
		err := json.Unmarshal(data, &firstUserMessageMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventFirstUserMessageToGql(firstUserMessageMessage), nil
	case events.RaidedSubject:
		raidedMessage := events.RaidedMessage{}
		err := json.Unmarshal(data, &raidedMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventRaidedToGql(raidedMessage), nil
	case events.PredictionBeginSubject:
		predictionBeginMessage := events.PredictionBeginMessage{}
		err := json.Unmarshal(data, &predictionBeginMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventPredictionBeginToGql(predictionBeginMessage), nil
	case events.PredictionProgressSubject:
		predictionProgressMessage := events.PredictionProgressMessage{}
		err := json.Unmarshal(data, &predictionProgressMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventPredictionProgressToGql(predictionProgressMessage), nil
	case events.PredictionLockSubject:
		predictionLockMessage := events.PredictionLockMessage{}
		err := json.Unmarshal(data, &predictionLockMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventPredictionLockToGql(predictionLockMessage), nil
	case events.PredictionEndSubject:
		predictionEndMessage := events.PredictionEndMessage{}
		err := json.Unmarshal(data, &predictionEndMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventPredictionEndToGql(predictionEndMessage), nil
	case events.StreamFirstUserJoinSubject:
		streamFirstUserJoinMessage := events.StreamFirstUserJoinMessage{}
		err := json.Unmarshal(data, &streamFirstUserJoinMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventStreamFirstUserJoinToGql(streamFirstUserJoinMessage), nil
	case events.ChannelBanSubject:
		channelBanMessage := events.ChannelBanMessage{}
		err := json.Unmarshal(data, &channelBanMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventChannelBanToGql(channelBanMessage), nil
	case events.ChannelUnbanRequestCreateSubject:
		channelUnbanRequestCreateMessage := events.ChannelUnbanRequestCreateMessage{}
		err := json.Unmarshal(data, &channelUnbanRequestCreateMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventChannelUnbanRequestCreateToGql(channelUnbanRequestCreateMessage), nil
	case events.ChannelUnbanRequestResolveSubject:
		channelUnbanRequestResolveMessage := events.ChannelUnbanRequestResolveMessage{}
		err := json.Unmarshal(data, &channelUnbanRequestResolveMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventChannelUnbanRequestResolveToGql(channelUnbanRequestResolveMessage), nil
	case events.ChannelMessageDeleteSubject:
		channelMessageDeleteMessage := events.ChannelMessageDeleteMessage{}
		err := json.Unmarshal(data, &channelMessageDeleteMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventChannelMessageDeleteToGql(channelMessageDeleteMessage), nil
	case events.PollProgressSubject:
		pollProgressMessage := events.PollProgressMessage{}
		err := json.Unmarshal(data, &pollProgressMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventPollProgressToGql(pollProgressMessage), nil
	case events.PollEndSubject:
		pollEndMessage := events.PollEndMessage{}
		err := json.Unmarshal(data, &pollEndMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventPollEndToGql(pollEndMessage), nil
	case events.StreamOnlineSubject:
		streamOnlineMessage := twitch.StreamOnlineMessage{}
		err := json.Unmarshal(data, &streamOnlineMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventStreamOnlineToGql(streamOnlineMessage), nil
	case events.StreamOfflineSubject:
		streamOfflineMessage := twitch.StreamOfflineMessage{}
		err := json.Unmarshal(data, &streamOfflineMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventStreamOfflineToGql(streamOfflineMessage), nil
	case events.VipAddedSubject:
		vipAddedMessage := events.VipAddedMessage{}
		err := json.Unmarshal(data, &vipAddedMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventVipAddedToGql(vipAddedMessage), nil
	case events.VipRemovedSubject:
		vipRemovedMessage := events.VipRemovedMessage{}
		err := json.Unmarshal(data, &vipRemovedMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventVipRemovedToGql(vipRemovedMessage), nil
	case events.ModeratorAddedSubject:
		moderatorAddedMessage := events.ModeratorAddedMessage{}
		err := json.Unmarshal(data, &moderatorAddedMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventModeratorAddedToGql(moderatorAddedMessage), nil
	case events.ModeratorRemovedSubject:
		moderatorRemovedMessage := events.ModeratorRemovedMessage{}
		err := json.Unmarshal(data, &moderatorRemovedMessage)
		if err != nil {
			return nil, err
		}
		return TwirEventModeratorRemovedToGql(moderatorRemovedMessage), nil
	}

	return nil, fmt.Errorf("unknown event: %s", eventName)
}
