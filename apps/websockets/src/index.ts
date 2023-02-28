import { listen } from './libs/grpc.js';
import { io } from './libs/io.js';
import './namespaces/youtube.js';
import { typeorm } from './libs/typeorm.js';

await typeorm.initialize();
await listen();

console.info('✅ Started');

io.listen(3004);

process
  .on('unhandledRejection', (reason, promise) => {
    console.error('Unhandled Rejection at:', promise, 'reason:', reason);
  })
  .on('uncaughtException', (err) => {
    console.error('Uncaught Exception thrown:', err);
  });