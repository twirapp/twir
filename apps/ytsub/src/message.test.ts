import { expect, test } from 'bun:test'

import { isYoutubeSubscriberBadge, normalizeYoutubeTextMessage } from './message.ts'

test('normalizeYoutubeTextMessage emits the generic Go chat schema for text chat', () => {
	const message = normalizeYoutubeTextMessage(
		{
			id: 'binding-1',
			channelId: 'channel-1',
			platformChannelId: 'UCbroadcaster',
			botPlatformId: null,
			userId: 'viewer-user-1',
			enabled: true,
		},
		'Broadcaster',
		{
			id: 'message-1',
			text: 'hello 😀',
			author: {
				id: 'UCviewer',
				name: 'Viewer',
				isModerator: true,
				badges: [{ id: 'moderator', set_id: 'moderator', text: 'Moderator' }],
			},
		}
	)

	expect(message).toEqual({
		message: {
			text: 'hello 😀',
			fragments: [{ text: 'hello 😀', position: { start: 0, end: 7 }, type: 0 }],
		},
		id: 'message-1',
		broadcaster_user_id: 'UCbroadcaster',
		broadcaster_user_name: 'Broadcaster',
		broadcaster_user_login: 'Broadcaster',
		chatter_user_id: 'UCviewer',
		chatter_user_name: 'Viewer',
		chatter_user_login: 'Viewer',
		message_type: 'text',
		platform: 'youtube',
		channel_id: 'channel-1',
		channel_binding_id: 'binding-1',
		user_id: 'viewer-user-1',
		platform_channel_id: 'UCbroadcaster',
		sender_id: 'UCviewer',
		sender_login: 'Viewer',
		sender_display_name: 'Viewer',
		message_id: 'message-1',
		text: 'hello 😀',
		badges: [{ id: 'moderator', set_id: 'moderator', text: 'Moderator' }],
		is_broadcaster: false,
		is_moderator: true,
		is_vip: false,
		is_subscriber: false,
		color: '',
	})
})

test('isYoutubeSubscriberBadge identifies membership badges without treating channel roles as subscriptions', () => {
	expect(isYoutubeSubscriberBadge({ id: 'owner', set_id: 'owner', text: 'Owner' })).toBe(false)
	expect(isYoutubeSubscriberBadge({ id: 'moderator', set_id: 'moderator', text: 'Moderator' })).toBe(false)
	expect(isYoutubeSubscriberBadge({ id: 'member', set_id: 'member', text: 'Member' })).toBe(true)
	expect(isYoutubeSubscriberBadge({ id: '', set_id: '', text: '' })).toBe(false)
})

test('normalizeYoutubeTextMessage marks members as subscribers', () => {
	const message = normalizeYoutubeTextMessage(
		{
			id: 'binding-1',
			channelId: 'channel-1',
			platformChannelId: 'UCbroadcaster',
			botPlatformId: null,
			userId: 'viewer-user-1',
			enabled: true,
		},
		'Broadcaster',
		{
			id: 'message-1',
			text: 'hello',
			author: {
				id: 'UCmember',
				name: 'Member',
				isModerator: false,
				badges: [{ id: 'member', set_id: 'member', text: 'Member' }],
			},
		}
	)

	expect(message.is_subscriber).toBe(true)
})
