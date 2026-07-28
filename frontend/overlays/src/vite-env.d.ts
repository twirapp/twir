/// <reference types="vite/client" />

declare global {
	interface Window {
		eruda: any
	}
}

/// <reference types="vite/client" />
/// <reference types="vite-svg-loader" />
declare module '*.vue' {
	import type { DefineComponent } from 'vue'
	const component: DefineComponent<object, object, unknown>
	export default component
}

export {}
