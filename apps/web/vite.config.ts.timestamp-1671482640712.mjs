// vite.config.ts
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import vue from "file:///workspace/node_modules/.pnpm/@vitejs+plugin-vue@3.2.0_vite@3.2.4+vue@3.2.45/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import { defineConfig } from "file:///workspace/node_modules/.pnpm/vite@3.2.4/node_modules/vite/dist/node/index.js";
import ssr from "file:///workspace/node_modules/.pnpm/vite-plugin-ssr@0.4.54_vite@3.2.4/node_modules/vite-plugin-ssr/dist/cjs/node/plugin/index.js";
import svg from "file:///workspace/node_modules/.pnpm/vite-svg-loader@3.6.0/node_modules/vite-svg-loader/index.js";
var __vite_injected_original_import_meta_url = "file:///workspace/apps/web/vite.config.ts";
var __dirname = dirname(fileURLToPath(__vite_injected_original_import_meta_url));
var vite_config_default = defineConfig({
  define: {
    __VUE_I18N_FULL_INSTALL__: false,
    __VUE_I18N_LEGACY_API__: false,
    __INTLIFY_PROD_DEVTOOLS__: false
  },
  plugins: [
    vue(),
    svg({
      svgo: false,
      defaultImport: "url"
    }),
    ssr({
      prerender: true
    })
  ],
  resolve: {
    alias: {
      "@": resolve(__dirname, "src")
    }
  },
  server: {
    host: true,
    port: Number(process.env.VITE_PORT ?? 3005),
    proxy: {
      "/api": {
        target: "http://127.0.0.1:3002",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ""),
        ws: true
      },
      "/dashboard": {
        target: "http://127.0.0.1:3006/dashboard",
        changeOrigin: false,
        ws: true,
        rewrite: (path) => path.replace(/^\/dashboard/, "")
      }
    }
  }
});
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCIvd29ya3NwYWNlL2FwcHMvd2ViXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ZpbGVuYW1lID0gXCIvd29ya3NwYWNlL2FwcHMvd2ViL3ZpdGUuY29uZmlnLnRzXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ltcG9ydF9tZXRhX3VybCA9IFwiZmlsZTovLy93b3Jrc3BhY2UvYXBwcy93ZWIvdml0ZS5jb25maWcudHNcIjtpbXBvcnQgeyBkaXJuYW1lLCByZXNvbHZlIH0gZnJvbSAnbm9kZTpwYXRoJztcbmltcG9ydCB7IGZpbGVVUkxUb1BhdGggfSBmcm9tICdub2RlOnVybCc7XG5cbmltcG9ydCB2dWUgZnJvbSAnQHZpdGVqcy9wbHVnaW4tdnVlJztcbmltcG9ydCB7IGRlZmluZUNvbmZpZyB9IGZyb20gJ3ZpdGUnO1xuaW1wb3J0IHNzciBmcm9tICd2aXRlLXBsdWdpbi1zc3IvcGx1Z2luJztcbmltcG9ydCBzdmcgZnJvbSAndml0ZS1zdmctbG9hZGVyJztcblxuY29uc3QgX19kaXJuYW1lID0gZGlybmFtZShmaWxlVVJMVG9QYXRoKGltcG9ydC5tZXRhLnVybCkpO1xuXG5leHBvcnQgZGVmYXVsdCBkZWZpbmVDb25maWcoe1xuICBkZWZpbmU6IHtcbiAgICBfX1ZVRV9JMThOX0ZVTExfSU5TVEFMTF9fOiBmYWxzZSxcbiAgICBfX1ZVRV9JMThOX0xFR0FDWV9BUElfXzogZmFsc2UsXG4gICAgX19JTlRMSUZZX1BST0RfREVWVE9PTFNfXzogZmFsc2UsXG4gIH0sXG4gIHBsdWdpbnM6IFtcbiAgICB2dWUoKSxcbiAgICBzdmcoe1xuICAgICAgc3ZnbzogZmFsc2UsXG4gICAgICBkZWZhdWx0SW1wb3J0OiAndXJsJyxcbiAgICB9KSxcbiAgICBzc3Ioe1xuICAgICAgcHJlcmVuZGVyOiB0cnVlLFxuICAgIH0pLFxuICBdLFxuICByZXNvbHZlOiB7XG4gICAgYWxpYXM6IHtcbiAgICAgICdAJzogcmVzb2x2ZShfX2Rpcm5hbWUsICdzcmMnKSxcbiAgICB9LFxuICB9LFxuICAvLyBiYXNlOiAnZGFzaGJvYXJkJyxcbiAgc2VydmVyOiB7XG4gICAgaG9zdDogdHJ1ZSxcbiAgICBwb3J0OiBOdW1iZXIocHJvY2Vzcy5lbnYuVklURV9QT1JUID8/IDMwMDUpLFxuICAgIHByb3h5OiB7XG4gICAgICAnL2FwaSc6IHtcbiAgICAgICAgdGFyZ2V0OiAnaHR0cDovLzEyNy4wLjAuMTozMDAyJyxcbiAgICAgICAgY2hhbmdlT3JpZ2luOiB0cnVlLFxuICAgICAgICByZXdyaXRlOiAocGF0aCkgPT4gcGF0aC5yZXBsYWNlKC9eXFwvYXBpLywgJycpLFxuICAgICAgICB3czogdHJ1ZSxcbiAgICAgIH0sXG4gICAgICAnL2Rhc2hib2FyZCc6IHtcbiAgICAgICAgdGFyZ2V0OiAnaHR0cDovLzEyNy4wLjAuMTozMDA2L2Rhc2hib2FyZCcsXG4gICAgICAgIGNoYW5nZU9yaWdpbjogZmFsc2UsXG4gICAgICAgIHdzOiB0cnVlLFxuICAgICAgICByZXdyaXRlOiAocGF0aCkgPT4gcGF0aC5yZXBsYWNlKC9eXFwvZGFzaGJvYXJkLywgJycpLFxuICAgICAgfSxcbiAgICB9LFxuICB9LFxufSk7XG4iXSwKICAibWFwcGluZ3MiOiAiO0FBQTJPLFNBQVMsU0FBUyxlQUFlO0FBQzVRLFNBQVMscUJBQXFCO0FBRTlCLE9BQU8sU0FBUztBQUNoQixTQUFTLG9CQUFvQjtBQUM3QixPQUFPLFNBQVM7QUFDaEIsT0FBTyxTQUFTO0FBTjhILElBQU0sMkNBQTJDO0FBUS9MLElBQU0sWUFBWSxRQUFRLGNBQWMsd0NBQWUsQ0FBQztBQUV4RCxJQUFPLHNCQUFRLGFBQWE7QUFBQSxFQUMxQixRQUFRO0FBQUEsSUFDTiwyQkFBMkI7QUFBQSxJQUMzQix5QkFBeUI7QUFBQSxJQUN6QiwyQkFBMkI7QUFBQSxFQUM3QjtBQUFBLEVBQ0EsU0FBUztBQUFBLElBQ1AsSUFBSTtBQUFBLElBQ0osSUFBSTtBQUFBLE1BQ0YsTUFBTTtBQUFBLE1BQ04sZUFBZTtBQUFBLElBQ2pCLENBQUM7QUFBQSxJQUNELElBQUk7QUFBQSxNQUNGLFdBQVc7QUFBQSxJQUNiLENBQUM7QUFBQSxFQUNIO0FBQUEsRUFDQSxTQUFTO0FBQUEsSUFDUCxPQUFPO0FBQUEsTUFDTCxLQUFLLFFBQVEsV0FBVyxLQUFLO0FBQUEsSUFDL0I7QUFBQSxFQUNGO0FBQUEsRUFFQSxRQUFRO0FBQUEsSUFDTixNQUFNO0FBQUEsSUFDTixNQUFNLE9BQU8sUUFBUSxJQUFJLGFBQWEsSUFBSTtBQUFBLElBQzFDLE9BQU87QUFBQSxNQUNMLFFBQVE7QUFBQSxRQUNOLFFBQVE7QUFBQSxRQUNSLGNBQWM7QUFBQSxRQUNkLFNBQVMsQ0FBQyxTQUFTLEtBQUssUUFBUSxVQUFVLEVBQUU7QUFBQSxRQUM1QyxJQUFJO0FBQUEsTUFDTjtBQUFBLE1BQ0EsY0FBYztBQUFBLFFBQ1osUUFBUTtBQUFBLFFBQ1IsY0FBYztBQUFBLFFBQ2QsSUFBSTtBQUFBLFFBQ0osU0FBUyxDQUFDLFNBQVMsS0FBSyxRQUFRLGdCQUFnQixFQUFFO0FBQUEsTUFDcEQ7QUFBQSxJQUNGO0FBQUEsRUFDRjtBQUNGLENBQUM7IiwKICAibmFtZXMiOiBbXQp9Cg==
