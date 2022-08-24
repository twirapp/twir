import { fileURLToPath } from 'node:url';

import load from '@proload/core';
import type { AstroIntegration } from 'astro';
import tailwindPlugin, { Config as TailwindConfig } from 'tailwindcss';

export default function tailwindIntegration(): AstroIntegration {
  return {
    name: 'tailwind',
    hooks: {
      'astro:config:setup': async ({ config, injectScript }) => {
        const userConfig = await load('tailwind', {
          mustExist: false,
          cwd: fileURLToPath(config.root),
        });

        if (!userConfig?.value) {
          throw new Error(`Could not find a Tailwind config. Does the file exist?`);
        }

        const tailwindConfig = userConfig.value as TailwindConfig;
        config.style.postcss.plugins.push(tailwindPlugin(tailwindConfig));

        injectScript('page-ssr', `import '../src/styles/tailwind.base.css';`);
      },
    },
  };
}
