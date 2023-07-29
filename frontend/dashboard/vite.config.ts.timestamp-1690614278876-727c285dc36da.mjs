// vite.config.ts
import path from "node:path";
import { fileURLToPath } from "node:url";
import VueI18nPlugin from "file:///home/satont/Projects/twir/node_modules/.pnpm/@intlify+unplugin-vue-i18n@0.12.2_rollup@2.79.1_vue-i18n@9.2.2/node_modules/@intlify/unplugin-vue-i18n/lib/vite.mjs";
import { webUpdateNotice } from "file:///home/satont/Projects/twir/node_modules/.pnpm/@plugin-web-update-notification+vite@1.6.4_vite@4.4.3/node_modules/@plugin-web-update-notification/vite/dist/index.js";
import vue from "file:///home/satont/Projects/twir/node_modules/.pnpm/@vitejs+plugin-vue@4.2.3_vite@4.4.3_vue@3.3.4/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import { defineConfig, loadEnv } from "file:///home/satont/Projects/twir/node_modules/.pnpm/vite@4.4.3_@types+node@20.4.2_sass@1.63.6/node_modules/vite/dist/node/index.js";
import { VitePWA } from "file:///home/satont/Projects/twir/node_modules/.pnpm/vite-plugin-pwa@0.16.4_vite@4.4.3_workbox-build@7.0.0_workbox-window@7.0.0/node_modules/vite-plugin-pwa/dist/index.js";
import svg from "file:///home/satont/Projects/twir/node_modules/.pnpm/vite-svg-loader@4.0.0/node_modules/vite-svg-loader/index.js";
var __vite_injected_original_dirname = "/home/satont/Projects/twir/frontend/dashboard";
var __vite_injected_original_import_meta_url = "file:///home/satont/Projects/twir/frontend/dashboard/vite.config.ts";
var vite_config_default = defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "");
  return {
    plugins: [
      vue({
        script: {
          defineModel: true
        }
      }),
      svg({ svgo: false }),
      VitePWA(),
      webUpdateNotice({
        notificationProps: {
          title: "New version",
          description: "An update available, please refresh the page to get latest features and bug fixes!",
          buttonText: "refresh",
          dismissButtonText: "cancel"
        },
        checkInterval: 1 * 60 * 1e3
      }),
      VueI18nPlugin({
        include: [
          path.resolve(__vite_injected_original_dirname, "./src/locales/**")
        ],
        strictMessage: false,
        escapeHtml: false
      })
    ],
    base: "/dashboard",
    resolve: {
      alias: {
        vue: "vue/dist/vue.esm-bundler.js",
        "@": fileURLToPath(new URL("./src", __vite_injected_original_import_meta_url))
      }
    },
    server: {
      port: 3006,
      host: true,
      hmr: {
        protocol: env.USE_WSS === "true" ? "wss" : "ws"
      },
      proxy: {
        "/api": {
          target: "http://127.0.0.1:3002",
          changeOrigin: true,
          rewrite: (path2) => path2.replace(/^\/api/, ""),
          ws: true
        },
        "/socket": {
          target: "http://127.0.0.1:3004",
          changeOrigin: true,
          ws: true,
          rewrite: (path2) => path2.replace(/^\/socket/, "")
        },
        "/p": {
          target: "http://127.0.0.1:3007",
          changeOrigin: true,
          ws: true
          // rewrite: (path) => path.replace(/^\/p/, ''),
        },
        "/overlays": {
          target: "http://127.0.0.1:3008",
          changeOrigin: true,
          ws: true
        }
      }
    },
    clearScreen: false
  };
});
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCIvaG9tZS9zYXRvbnQvUHJvamVjdHMvdHdpci9mcm9udGVuZC9kYXNoYm9hcmRcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZmlsZW5hbWUgPSBcIi9ob21lL3NhdG9udC9Qcm9qZWN0cy90d2lyL2Zyb250ZW5kL2Rhc2hib2FyZC92aXRlLmNvbmZpZy50c1wiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9pbXBvcnRfbWV0YV91cmwgPSBcImZpbGU6Ly8vaG9tZS9zYXRvbnQvUHJvamVjdHMvdHdpci9mcm9udGVuZC9kYXNoYm9hcmQvdml0ZS5jb25maWcudHNcIjtpbXBvcnQgcGF0aCBmcm9tICdub2RlOnBhdGgnO1xuaW1wb3J0IHsgZmlsZVVSTFRvUGF0aCB9IGZyb20gJ25vZGU6dXJsJztcblxuaW1wb3J0IFZ1ZUkxOG5QbHVnaW4gZnJvbSAnQGludGxpZnkvdW5wbHVnaW4tdnVlLWkxOG4vdml0ZSc7XG5pbXBvcnQgeyB3ZWJVcGRhdGVOb3RpY2UgfSBmcm9tICdAcGx1Z2luLXdlYi11cGRhdGUtbm90aWZpY2F0aW9uL3ZpdGUnO1xuaW1wb3J0IHZ1ZSBmcm9tICdAdml0ZWpzL3BsdWdpbi12dWUnO1xuaW1wb3J0IHsgZGVmaW5lQ29uZmlnLCBsb2FkRW52IH0gZnJvbSAndml0ZSc7XG5pbXBvcnQgeyBWaXRlUFdBIH0gZnJvbSAndml0ZS1wbHVnaW4tcHdhJztcbmltcG9ydCBzdmcgZnJvbSAndml0ZS1zdmctbG9hZGVyJztcblxuLy8gaHR0cHM6Ly92aXRlanMuZGV2L2NvbmZpZy9cbmV4cG9ydCBkZWZhdWx0IGRlZmluZUNvbmZpZygoeyBtb2RlIH0pID0+IHtcblx0Y29uc3QgZW52ID0gbG9hZEVudihtb2RlLCBwcm9jZXNzLmN3ZCgpLCAnJyk7XG5cblx0cmV0dXJuIHtcblx0XHRwbHVnaW5zOiBbXG5cdFx0XHR2dWUoe1xuXHRcdFx0XHRzY3JpcHQ6IHtcblx0XHRcdFx0XHRkZWZpbmVNb2RlbDogdHJ1ZSxcblx0XHRcdFx0fSxcblx0XHRcdH0pLFxuXHRcdFx0c3ZnKHsgc3ZnbzogZmFsc2UgfSksXG5cdFx0XHRWaXRlUFdBKCksXG5cdFx0XHR3ZWJVcGRhdGVOb3RpY2Uoe1xuXHRcdFx0XHRub3RpZmljYXRpb25Qcm9wczoge1xuXHRcdFx0XHRcdHRpdGxlOiAnTmV3IHZlcnNpb24nLFxuXHRcdFx0XHRcdGRlc2NyaXB0aW9uOiAnQW4gdXBkYXRlIGF2YWlsYWJsZSwgcGxlYXNlIHJlZnJlc2ggdGhlIHBhZ2UgdG8gZ2V0IGxhdGVzdCBmZWF0dXJlcyBhbmQgYnVnIGZpeGVzIScsXG5cdFx0XHRcdFx0YnV0dG9uVGV4dDogJ3JlZnJlc2gnLFxuXHRcdFx0XHRcdGRpc21pc3NCdXR0b25UZXh0OiAnY2FuY2VsJyxcblx0XHRcdFx0fSxcblx0XHRcdFx0Y2hlY2tJbnRlcnZhbDogMSAqIDYwICogMTAwMCxcblx0XHRcdH0pLFxuXHRcdFx0VnVlSTE4blBsdWdpbih7XG5cdFx0XHRcdGluY2x1ZGU6IFtcblx0XHRcdFx0XHRwYXRoLnJlc29sdmUoX19kaXJuYW1lLCAnLi9zcmMvbG9jYWxlcy8qKicpLFxuXHRcdFx0XHRdLFxuXHRcdFx0XHRzdHJpY3RNZXNzYWdlOiBmYWxzZSxcblx0XHRcdFx0ZXNjYXBlSHRtbDogZmFsc2UsXG5cdFx0XHR9KSxcblx0XHRdLFxuXHRcdGJhc2U6ICcvZGFzaGJvYXJkJyxcblx0XHRyZXNvbHZlOiB7XG5cdFx0XHRhbGlhczoge1xuXHRcdFx0XHR2dWU6ICd2dWUvZGlzdC92dWUuZXNtLWJ1bmRsZXIuanMnLFxuXHRcdFx0XHQnQCc6IGZpbGVVUkxUb1BhdGgobmV3IFVSTCgnLi9zcmMnLCBpbXBvcnQubWV0YS51cmwpKSxcblx0XHRcdH0sXG5cdFx0fSxcblx0XHRzZXJ2ZXI6IHtcblx0XHRcdHBvcnQ6IDMwMDYsXG5cdFx0XHRob3N0OiB0cnVlLFxuXHRcdFx0aG1yOiB7XG5cdFx0XHRcdHByb3RvY29sOiBlbnYuVVNFX1dTUyA9PT0gJ3RydWUnID8gJ3dzcycgOiAnd3MnLFxuXHRcdFx0fSxcblx0XHRcdHByb3h5OiB7XG5cdFx0XHRcdCcvYXBpJzoge1xuXHRcdFx0XHRcdHRhcmdldDogJ2h0dHA6Ly8xMjcuMC4wLjE6MzAwMicsXG5cdFx0XHRcdFx0Y2hhbmdlT3JpZ2luOiB0cnVlLFxuXHRcdFx0XHRcdHJld3JpdGU6IChwYXRoKSA9PiBwYXRoLnJlcGxhY2UoL15cXC9hcGkvLCAnJyksXG5cdFx0XHRcdFx0d3M6IHRydWUsXG5cdFx0XHRcdH0sXG5cdFx0XHRcdCcvc29ja2V0Jzoge1xuXHRcdFx0XHRcdHRhcmdldDogJ2h0dHA6Ly8xMjcuMC4wLjE6MzAwNCcsXG5cdFx0XHRcdFx0Y2hhbmdlT3JpZ2luOiB0cnVlLFxuXHRcdFx0XHRcdHdzOiB0cnVlLFxuXHRcdFx0XHRcdHJld3JpdGU6IChwYXRoKSA9PiBwYXRoLnJlcGxhY2UoL15cXC9zb2NrZXQvLCAnJyksXG5cdFx0XHRcdH0sXG5cdFx0XHRcdCcvcCc6IHtcblx0XHRcdFx0XHR0YXJnZXQ6ICdodHRwOi8vMTI3LjAuMC4xOjMwMDcnLFxuXHRcdFx0XHRcdGNoYW5nZU9yaWdpbjogdHJ1ZSxcblx0XHRcdFx0XHR3czogdHJ1ZSxcblx0XHRcdFx0XHQvLyByZXdyaXRlOiAocGF0aCkgPT4gcGF0aC5yZXBsYWNlKC9eXFwvcC8sICcnKSxcblx0XHRcdFx0fSxcblx0XHRcdFx0Jy9vdmVybGF5cyc6IHtcblx0XHRcdFx0XHR0YXJnZXQ6ICdodHRwOi8vMTI3LjAuMC4xOjMwMDgnLFxuXHRcdFx0XHRcdGNoYW5nZU9yaWdpbjogdHJ1ZSxcblx0XHRcdFx0XHR3czogdHJ1ZSxcblx0XHRcdFx0fSxcblx0XHRcdH0sXG5cdFx0fSxcblx0XHRjbGVhclNjcmVlbjogZmFsc2UsXG5cdH07XG59KTtcbiJdLAogICJtYXBwaW5ncyI6ICI7QUFBeVQsT0FBTyxVQUFVO0FBQzFVLFNBQVMscUJBQXFCO0FBRTlCLE9BQU8sbUJBQW1CO0FBQzFCLFNBQVMsdUJBQXVCO0FBQ2hDLE9BQU8sU0FBUztBQUNoQixTQUFTLGNBQWMsZUFBZTtBQUN0QyxTQUFTLGVBQWU7QUFDeEIsT0FBTyxTQUFTO0FBUmhCLElBQU0sbUNBQW1DO0FBQXlKLElBQU0sMkNBQTJDO0FBV25QLElBQU8sc0JBQVEsYUFBYSxDQUFDLEVBQUUsS0FBSyxNQUFNO0FBQ3pDLFFBQU0sTUFBTSxRQUFRLE1BQU0sUUFBUSxJQUFJLEdBQUcsRUFBRTtBQUUzQyxTQUFPO0FBQUEsSUFDTixTQUFTO0FBQUEsTUFDUixJQUFJO0FBQUEsUUFDSCxRQUFRO0FBQUEsVUFDUCxhQUFhO0FBQUEsUUFDZDtBQUFBLE1BQ0QsQ0FBQztBQUFBLE1BQ0QsSUFBSSxFQUFFLE1BQU0sTUFBTSxDQUFDO0FBQUEsTUFDbkIsUUFBUTtBQUFBLE1BQ1IsZ0JBQWdCO0FBQUEsUUFDZixtQkFBbUI7QUFBQSxVQUNsQixPQUFPO0FBQUEsVUFDUCxhQUFhO0FBQUEsVUFDYixZQUFZO0FBQUEsVUFDWixtQkFBbUI7QUFBQSxRQUNwQjtBQUFBLFFBQ0EsZUFBZSxJQUFJLEtBQUs7QUFBQSxNQUN6QixDQUFDO0FBQUEsTUFDRCxjQUFjO0FBQUEsUUFDYixTQUFTO0FBQUEsVUFDUixLQUFLLFFBQVEsa0NBQVcsa0JBQWtCO0FBQUEsUUFDM0M7QUFBQSxRQUNBLGVBQWU7QUFBQSxRQUNmLFlBQVk7QUFBQSxNQUNiLENBQUM7QUFBQSxJQUNGO0FBQUEsSUFDQSxNQUFNO0FBQUEsSUFDTixTQUFTO0FBQUEsTUFDUixPQUFPO0FBQUEsUUFDTixLQUFLO0FBQUEsUUFDTCxLQUFLLGNBQWMsSUFBSSxJQUFJLFNBQVMsd0NBQWUsQ0FBQztBQUFBLE1BQ3JEO0FBQUEsSUFDRDtBQUFBLElBQ0EsUUFBUTtBQUFBLE1BQ1AsTUFBTTtBQUFBLE1BQ04sTUFBTTtBQUFBLE1BQ04sS0FBSztBQUFBLFFBQ0osVUFBVSxJQUFJLFlBQVksU0FBUyxRQUFRO0FBQUEsTUFDNUM7QUFBQSxNQUNBLE9BQU87QUFBQSxRQUNOLFFBQVE7QUFBQSxVQUNQLFFBQVE7QUFBQSxVQUNSLGNBQWM7QUFBQSxVQUNkLFNBQVMsQ0FBQ0EsVUFBU0EsTUFBSyxRQUFRLFVBQVUsRUFBRTtBQUFBLFVBQzVDLElBQUk7QUFBQSxRQUNMO0FBQUEsUUFDQSxXQUFXO0FBQUEsVUFDVixRQUFRO0FBQUEsVUFDUixjQUFjO0FBQUEsVUFDZCxJQUFJO0FBQUEsVUFDSixTQUFTLENBQUNBLFVBQVNBLE1BQUssUUFBUSxhQUFhLEVBQUU7QUFBQSxRQUNoRDtBQUFBLFFBQ0EsTUFBTTtBQUFBLFVBQ0wsUUFBUTtBQUFBLFVBQ1IsY0FBYztBQUFBLFVBQ2QsSUFBSTtBQUFBO0FBQUEsUUFFTDtBQUFBLFFBQ0EsYUFBYTtBQUFBLFVBQ1osUUFBUTtBQUFBLFVBQ1IsY0FBYztBQUFBLFVBQ2QsSUFBSTtBQUFBLFFBQ0w7QUFBQSxNQUNEO0FBQUEsSUFDRDtBQUFBLElBQ0EsYUFBYTtBQUFBLEVBQ2Q7QUFDRCxDQUFDOyIsCiAgIm5hbWVzIjogWyJwYXRoIl0KfQo=
