import { join } from 'node:path'

import { defineNuxtModule } from '@nuxt/kit'
import { debounce } from '@std/async/debounce'
import { $, Glob } from 'bun'

import codegenConfig from '../codegen'

const cwd = join(import.meta.path, '..', '..')
if (!codegenConfig.documents || !Array.isArray(codegenConfig.documents)) throw new Error('codegenConfig.documents is required')

const globs = codegenConfig.documents!.map((doc) => {
	return new Glob(join(cwd, doc as string))
})
const runBuild = debounce(async () => await $`bun run graphql-codegen`, 1000)

export default defineNuxtModule((options, nuxt) => {
	nuxt.hook('build:before', async () => {
		await $`bun run graphql-codegen`
	})

	if (nuxt.options.dev) {
		nuxt.hook('builder:watch', async (event, path) => {
			if (!globs.some((glob) => glob.match(path))) return

			runBuild()
		})
	}
})
