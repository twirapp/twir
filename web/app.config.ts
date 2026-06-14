export default defineAppConfig({
	icon: {
		customize: (content: string, _name?: string, prefix?: string) => {
			if (prefix?.startsWith('twir-')) {
				return content
			}

			return content
				.replace(/stroke-width="[^"]*"/g, 'stroke-width="1.5"')
				.replace(/fill="[^"]*"/g, 'fill="currentColor"')
		},
	},
})
