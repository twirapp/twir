import Redis from 'redis';

export const createPubSub = async (url) => {
  const subscriber = Redis.createClient({
    url,
  });
  const publisher = Redis.createClient({
    url,
  });
  await subscriber.connect();
  await publisher.connect();

  return {
    publish: (topic, data) => {
      publisher.publish(topic, JSON.stringify(data));
    },
    subscribe: (topic, cb) => {
      subscriber.subscribe(topic, cb);
    },
  };
};
