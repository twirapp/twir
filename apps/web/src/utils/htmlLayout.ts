import type { Readable } from 'stream';

import { escapeInject } from 'vite-plugin-ssr';

import svgFavicon from '@/assets/NewLogo.svg';
import { author, ogImage } from '@/data/seo.js';
import type { Locale } from '@/locales';

export const htmlLayout = (data: {
  title: string;
  description: string;
  keywords: string[];
  content?: Readable;
  urlCanonical?: string;
  urlOriginal: string;
  locale: Locale;
}) => escapeInject`<!DOCTYPE html>
    <html lang="${data.locale}">
      <head>
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <meta charset="utf-8" />
        ${data.urlCanonical ? escapeInject`<link rel="canonical" href="${data.urlCanonical}">` : ''}
        <link rel="icon" href="${svgFavicon}" sizes="any" type="image/svg+xml">
        
        <title>${data.title}</title>
        <meta name="description" content="${data.description}" >
        <meta name="keywords" content="${data.keywords.join(', ')}">
        <meta name="author" content="${author}">

        <meta property="og:title" content="${data.title}" >
        <meta property="og:url" content="${data.urlOriginal}" />
        <meta property="og:description" content="${data.description}" />
        <meta property="og:image" content="${ogImage}" />
      </head>
      <body> 
        <div id="app">${data.content || ''}</div>
      </body>
    </html>`;
