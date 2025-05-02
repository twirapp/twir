export const FollowSubject = 'events.follow'
export const SubscribeSubject = 'events.subscribe'
export const SubGiftSubject = 'events.sub_gift'
export const ReSubscribeSubject = 'events.resubscribe'
export const RedemptionCreatedSubject = 'events.redemption_created'
export const CommandUsedSubject = 'events.command_used'
export const FirstUserMessageSubject = 'events.first_user_message'
export const RaidedSubject = 'events.raided'
export const TitleOrCategoryChangedSubject = 'events.title_or_category_changed'
export const StreamOnlineSubject = 'events.stream_online'
export const StreamOfflineSubject = 'events.stream_offline'
export const ChatClearSubject = 'events.chat_clear'
export const DonateSubject = 'events.donate'
export const KeywordMatchedSubject = 'events.keyword_matched'
export const GreetingSendedSubject = 'events.greeting_sended'
export const PollBeginSubject = 'events.poll_begin'
export const PollProgressSubject = 'events.poll_progress'
export const PollEndSubject = 'events.poll_end'
export const PredictionBeginSubject = 'events.prediction_begin'
export const PredictionProgressSubject = 'events.prediction_progress'
export const PredictionLockSubject = 'events.prediction_lock'
export const PredictionEndSubject = 'events.prediction_end'
export const StreamFirstUserJoinSubject = 'events.stream_first_user_join'
export const ChannelBanSubject = 'events.channel_ban'
export const ChannelUnbanRequestCreateSubject = 'events.channel_unban_request_create'
export const ChannelUnbanRequestResolveSubject = 'events.channel_unban_request_resolve'
export const ChannelMessageDeleteSubject = 'events.channel_message_delete'
export const ModeratorAddedSubject = 'events.moderator_added'
export const ModeratorRemovedSubject = 'events.moderator_removed'
export const VipAddedSubject = 'events.vip_added'
export const VipRemovedSubject = 'events.vip_removed'

export interface BaseInfo {
	channel_id: string
	channel_name: string
}

export interface FollowMessage {
	base_info: BaseInfo
	user_name: string
	user_display_name: string
	user_id: string
}

export interface SubscribeMessage {
	base_info: BaseInfo
	user_name: string
	user_display_name: string
	level: string
	user_id: string
}

export interface SubGiftMessage {
	base_info: BaseInfo
	sender_user_name: string
	sender_display_name: string
	target_user_name: string
	target_display_name: string
	level: string
	sender_user_id: string
}

export interface ReSubscribeMessage {
	base_info: BaseInfo
	user_name: string
	user_display_name: string
	months: number
	streak: number
	is_prime: boolean
	message: string
	level: string
	user_id: string
}

export interface RedemptionCreatedMessage {
	base_info: BaseInfo
	user_name: string
	user_display_name: string
	id: string
	reward_name: string
	reward_cost: string
	input?: string
	user_id: string
}

export interface CommandUsedMessage {
	base_info: BaseInfo
	command_id: string
	command_name: string
	user_name: string
	user_display_name: string
	command_input: string
	user_id: string
	is_default: boolean
	default_command_name: string
	message_id: string
}

export interface FirstUserMessageMessage {
	base_info: BaseInfo
	user_id: string
	user_name: string
	user_display_name: string
	message_id: string
}

export interface RaidedMessage {
	base_info: BaseInfo
	user_name: string
	user_display_name: string
	viewers: number
	user_id: string
}

export interface TitleOrCategoryChangedMessage {
	base_info: BaseInfo
	old_title: string
	new_title: string
	old_category: string
	new_category: string
}

export interface StreamOnlineMessage {
	base_info: BaseInfo
	title: string
	category: string
}

export interface StreamOfflineMessage {
	base_info: BaseInfo
}

export interface ChatClearMessage {
	base_info: BaseInfo
}

export interface DonateMessage {
	base_info: BaseInfo
	user_name: string
	amount: string
	currency: string
	message: string
}

export interface KeywordMatchedMessage {
	base_info: BaseInfo
	keyword_id: string
	keyword_name: string
	keyword_response: string
	user_id: string
	user_name: string
	user_display_name: string
}

export interface GreetingSendedMessage {
	base_info: BaseInfo
	user_id: string
	user_name: string
	user_display_name: string
	greeting_text: string
}

export interface PollChoice {
	id: string
	title: string
	bits_votes: number
	channels_points_votes: number
	votes: number
}

export interface PollBitsVotes {
	enabled: boolean
	amount_per_vote: number
}

export interface PollChannelPointsVotes {
	enabled: boolean
	amount_per_vote: number
}

export interface PollInfo {
	title: string
	choices: PollChoice[]
	bits_voting: PollBitsVotes
	channels_points_voting: PollChannelPointsVotes
}

export interface PollBeginMessage {
	base_info: BaseInfo
	user_name: string
	user_display_name: string
	info: PollInfo
}

export interface PollProgressMessage {
	base_info: BaseInfo
	user_name: string
	user_display_name: string
	info: PollInfo
}

export interface PollEndMessage {
	base_info: BaseInfo
	user_name: string
	user_display_name: string
	info: PollInfo
}

export interface PredictionTopPredictor {
	user_name: string
	user_display_name: string
	user_id: string
	points_used: number
	points_win?: number
}

export interface PredictionOutcome {
	id: string
	title: string
	color: string
	users: number
	channel_points: number
	top_predictors: PredictionTopPredictor[]
}

export interface PredictionInfo {
	title: string
	outcomes: PredictionOutcome[]
}

export interface PredictionBeginMessage {
	base_info: BaseInfo
	user_name: string
	user_display_name: string
	info: PredictionInfo
}

export interface PredictionProgressMessage {
	base_info: BaseInfo
	user_name: string
	user_display_name: string
	info: PredictionInfo
}

export interface PredictionLockMessage {
	base_info: BaseInfo
	user_name: string
	user_display_name: string
	info: PredictionInfo
}

export interface PredictionEndMessage {
	base_info: BaseInfo
	user_name: string
	user_display_name: string
	info: PredictionInfo
	winning_outcome_id: string
}

export interface StreamFirstUserJoinMessage {
	base_info: BaseInfo
	user_name: string
}

export interface ChannelBanMessage {
	base_info: BaseInfo
	user_name: string
	user_login: string
	broadcaster_user_name: string
	broadcaster_user_login: string
	moderator_user_name: string
	moderator_user_login: string
	reason: string
	ends_at: string
	is_permanent: boolean
}

export interface ChannelUnbanRequestCreateMessage {
	base_info: BaseInfo
	user_name: string
	user_login: string
	broadcaster_user_name: string
	broadcaster_user_login: string
	text: string
}

export interface ChannelUnbanRequestResolveMessage {
	base_info: BaseInfo
	user_name: string
	user_login: string
	broadcaster_user_name: string
	broadcaster_user_login: string
	moderator_user_name: string
	moderator_user_login: string
}

export interface ChannelMessageDeleteMessage {
	base_info: BaseInfo
	user_id: string
	user_name: string
	user_login: string
	broadcaster_user_name: string
	broadcaster_user_login: string
	message_id: string
}

export interface ModeratorAddedMessage {
	base_info: BaseInfo
	user_id: string
	user_name: string
}

export interface ModeratorRemovedMessage {
	base_info: BaseInfo
	user_id: string
	user_name: string
}

export interface VipAddedMessage {
	base_info: BaseInfo
	user_id: string
	user_name: string
}

export interface VipRemovedMessage {
	base_info: BaseInfo
	user_id: string
	user_name: string
}
