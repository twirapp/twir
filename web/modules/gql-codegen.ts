import { defineNuxtModule } from '@nuxt/kit'
import { $, Glob } from 'bun'
import { join } from 'node:path'

import codegenConfig from '../codegen'

const cwd = join(import.meta.path, '..', '..')
if (!codegenConfig.documents || !Array.isArray(codegenConfig.documents)) {
	throw new Error('codegenConfig.documents is required')
}

const globs = codegenConfig.documents!.map((doc) => {
	return new Glob(join(cwd, doc as string))
})
const runBuild = debounce(async () => await $`bun run graphql-codegen`, 1000)

let alreadyBuilt = false

export default defineNuxtModule((_options, nuxt) => {
	nuxt.hook('build:before', async () => {
		if (alreadyBuilt) return
		await $`bun run graphql-codegen`
		alreadyBuilt = true
	})

	if (nuxt.options.dev) {
		nuxt.hook('builder:watch', async (event, path) => {
			if (!globs.some((glob) => glob.match(path))) return

			runBuild()
		})
	}
})

interface DebouncedFunction<T extends Array<unknown>> {
	(...args: T): void
	/** Clears the debounce timeout and omits calling the debounced function. */
	clear(): void
	/** Clears the debounce timeout and calls the debounced function immediately. */
	flush(): void
	/** Returns a boolean whether a debounce call is pending or not. */
	readonly pending: boolean
}

function debounce<T extends Array<any>>(
	fn: (this: DebouncedFunction<T>, ...args: T) => void,
	wait: number
): DebouncedFunction<T> {
	let timeout: number | null = null
	let flush: (() => void) | null = null

	const debounced: DebouncedFunction<T> = ((...args: T) => {
		debounced.clear()
		flush = () => {
			debounced.clear()
			fn.call(debounced, ...args)
		}
		timeout = Number(setTimeout(flush, wait))
	}) as DebouncedFunction<T>

	debounced.clear = () => {
		if (typeof timeout === 'number') {
			clearTimeout(timeout)
			timeout = null
			flush = null
		}
	}

	debounced.flush = () => {
		flush?.()
	}

	Object.defineProperty(debounced, 'pending', {
		get: () => typeof timeout === 'number',
	})

	return debounced
}
