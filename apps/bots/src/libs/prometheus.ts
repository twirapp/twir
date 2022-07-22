import client from 'prom-client';

client.collectDefaultMetrics();

export const messagesCounter = new client.Counter({
  name: 'chat_messages',
  help: 'Chat messages sended',
});

export const commandsCounter = new client.Counter({
  name: 'chat_commands',
  help: 'Chat commands usage',
});

export const greetingsCounter = new client.Counter({
  name: 'chat_greetings',
  help: 'Chat greetings',
});

export const keywordsCounter = new client.Counter({
  name: 'chat_keywords',
  help: 'Chat keywords',
});

const responseTimeBuckets = [0.10, 5, 15, 50, 100, 200, 300, 400, 500, 1000, 1500, 2000];

export const messageParseTime = new client.Histogram({
  name: 'message_parse_time',
  help: 'Duration of message parse time in ms',
  buckets: responseTimeBuckets,
});


export const commandsResponseTime = new client.Histogram({
  name: 'commands_response_time',
  help: 'Duration of commands response time in ms',
  buckets: responseTimeBuckets,
  labelNames: ['channel', 'commandName'],
});

export const moderationParseTime = new client.Histogram({
  name: 'moderation_parse_time',
  help: 'Duration of moderation parse',
  buckets: responseTimeBuckets,
});

export const greetingsParseTime = new client.Histogram({
  name: 'greetings_parse_time',
  help: 'Duration of greetings parse',
  buckets: responseTimeBuckets,
});

export const keywordsParseTime = new client.Histogram({
  name: 'keywords_parse_time',
  help: 'Duration of keywords parse',
  buckets: responseTimeBuckets,
});

export const prometheus = client;