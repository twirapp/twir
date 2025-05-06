import * as Events from './events/events.ts'
import { Queue } from './queue.js'

import type { NatsConnection } from 'nats'

export function newBus(nc: NatsConnection) {
	return {
		Events: {
			Follow: new Queue<Events.FollowMessage, any>(nc, Events.FollowSubject),
			Subscribe: new Queue<Events.SubscribeMessage, any>(nc, Events.SubscribeSubject),
			SubGift: new Queue<Events.SubGiftMessage, any>(nc, Events.SubGiftSubject),
			ReSubscribe: new Queue<Events.ReSubscribeMessage, any>(nc, Events.ReSubscribeSubject),
			RedemptionCreated: new Queue<Events.RedemptionCreatedMessage, any>(nc, Events.RedemptionCreatedSubject),
			CommandUsed: new Queue<Events.CommandUsedMessage, any>(nc, Events.CommandUsedSubject),
			FirstUserMessage: new Queue<Events.FirstUserMessageMessage, any>(nc, Events.FirstUserMessageSubject),
			Raided: new Queue<Events.RaidedMessage, any>(nc, Events.RaidedSubject),
			TitleOrCategoryChanged: new Queue<Events.TitleOrCategoryChangedMessage, any>(nc, Events.TitleOrCategoryChangedSubject),
			StreamOnline: new Queue<Events.StreamOnlineMessage, any>(nc, Events.StreamOnlineSubject),
			StreamOffline: new Queue<Events.StreamOfflineMessage, any>(nc, Events.StreamOfflineSubject),
			ChatClear: new Queue<Events.ChatClearMessage, any>(nc, Events.ChatClearSubject),
			Donate: new Queue<Events.DonateMessage, any>(nc, Events.DonateSubject),
			KeywordMatched: new Queue<Events.KeywordMatchedMessage, any>(nc, Events.KeywordMatchedSubject),
			GreetingSended: new Queue<Events.GreetingSendedMessage, any>(nc, Events.GreetingSendedSubject),
			PollBegin: new Queue<Events.PollBeginMessage, any>(nc, Events.PollBeginSubject),
			PollProgress: new Queue<Events.PollProgressMessage, any>(nc, Events.PollProgressSubject),
			PollEnd: new Queue<Events.PollEndMessage, any>(nc, Events.PollEndSubject),
			PredictionBegin: new Queue<Events.PredictionBeginMessage, any>(nc, Events.PredictionBeginSubject),
			PredictionProgress: new Queue<Events.PredictionProgressMessage, any>(nc, Events.PredictionProgressSubject),
			PredictionLock: new Queue<Events.PredictionLockMessage, any>(nc, Events.PredictionLockSubject),
			PredictionEnd: new Queue<Events.PredictionEndMessage, any>(nc, Events.PredictionEndSubject),
			StreamFirstUserJoin: new Queue<Events.StreamFirstUserJoinMessage, any>(nc, Events.StreamFirstUserJoinSubject),
			ChannelBan: new Queue<Events.ChannelBanMessage, any>(nc, Events.ChannelBanSubject),
			ChannelUnbanRequestCreate: new Queue<Events.ChannelUnbanRequestCreateMessage, any>(nc, Events.ChannelUnbanRequestCreateSubject),
			ChannelUnbanRequestResolve: new Queue<Events.ChannelUnbanRequestResolveMessage, any>(nc, Events.ChannelUnbanRequestResolveSubject),
			ChannelMessageDelete: new Queue<Events.ChannelMessageDeleteMessage, any>(nc, Events.ChannelMessageDeleteSubject),
			VipAdded: new Queue<Events.VipAddedMessage, any>(nc, Events.VipAddedSubject),
			VipRemoved: new Queue<Events.VipRemovedMessage, any>(nc, Events.VipRemovedSubject),
			ModeratorAdded: new Queue<Events.ModeratorAddedMessage, any>(nc, Events.ModeratorAddedSubject),
			ModeratorRemoved: new Queue<Events.ModeratorRemovedMessage, any>(nc, Events.ModeratorRemovedSubject),
		},
	}
}
