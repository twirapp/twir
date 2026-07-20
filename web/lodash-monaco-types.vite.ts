import type { Plugin } from 'vite'

import { readFileSync } from 'node:fs'
import { createRequire } from 'node:module'
import { dirname, join } from 'node:path'

const virtualModuleId = 'virtual:lodash-monaco-types'
const resolvedVirtualModuleId = `\0${virtualModuleId}`
const declarationFiles = [
	'index.d.ts',
	'common/common.d.ts',
	'common/array.d.ts',
	'common/collection.d.ts',
	'common/date.d.ts',
	'common/function.d.ts',
	'common/lang.d.ts',
	'common/math.d.ts',
	'common/number.d.ts',
	'common/object.d.ts',
	'common/seq.d.ts',
	'common/string.d.ts',
	'common/util.d.ts',
]

export function lodashMonacoTypesPlugin(): Plugin {
	return {
		name: 'twir-lodash-monaco-types',
		resolveId(id) {
			if (id === virtualModuleId) return resolvedVirtualModuleId
		},
		load(id) {
			if (id !== resolvedVirtualModuleId) return

			const require = createRequire(import.meta.url)
			const typesDirectory = dirname(require.resolve('@types/lodash/package.json'))
			const definitions = Object.fromEntries(
				declarationFiles.map((file) => [
					`file:///node_modules/@types/lodash/${file}`,
					readFileSync(join(typesDirectory, file), 'utf8'),
				])
			)

			return `export default ${JSON.stringify(definitions)}`
		},
	}
}
