import { rename } from 'node:fs/promises'
import { arch, platform } from 'node:process'

import { build } from 'bun'

function isMusl() {
	const report = process.report?.getReport() as any

	return !report.header.glibcVersionRuntime
}

const isolatedVmCondition = `${platform}-${arch}-${isMusl() ? 'musl' : 'gnu'}`

await build({
	entrypoints: ['./src/index.ts'],
	outdir: '.out',
	compile: true,
	minify: true,
	sourcemap: 'inline',
	conditions: [isolatedVmCondition],
})

await rename('.out/src', '.out/twir-executron')
