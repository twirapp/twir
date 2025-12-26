import { join, resolve } from 'node:path'
import * as process from 'node:process'

import type { CodegenConfig } from '@graphql-codegen/cli'

const schemaDir = resolve(
	join(
		process.cwd(),
		'..',
		'..',
		'apps',
		'api-gql',
		'internal',
		'delivery',
		'gql',
		'schema',
		'**',
		'*.graphql'
	)
)

const config: CodegenConfig = {
	config: {
		scalars: {
			Upload: 'File',
		},
	},
	schema: schemaDir,
	documents: ['src/**/*.ts'],
	ignoreNoDocuments: true, // for better experience with the watcher
	generates: {
		'./src/gql/': {
			preset: 'client',
			config: {
				useTypeImports: true,
			},
		},
	},
}

export default config
