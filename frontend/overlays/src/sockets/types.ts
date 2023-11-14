export type TwirWebSocketEvent<T = Record<string, any>> = {
	eventName: string,
	data: T,
	createdAt: string
}
