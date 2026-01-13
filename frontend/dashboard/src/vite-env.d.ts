/// <reference types="vite/client" />
/// <reference types="@twirapp/vite-plugin-svg-spritemap/client" />

declare module '*.vue' {
	import type { DefineComponent } from 'vue'
	const component: DefineComponent<object, object, unknown>
	export default component
}
