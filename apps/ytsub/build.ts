import { rename } from 'node:fs/promises'

import { build } from 'bun'

await build({
	entrypoints: ['./src/index.ts'],
	outdir: '.out',
	compile: true,
	minify: true,
	sourcemap: 'inline',
})

await rename('.out/src', '.out/twir-ytsub')
