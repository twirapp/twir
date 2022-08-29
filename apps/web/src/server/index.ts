import compression from 'compression';
import express from 'express';
import { renderPage } from 'vite-plugin-ssr';

const isProduction = 'production';
const root = `${__dirname}/../..`;

console.log(root);

startServer();

async function startServer() {
  const app = express();

  app.use(compression());

  if (isProduction) {
    console.log('production');
    const sirv = require('sirv');
    app.use(sirv(`${root}/dist/client`));
  } else {
    const vite = require('vite');
    const viteDevMiddleware = (
      await vite.createServer({
        root,
        server: { middlewareMode: true },
      })
    ).middlewares;
    app.use(viteDevMiddleware);
  }

  app.get('*', async (req, res, next) => {
    const pageContextInit = {
      urlOriginal: req.originalUrl,
    };
    const pageContext = await renderPage(pageContextInit);
    const { httpResponse } = pageContext;
    if (!httpResponse) return next();
    const { statusCode, contentType } = httpResponse;
    const body = await httpResponse.getBody();
    res.status(statusCode).type(contentType).send(body);
  });

  const port = process.env.PORT || 3000;
  app.listen(port);
  console.log(`Server running at http://localhost:${port}`);
}
