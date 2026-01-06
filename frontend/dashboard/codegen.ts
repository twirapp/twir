import type { CodegenConfig } from '@graphql-codegen/cli'

import { join, resolve } from 'node:path'
import * as process from 'node:process'

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
	documents: ['src/**/*.{vue,ts}'],
	ignoreNoDocuments: true, //for better experience with the watcher
	generates: {
		'./src/gql/': {
			preset: 'client',
			config: {
				useTypeImports: true,
			},
			presetConfig: {
				// persistedDocuments: {
				// 	hashAlgorithm: (operation: string) => {
				// 		const h = new Bun.CryptoHasher('sha256')
				// 		h.update(operation)
				// 		return h.digest('hex')
				// 	},
				// },
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
