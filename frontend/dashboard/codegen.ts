import { join, resolve } from 'node:path'
import process from 'node:process'

import { addTypenameSelectionDocumentTransform } from '@graphql-codegen/client-preset'

import { generatePersistHash } from './codegen-persist-hash'

import type { CodegenConfig } from '@graphql-codegen/cli'

const schemaDir = resolve(join(process.cwd(), '..', '..', 'apps', 'api-gql', 'schema', '*.graphqls'))

const config: CodegenConfig = {
	config: {
		scalars: {
			Upload: 'File',
		},
	},
	schema: schemaDir,
	documents: ['src/api/**/*.ts'],
	ignoreNoDocuments: true, // for better experience with the watcher
	generates: {
		'./src/gql/': {
			preset: 'client',
			config: {
				useTypeImports: true,
			},
			presetConfig: {
				onExecutableDocumentNode: generatePersistHash,
			},
			documentTransforms: [
				addTypenameSelectionDocumentTransform,
			],
		},
	},
}

export default config
