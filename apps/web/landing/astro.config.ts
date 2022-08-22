import vue from '@astrojs/vue';
import { defineConfig } from 'astro/config';

import tailwind from './integrations/tailwind';

// https://astro.build/config
export default defineConfig({
  integrations: [vue(), tailwind()],
});
