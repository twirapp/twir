import { copyFileSync, existsSync, readdirSync } from 'node:fs'
import { join } from 'node:path'
import { arch, platform } from 'node:process'

function detectHostPlatform(): string {
	let p = platform
	if (p === 'linux') {
		try {
			const fs = require('node:fs')
			if (fs.readFileSync('/usr/bin/ldd', 'latin1').includes('ld-musl-')) {
				return `linux-${arch}-musl`
			}
		} catch {}
		return `linux-${arch}-gnu`
	}
	return `${p}-${arch}`
}

const target = process.env.BUILD_TARGET || detectHostPlatform()
console.log(`Target platform: ${target}`)

const bunDir = join(import.meta.dirname, '..', '..', '..', 'node_modules', '.bun')
const outDir = join(import.meta.dirname, '..', 'src')

function tryCopy(plat: string): boolean {
	const prefix = `@isolated-vm+experimental-${plat}@`
	const entries = readdirSync(bunDir)
	const pkgEntry = entries.find((e) => e.startsWith(prefix))

	if (!pkgEntry) return false

	const nodeFile = `backend_napi_v8.${plat}.node`
	const srcPath = join(
		bunDir,
		pkgEntry,
		'node_modules',
		'@isolated-vm',
		`experimental-${plat}`,
		nodeFile
	)

	if (!existsSync(srcPath)) return false

	copyFileSync(srcPath, join(outDir, 'backend.node'))
	console.log(`Copied ${nodeFile} -> backend.node`)
	return true
}

if (!tryCopy(target)) {
	const fallback = target.endsWith('-gnu')
		? target.replace('-gnu', '-musl')
		: target.replace('-musl', '-gnu')

	console.log(`${target} not found, trying ${fallback}...`)

	if (!tryCopy(fallback)) {
		console.error(`Could not find isolated-vm native addon for ${target} or ${fallback}`)
		process.exit(1)
	}
}
