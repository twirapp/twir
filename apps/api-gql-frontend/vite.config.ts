import vue from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
	server: {
		proxy: {
			'/query': {
				target: 'http://localhost:3009',
				changeOrigin: true,
				ws: true,
			},
		},
	},
});
