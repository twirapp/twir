import Express from 'express';

import { initChannels } from './libs/initChannels.js';
import { eventSubMiddleware } from './libs/middleware.js';

const app = Express();
await eventSubMiddleware.apply(app);

app.listen(3003, async () => {
  await eventSubMiddleware.markAsReady();
  await initChannels();
});
