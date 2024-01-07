import Redis from 'redis';

/**
 *
 * @param url
 * @returns {Promise<{subscribe: *, publish: *}>}
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
		/**
		 *
		 * @param {string} topic
		 * @param {PublishData} data
		 */
    publish: (topic, data) => {
      publisher.publish(topic, JSON.stringify(data));
    },
		/**
		 *
		 * @param {string} topic
		 * @param {SubscribeCallback} cb
		 */
    subscribe: (topic, cb) => {
      subscriber.subscribe(topic, cb);
    },
  };
};
