import * as EventSub from '@tsuwari/grpc/generated/eventsub/eventsub';
import { PORTS } from '@tsuwari/grpc/servers/constants';
import Express from 'express';
import Ngrok from 'ngrok';
import { createServer } from 'nice-grpc';

import { initChannels } from './libs/initChannels.js';
import { eventSubMiddleware, subscribeToEvents } from './libs/middleware.js';

const app = Express();
app.get('/', (req, res) => {
  res.send('Twir eventsub home.');
});
await eventSubMiddleware.apply(app);

app.listen(3003, async () => {
  await eventSubMiddleware.markAsReady();
  await initChannels();
});

const eventSubService: EventSub.EventSubServiceImplementation = {
  async subscribeToEvents(request: EventSub.SubscribeToEventsRequest) {
    subscribeToEvents(request.channelId);
    return {};
  },
};

const server = createServer({
  'grpc.keepalive_time_ms': 1 * 60 * 1000,
});

server.add(EventSub.EventSubDefinition, eventSubService);

await server.listen(`0.0.0.0:${PORTS.EVENTSUB_SERVER_PORT}`);

async function close() {
  server.forceShutdown();
  Ngrok.disconnect().catch();
}

process.on('SIGTERM', close).on('SIGINT', close);
