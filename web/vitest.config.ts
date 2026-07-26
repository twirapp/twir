import { defineVitestConfig } from '@nuxt/test-utils/config'

export default defineVitestConfig({
	test: {
		environment: 'happy-dom',
		include: ['layers/dashboard/features/**/*.spec.ts', 'layers/dashboard/components/**/*.spec.ts'],
	},
})
