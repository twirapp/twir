import { Server } from 'socket.io';

import { typeorm } from './libs/typeorm.js';
import { authMiddleware } from './middlewares/auth.js';
import { createYoutubeNameSpace } from './namespaces/youtube.js';

const io = new Server();
await typeorm.initialize();

io.use(authMiddleware);
createYoutubeNameSpace(io);

console.info('✅ Started');

io.listen(3004);
