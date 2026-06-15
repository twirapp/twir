import type { ChatOverlay } from '@/gql/graphql'

import { ChatOverlayAnimation } from '@/gql/graphql'

export type ChatSettingsWithOptionalId = Omit<ChatOverlay, 'id'> & { id?: string }

export const defaultChatSettings: ChatSettingsWithOptionalId = {
	fontFamily: 'inter',
	fontSize: 20,
	fontWeight: 400,
	fontStyle: 'normal',
	hideBots: false,
	hideCommands: false,
	messageHideTimeout: 0,
	messageShowDelay: 0,
	preset: 'clean',
	showBadges: true,
	showAnnounceBadge: true,
	textShadowColor: 'rgba(0,0,0,1)',
	textShadowSize: 0,
	chatBackgroundColor: 'rgba(0, 0, 0, 0)',
	direction: 'top',
	paddingContainer: 0,
	animation: ChatOverlayAnimation.Default,
}
