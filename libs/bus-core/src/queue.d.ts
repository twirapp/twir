export type QueueResponse<T> = {
	data: T;
};

export type QueueSubscribeCallback<Req, Res> = (data: Req) => Res | Promise<Res>;

export interface Queue<Req, Res> {
	publish(data: Req): void;
	request(data: Req): Promise<QueueResponse<Res>>;
	subscribeGroup(queueGroup: string, data: QueueSubscribeCallback<Req, Res>): void;
	subscribe(data: QueueSubscribeCallback<Req, Res>): void;
	unsubscribe(): void;
}
