import { resolve } from 'node:path';

import fastify from 'fastify';
import { renderPage } from 'vite-plugin-ssr';

import type { PageContext } from '@/types/pageContext.js';

const PORT = Number(process.env.PORT) || 3000;
const isProduction = process.env.NODE_ENV === 'production' || false;
const root = resolve(__dirname, '../..');

startServer();

async function startServer() {
  try {
    const app = fastify({
      logger: isProduction ? true : false,
    });

    await app.register(import('@fastify/middie'));
    await app.register(import('@fastify/compress'), { global: false });

    if (isProduction) {
      await app.register(import('@fastify/static'), {
        root: `${root}/dist/client`,
      });
    } else {
      const { createServer } = await import('vite');

      const viteDevMiddleware = (
        await createServer({
          root,
          server: { middlewareMode: true },
        })
      ).middlewares;

      app.use(viteDevMiddleware);
    }

    app.get(isProduction ? '/app/*' : '*', async (req, res) => {
      const urlOriginal = `${req.protocol}://${req.hostname + req.url}`;

      const pageContextInit: Partial<PageContext> = {
        urlOriginal,
      };

      const { httpResponse } = await renderPage(pageContextInit);
      if (!httpResponse) return;

      const { statusCode, contentType } = httpResponse;
      const body = await httpResponse.getBody();

      res.status(statusCode).type(contentType).send(body);
    });

    app.listen({ port: PORT });
    console.log(`Server running at http://localhost:${PORT}`);
  } catch (error) {
    console.error(error);
    process.exit(1);
  }
}
