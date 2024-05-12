declare module '@twir/config' {
	type Config = typeof import('./index.js').config
	const config: Config
	export { config }
}
