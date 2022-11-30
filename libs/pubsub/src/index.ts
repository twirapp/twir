import Redis from 'redis';

export const createPubSub = async () => {
  const subscriber = Redis.createClient();
  const publisher = Redis.createClient();

  return {
    publish: (topic: string, data: Record<any, any> | any[]) => {
      publisher.publish(topic, JSON.stringify(data));
    },
    subscribe: (topic: string, cb: (data: string) => void) => {
      subscriber.subscribe(topic, cb);
    },
  };
};
