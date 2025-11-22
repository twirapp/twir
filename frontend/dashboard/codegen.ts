import { join, resolve } from 'node:path'
import process from 'node:process'

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
	documents: ['src/**/*.{vue,ts}'],
	ignoreNoDocuments: true, //for better experience with the watcher
	generates: {
		'./src/gql/': {
			preset: 'client',
			config: {
				useTypeImports: true,
			},
			// presetConfig: {
			// 	onExecutableDocumentNode: generatePersistHash,
			// },
			// documentTransforms: [
			// 	addTypenameSelectionDocumentTransform,
			// ],
		},
		'./src/gql/validation-schemas.ts': {
			plugins: ['./codegen-plugins/zod.ts'],
		},
	},
}

export default config
