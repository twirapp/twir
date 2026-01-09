export default defineNuxtConfig({
	imports: {
		dirs: ['composables', 'utils/**', 'stores'],
	},

	components: [
		{
			path: '#layers/dashboard/components',
			pathPrefix: false,
			global: true,
		},
	],
})
