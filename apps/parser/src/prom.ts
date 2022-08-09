import client from 'prom-client';

client.collectDefaultMetrics();

export const parseResponseCounter = new client.Counter({
  name: 'parse_response',
  help: 'Parse responses requests',
});

export const parseChatMessageCounter = new client.Counter({
  name: 'parse_chat_message',
  help: 'Parse chat message requests',
});

export const prometheus = client;
