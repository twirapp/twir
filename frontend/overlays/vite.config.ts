import { webUpdateNotice } from '@plugin-web-update-notification/vite';
import vue from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
		vue(),
		webUpdateNotice({
      hiddenDefaultNotification: true,
      checkInterval: 1 * 60 * 1000,
    }),
	],
	base: '/overlays',
  server: {
    host: true,
    port: 3008,
  },
	clearScreen: false,
});
