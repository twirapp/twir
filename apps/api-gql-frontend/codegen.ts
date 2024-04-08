import { join, resolve } from 'path';

import type { CodegenConfig } from '@graphql-codegen/cli';

const schemaDir = resolve(join(process.cwd(), '..', 'api-gql', 'schema', '*.graphqls'));

const config: CodegenConfig = {
	schema: schemaDir,
	watch: true,
	documents: ['src/**/*.vue', 'src/**/*.ts'],
	ignoreNoDocuments: true, // for better experience with the watcher
	generates: {
		'./src/gql/': {
			preset: 'client',
			config: {
				useTypeImports: true,
			},
		},
	},
};

export default config;
