export const EventsubSubscribeAllSubject = 'eventsub.subscribeAll'
export const EventsubUnsubscribeSubject = 'eventsub.unsubscribe'

export interface EventsubSubscribeToAllEventsRequest {
	readonly ChannelID: string
	readonly Platform: string
}

export interface EventsubBindingSnapshot {
	readonly ID: string
	readonly UserID: string
	readonly PlatformChannelID: string
}

export interface EventsubUnsubscribeRequest {
	readonly ChannelID: string
	readonly Platform: string
	readonly Binding?: EventsubBindingSnapshot
}
