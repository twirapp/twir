import './main.css'
import './assets/index.css'
import { install as VueMonacoEditorPlugin } from '@guolao/vue-monaco-editor'
import { broadcastQueryClient } from '@tanstack/query-broadcast-client-experimental'
import { QueryClient, VueQueryPlugin } from '@tanstack/vue-query'
import { createApp } from 'vue'

import App from '@/App.vue'

import { i18n } from './plugins/i18n.js'
import { newRouter } from './plugins/router.js'

const app = createApp(App)

const queryClient = new QueryClient({
	defaultOptions: {
		queries: {
			refetchOnWindowFocus: false,
			refetchOnMount: false,
			refetchOnReconnect: false,
			retry: false,
		},
	},
})

broadcastQueryClient({
	queryClient,
	broadcastChannel: 'twir-dashboard',
})

VueQueryPlugin.install(app, {
	queryClient,
})

app
	.use(i18n)
	.use(VueMonacoEditorPlugin, {
		paths: {
			// The recommended CDN config
			vs: 'https://cdn.jsdelivr.net/npm/monaco-editor@0.52.0/min/vs',
		},
	})
	.use(newRouter())
	.mount('#app')

if (import.meta.env.DEV) {
	document.title = 'Twir (dev)'
}
