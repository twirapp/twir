import path from 'node:path'
import process from 'node:process'
import url from 'node:url'

import { config, readEnv } from '@twir/config'

readEnv(path.join(process.cwd(), '..', '..', '.env'))

const https = config.TWITCH_CALLBACKURL.startsWith('https')

const currentDir = path.dirname(url.fileURLToPath(import.meta.url))

export default defineNuxtConfig({
	modules: ['@bicou/nuxt-urql'],

	urql: {
		endpoint: process.env.NODE_ENV !== 'production'
			? `${https ? 'https' : 'http'}://${config.SITE_BASE_URL}/api/query`
			: 'http://api-gql:3009/query',
		client: path.join(currentDir, 'urql.ts'),
		ssr: {
			endpoint: process.env.NODE_ENV !== 'production'
				? `${https ? 'https' : 'http'}://${config.SITE_BASE_URL}/api/query`
				: 'http://api-gql:3009/query',
			key: 'urql-ssr'
		},
	},
})
