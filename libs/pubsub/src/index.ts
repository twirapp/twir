import Redis from 'redis';

export const createPubSub = async (url: string) => {
  const subscriber = Redis.createClient({
    url,
  });
  const publisher = Redis.createClient({
    url,
  });
  await subscriber.connect();
  await publisher.connect();

  return {
    publish: (topic: string, data: Record<any, any> | any[] | string) => {
      publisher.publish(topic, JSON.stringify(data));
    },
    subscribe: (topic: string, cb: (data: string) => void) => {
      subscriber.subscribe(topic, (msg) => cb(msg));
    },
  };
};
