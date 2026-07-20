# Executron Lodash Global Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Expose full Lodash 4.17.21 as the global `_` in Executron scripts and provide matching Monaco autocomplete in the variables editor.

**Architecture:** Bun embeds Lodash's official browser build as source text in the Executron binary, and each isolated realm evaluates it before user code. The web app registers the official `@types/lodash` declaration set as Monaco virtual files alongside the existing dynamic Twir declarations.

**Tech Stack:** Bun, TypeScript, `@isolated-vm/experimental`, Lodash 4.17.21, Monaco Editor, Nuxt 4/Vite 8

---

## File Structure

- Create `apps/executron/src/libs/executor.test.ts`: verifies static and chain Lodash APIs inside the actual isolated realm.
- Create `apps/executron/src/types/lodash-source.d.ts`: types the browser-build text import.
- Modify `apps/executron/src/libs/executor.ts`: embeds and initializes Lodash before user code.
- Modify `apps/executron/package.json`: adds the pinned runtime dependency.
- Create `web/layers/dashboard/features/variables/lodash-monaco-types.ts`: maps official Lodash declarations to Monaco virtual file paths.
- Modify `web/layers/dashboard/features/variables/variables-edit.vue`: registers and disposes all Twir and Lodash extra libraries.
- Modify `web/package.json`: adds the pinned Lodash type dependency.
- Modify `bun.lock`: records both workspace dependency changes.

### Task 1: Prove Lodash Runtime Behavior

**Files:**
- Create: `apps/executron/src/libs/executor.test.ts`

- [ ] **Step 1: Write the failing sandbox test**

```typescript
import { describe, expect, test } from 'bun:test'

import { executeCode } from './executor'

describe('executeCode Lodash global', () => {
	test('supports static and chain APIs through _', async () => {
		const execution = await executeCode(
			`
				const chunks = _.chunk([1, 2, 3, 4, 5], 2);
				const chained = _.chain([1, 2, 3, 4])
					.filter((value) => value % 2 === 0)
					.map((value) => value * 10)
					.value();

				return JSON.stringify({ chunks, chained, version: _.VERSION });
			`,
			'test-channel',
			new Map()
		)

		expect(execution.error).toBe('')
		expect(JSON.parse(execution.result)).toEqual({
			chunks: [[1, 2], [3, 4], [5]],
			chained: [20, 40],
			version: '4.17.21',
		})
	})
})
```

- [ ] **Step 2: Run the test and verify the missing global**

Run from `apps/executron`:

```bash
bun test src/libs/executor.test.ts
```

Expected: FAIL because the sandbox reports `_ is not defined`.

### Task 2: Embed Lodash in Executron

**Files:**
- Modify: `apps/executron/package.json`
- Modify: `bun.lock`
- Create: `apps/executron/src/types/lodash-source.d.ts`
- Modify: `apps/executron/src/libs/executor.ts:1-24,260-368`
- Test: `apps/executron/src/libs/executor.test.ts`

- [ ] **Step 1: Install the pinned runtime dependency**

Run from `apps/executron`:

```bash
bun add lodash@4.17.21
```

Expected: `lodash: 4.17.21` appears in `apps/executron/package.json`, and `bun.lock` changes.

- [ ] **Step 2: Type the source-text import**

Create `apps/executron/src/types/lodash-source.d.ts`:

```typescript
declare module 'lodash/lodash.min.js' {
	const source: string
	export default source
}
```

- [ ] **Step 3: Import the browser build as embedded text**

Add near the external imports in `apps/executron/src/libs/executor.ts`:

```typescript
import lodashSource from 'lodash/lodash.min.js' with { type: 'text' }
```

- [ ] **Step 4: Initialize `_` before sandbox APIs and user code**

Add immediately after `const __secrets = ${secretsJson};` in `wrappedCode`:

```typescript
			${lodashSource}
			const _ = globalThis._;
```

This evaluates the official build inside the isolated realm. The lexical `const _` guarantees that nested user code resolves the global API without importing a module or crossing the host boundary.

- [ ] **Step 5: Run the runtime test**

Run from `apps/executron`:

```bash
bun test src/libs/executor.test.ts
```

Expected: PASS with one test, including `version: '4.17.21'`.

- [ ] **Step 6: Run Executron type checking and compilation**

Run from `apps/executron`:

```bash
bun run build
```

Expected: TypeScript exits successfully and `.out/twir-executron` is rebuilt with Lodash embedded.

### Task 3: Register Official Lodash Types in Monaco

**Files:**
- Modify: `web/package.json`
- Modify: `bun.lock`
- Create: `web/layers/dashboard/features/variables/lodash-monaco-types.ts`
- Modify: `web/layers/dashboard/features/variables/variables-edit.vue:25-61`

- [ ] **Step 1: Install the compatible pinned declarations**

Run from `web`:

```bash
bun add --dev @types/lodash@4.17.24
```

Expected: `@types/lodash: 4.17.24` appears in `web/package.json`, and `bun.lock` changes.

- [ ] **Step 2: Create the Monaco declaration registry**

Create `web/layers/dashboard/features/variables/lodash-monaco-types.ts`:

```typescript
import arrayTypes from '@types/lodash/common/array.d.ts?raw'
import collectionTypes from '@types/lodash/common/collection.d.ts?raw'
import commonTypes from '@types/lodash/common/common.d.ts?raw'
import dateTypes from '@types/lodash/common/date.d.ts?raw'
import functionTypes from '@types/lodash/common/function.d.ts?raw'
import langTypes from '@types/lodash/common/lang.d.ts?raw'
import mathTypes from '@types/lodash/common/math.d.ts?raw'
import numberTypes from '@types/lodash/common/number.d.ts?raw'
import objectTypes from '@types/lodash/common/object.d.ts?raw'
import seqTypes from '@types/lodash/common/seq.d.ts?raw'
import stringTypes from '@types/lodash/common/string.d.ts?raw'
import utilTypes from '@types/lodash/common/util.d.ts?raw'
import indexTypes from '@types/lodash/index.d.ts?raw'

export const lodashMonacoTypeDefinitions: Record<string, string> = {
	'file:///node_modules/@types/lodash/index.d.ts': indexTypes,
	'file:///node_modules/@types/lodash/common/common.d.ts': commonTypes,
	'file:///node_modules/@types/lodash/common/array.d.ts': arrayTypes,
	'file:///node_modules/@types/lodash/common/collection.d.ts': collectionTypes,
	'file:///node_modules/@types/lodash/common/date.d.ts': dateTypes,
	'file:///node_modules/@types/lodash/common/function.d.ts': functionTypes,
	'file:///node_modules/@types/lodash/common/lang.d.ts': langTypes,
	'file:///node_modules/@types/lodash/common/math.d.ts': mathTypes,
	'file:///node_modules/@types/lodash/common/number.d.ts': numberTypes,
	'file:///node_modules/@types/lodash/common/object.d.ts': objectTypes,
	'file:///node_modules/@types/lodash/common/seq.d.ts': seqTypes,
	'file:///node_modules/@types/lodash/common/string.d.ts': stringTypes,
	'file:///node_modules/@types/lodash/common/util.d.ts': utilTypes,
}
```

The virtual paths intentionally mirror the package layout so the triple-slash references in `index.d.ts` resolve inside Monaco.

- [ ] **Step 3: Import the declaration registry into the variables editor**

Add in `web/layers/dashboard/features/variables/variables-edit.vue` beside `useTwirMonacoTypes`:

```typescript
import { lodashMonacoTypeDefinitions } from './lodash-monaco-types'
```

- [ ] **Step 4: Register and dispose every virtual declaration file**

Replace the single `extraLib` and `registerExtraLib` implementation with:

```typescript
let extraLibs: Array<{ dispose: () => void }> = []

function registerExtraLibs(monaco: any, twirDefinitions: string) {
	for (const extraLib of extraLibs) {
		extraLib.dispose()
	}

	extraLibs = [
		monaco.languages.typescript.javascriptDefaults.addExtraLib(
			twirDefinitions,
			'file:///twir-globals.d.ts',
		),
		...Object.entries(lodashMonacoTypeDefinitions).map(([filePath, definitions]) =>
			monaco.languages.typescript.javascriptDefaults.addExtraLib(definitions, filePath)
		),
	]
}
```

Update both call sites from `registerExtraLib(monaco, defs)` to `registerExtraLibs(monaco, defs)` and from `registerExtraLib(monacoInstance, defs)` to `registerExtraLibs(monacoInstance, defs)`.

- [ ] **Step 5: Run the web type checker**

Run from `web`:

```bash
bunx --bun nuxi typecheck
```

Expected: PASS with no unresolved `?raw` imports and no Vue/TypeScript errors from the changed files.

- [ ] **Step 6: Build the Nuxt app**

Run from `web`:

```bash
bun run build
```

Expected: PASS; Vite bundles all Lodash declaration text into the variables editor chunk.

### Task 4: Final Verification

**Files:**
- Verify: `apps/executron/src/libs/executor.test.ts`
- Verify: `apps/executron/src/libs/executor.ts`
- Verify: `web/layers/dashboard/features/variables/lodash-monaco-types.ts`
- Verify: `web/layers/dashboard/features/variables/variables-edit.vue`

- [ ] **Step 1: Re-run focused runtime verification**

Run from `apps/executron`:

```bash
bun test src/libs/executor.test.ts && bun run build
```

Expected: both commands pass.

- [ ] **Step 2: Re-run focused frontend verification**

Run from `web`:

```bash
bunx --bun nuxi typecheck && bun run build
```

Expected: both commands pass.

- [ ] **Step 3: Inspect only the intended diff**

Run from the repository root:

```bash
git diff -- apps/executron web/layers/dashboard/features/variables web/package.json bun.lock docs/superpowers
```

Expected: only the Lodash runtime, Monaco types, tests, dependency lock changes, design, and plan are present. Do not modify or revert unrelated working-tree changes.

No commit is created unless the user explicitly requests one.
