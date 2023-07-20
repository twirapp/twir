export type Integration = {
	id: string,
	enabled: true,
	accessToken: string | null,
	refreshToken: string | null,
	clientId: string | null,
	clientSecret: string | null,
	apiKey: string | null,
	data: Record<string, any> | null,
	channelId: string,
	integrationId: string,
	integration: {
		id: string,
		service: string,
		accessToken: string | null,
		refreshToken: string | null,
		clientId: string | null,
		clientSecret: string | null,
		apiKey: string | null,
		redirectUrl: string | null
	},
	channel: {
		id: string,
		isEnabled: boolean,
		isTwitchBanned: boolean,
		isBanned: boolean,
		botId: string,
		isBotMod: boolean
	}
}
