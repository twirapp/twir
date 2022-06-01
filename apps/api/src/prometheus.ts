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

export const prometheus = client;
