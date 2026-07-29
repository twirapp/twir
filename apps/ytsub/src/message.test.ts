import { expect, test } from 'bun:test'
import { YTNodes } from 'youtubei.js'

import { isYoutubeSubscriberBadge, normalizeYoutubeTextMessage, toYoutubeTextChatMessage } from './message.ts'

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
			fragments: [{ text: 'hello 😀' }],
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
		emotes: [],
	})
})

test('isYoutubeSubscriberBadge identifies membership badges without treating channel roles as subscriptions', () => {
	expect(isYoutubeSubscriberBadge({ id: 'owner', set_id: 'owner', text: 'Owner' })).toBe(false)
	expect(isYoutubeSubscriberBadge({ id: 'moderator', set_id: 'moderator', text: 'Moderator' })).toBe(false)
	expect(isYoutubeSubscriberBadge({ id: 'member', set_id: 'member', text: 'Member' })).toBe(true)
	expect(isYoutubeSubscriberBadge({ id: '', set_id: 'BADGE_STYLE_TYPE_MEMBER', text: '' })).toBe(true)
	expect(isYoutubeSubscriberBadge({ id: 'verified', set_id: 'BADGE_STYLE_TYPE_VERIFIED', text: '' })).toBe(false)
	expect(isYoutubeSubscriberBadge({ id: '', set_id: '', text: '' })).toBe(false)
})

test('toYoutubeTextChatMessage maps custom-thumbnail membership badges to member', () => {
	const item = {
		id: 'message-1',
		message: { toString: () => 'hello' },
		author: {
			id: 'UCmember',
			name: 'Member',
			badges: [{
				is: (node: unknown) => node === YTNodes.LiveChatAuthorBadge,
				icon_type: undefined,
				style: undefined,
				label: undefined,
				tooltip: 'Member',
				custom_thumbnail: [{ url: 'https://example.com/member-badge.png' }],
			}],
		},
	} as unknown as YTNodes.LiveChatTextMessage

	const message = toYoutubeTextChatMessage(item)

	expect(message.author.badges[0]?.id).toBe('member')
	expect(message.author.badges.some(isYoutubeSubscriberBadge)).toBe(true)
})

test('toYoutubeTextChatMessage maps emoji runs to emote fragments with shortcuts and image urls', () => {
	const item = {
		id: 'message-1',
		message: {
			toString: () => 'hi :face-blue-smiling:',
			runs: [
				{ text: 'hi ' },
				{
					text: ':face-blue-smiling:',
					emoji: {
						emoji_id: 'UCchannel/emoji42',
						shortcuts: [':face-blue-smiling:'],
						image: [{ url: 'https://yt3.ggpht.com/emoji42.png' }],
						is_custom: true,
					},
				},
			],
		},
		author: { id: 'UCviewer', name: 'Viewer', badges: [] },
	} as unknown as YTNodes.LiveChatTextMessage

	const message = toYoutubeTextChatMessage(item)

	expect(message.text).toBe('hi :face-blue-smiling:')
	expect(message.fragments).toEqual([
		{ text: 'hi ' },
		{ text: ':face-blue-smiling:', emote: { id: 'UCchannel/emoji42', url: 'https://yt3.ggpht.com/emoji42.png' } },
	])
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
			fragments: [{ text: 'hello' }],
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
