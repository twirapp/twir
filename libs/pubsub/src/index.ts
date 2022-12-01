import Redis from 'redis';

export const createPubSub = async (url: string) => {
  const subscriber = Redis.createClient({
    url,
  });
  const publisher = Redis.createClient({
    url,
  });

  return {
    publish: (topic: string, data: Record<any, any> | any[]) => {
      publisher.publish(topic, JSON.stringify(data));
    },
    subscribe: (topic: string, cb: (data: string) => void) => {
      subscriber.subscribe(topic, cb);
    },
  };
};
