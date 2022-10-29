import type { Readable } from 'stream';

import { escapeInject } from 'vite-plugin-ssr';

import svgFavicon from '@/assets/brand/TsuwariInCircle.svg';
import { author, ogImage, SeoPageProps } from '@/data/seo.js';
import type { PageContext } from '@/utils/pageContext.js';

export const htmlLayout = (seo: SeoPageProps, pageContext: PageContext, content?: Readable) => {
  const urlCanonical = pageContext.urlParsed.origin
    ? escapeInject`<link rel="canonical" href="${pageContext.urlParsed.origin}">`
    : '';

  return escapeInject`<!DOCTYPE html>
    <html lang="${pageContext.locale}">
      <head>
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <meta charset="utf-8" />
        ${urlCanonical}
        <link rel="icon" href="${svgFavicon}" sizes="any" type="image/svg+xml">
        
        <title>${seo.title}</title>
        <meta name="description" content="${seo.description}" >
        <meta name="keywords" content="${seo.keywords.join(', ')}">
        <meta name="author" content="${author}">

        <meta property="og:title" content="${seo.title}" >
        <meta property="og:url" content="${pageContext.urlOriginal}" />
        <meta property="og:description" content="${seo.description}" />
        <meta property="og:image" content="${ogImage}" />
      </head>
      <body> 
        <div id="app">${content || ''}</div>
      </body>
    </html>`;
};
