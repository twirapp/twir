import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

import compress from '@fastify/compress';
import middie from '@fastify/middie';
import { fastify } from 'fastify';
import { PageContextBuiltIn, renderPage } from 'vite-plugin-ssr';

const __dirname = dirname(fileURLToPath(import.meta.url));

const PORT = Number(process.env.PORT) || 3000;
const isProduction = process.env.NODE_ENV === 'production' || false;
const apiProxy = process.env.API_PROXY;
const root = resolve(__dirname, '../..');

async function startServer() {
  try {
    const app = fastify();

    await app.register(middie);
    await app.register(compress, { global: false });

    if (apiProxy) {
      await app.register((await import('@fastify/http-proxy')).default, {
        upstream: apiProxy,
        prefix: '/api',
        http2: false,
      });
    }

    if (isProduction) {
      await app.register((await import('@fastify/static')).default, {
        root: `${root}/dist/client`,
      });
    } else {
      const { createServer } = await import('vite');

      const { middlewares: viteMiddlewares } = await createServer({
        root,
        server: { middlewareMode: true },
      });

      app.use(viteMiddlewares);
    }

    app.get(isProduction ? '/app/*' : '*', async (req, res) => {
      const urlOriginal = `${req.protocol}://${req.hostname + req.url}`;

      const pageContextInit: Partial<PageContextBuiltIn> = {
        urlOriginal,
      };

      const { httpResponse } = await renderPage(pageContextInit);
      if (!httpResponse) return;

      const { statusCode, contentType } = httpResponse;
      const body = await httpResponse.getBody();

      res.status(statusCode).type(contentType).send(body);
    });

    app.listen({ port: PORT, host: '0.0.0.0' });
    await app.ready();

    console.log(`Server running at http://localhost:${PORT}`);
  } catch (error) {
    console.error(error);
    process.exit(1);
  }
}

startServer();
