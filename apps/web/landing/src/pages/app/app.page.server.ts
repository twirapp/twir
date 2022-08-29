import { escapeInject } from 'vite-plugin-ssr';

import type { PageContext } from '@/types/pageContext.js';

export { render };

async function render(pageContext: PageContext) {
  return escapeInject`<!DOCTYPE html>
    <html>
      <head>
        <title>App</title>
      </head>
      <body>
        <div id="app">
        </div>
      </body>
    </html>`;
}
