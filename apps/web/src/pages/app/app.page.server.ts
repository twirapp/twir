import { escapeInject } from 'vite-plugin-ssr';

import type { PageContext } from '@/types/pageContext.js';

export { render };

async function render(_pageContext: PageContext) {
  return escapeInject`<!DOCTYPE html>
    <html>
      <head>
        <meta charset="utf-8" />
        <title>App</title>
      </head>
      <body>
        <div id="app">
        </div>
      </body>
    </html>`;
}
