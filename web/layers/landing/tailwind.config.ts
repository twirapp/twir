import path from 'node:path'
import url from 'node:url'

import type { Config } from 'tailwindcss'

const currentDir = path.dirname(url.fileURLToPath(import.meta.url))

export default <Config>{
	content: [`${currentDir}/**/*.{html,js,vue,ts,tsx}`],
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
