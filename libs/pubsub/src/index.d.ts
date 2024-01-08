declare module '@twir/pubsub' {
	export const createPubSub: (url: string) => Promise<{
		publish: (topic: string, data: Record<any, any> | any[] | string) => void;
		subscribe: (topic: string, callback: (data: string) => void | Promise<void>) => void;
	}>;
}
