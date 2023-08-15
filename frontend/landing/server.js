import { fileURLToPath } from 'node:url';

import fastifyMiddie from '@fastify/middie';
import fastifyStatic from '@fastify/static';
import Fastify from 'fastify';

import { handler as ssrHandler } from './dist/server/entry.mjs';

// eslint-disable-next-line no-undef
process
	.on('uncaughtException', console.error)
	.on('unhandledRejection', console.error);

const app = Fastify({ logger: true });

await app
	.register(fastifyStatic, {
		root: fileURLToPath(new URL('./dist/client', import.meta.url)),
	})
	.register(fastifyMiddie);
app.use(ssrHandler);


await app.listen({ port: 3005 });
