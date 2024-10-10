declare module '@twir/config' {
	type Config = typeof import('./index.js').config
	const config: Config
	function readEnv(path: string)
	export { config, readEnv }
}
