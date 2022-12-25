import { listen } from './libs/grpc.js';
import { io } from './libs/io.js';
import './namespaces/youtube.js';
import { typeorm } from './libs/typeorm.js';

await typeorm.initialize();
await listen();

console.info('✅ Started');

io.listen(3004);
