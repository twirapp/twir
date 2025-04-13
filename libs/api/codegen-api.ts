import { fileURLToPath } from 'node:url'

import { generateApi } from 'swagger-typescript-api'

async function generateSwagger() {
	let tryCount = 0
	while (tryCount < 10) {
		try {
			await generateApi({
				name: 'openapi.ts',
				url: 'http://localhost:3009/docs/openapi.json',
				output: fileURLToPath(new URL(`./`, import.meta.url)),
				generateClient: true,
				httpClientType: 'fetch',
				singleHttpClient: true,
				extractEnums: true,
				defaultResponseAsSuccess: true,
				patch: true,
				generateResponses: true,
			})
			break
		} catch {
			tryCount++
			await new Promise((resolve) => setTimeout(resolve, 1000))
		}
	}
}

await generateSwagger()
