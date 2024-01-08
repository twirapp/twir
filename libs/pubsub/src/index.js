import Redis from 'redis';

/**
 *
 * @param url
 * @returns {Promise<{
 * subscribe: (topic: string, callback: (data: string) => void | Promise<void>) => void,
 * publish: (topic: string, data: Record<any, any> | any[] | string) => void;
 * }>}
 */
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
