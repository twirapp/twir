/**
 * @see https://prettier.io/docs/configuration
 * @type {import("prettier").Config}
 */
const config = {
	trailingComma: 'es5',
	tabWidth: 2,
	semi: false,
	singleQuote: true,
	plugins: ['@prettier/plugin-oxc'],
	printWidth: 100,
	useTabs: true,
}

export default config
