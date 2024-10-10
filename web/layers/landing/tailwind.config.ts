import type { Config } from 'tailwindcss'

export default <Config>{
	content: ['./**/*.{html,js,astro,vue,ts,tsx}'],
	theme: {
		extend: {
			screens: {
				ss: '360px',
				xs: '480px',
			},
			container: {
				center: true,
				screens: {
					sm: '640px',
					md: '768px',
					lg: '1024px',
					xl: '1280px',
				},
			},
			fontFamily: {
				sans: ['Inter', 'sans-serif'],
			},
		},
	},
}
