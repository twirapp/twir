type SubscribeCallback = (topic: string, callback: (data: string) => void | Promise<void>) => void;
type PublishData = (topic: string, data: Record<any, any> | any[] | string) => void;
