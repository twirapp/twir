# Executron Lodash Global Design

## Goal

Expose the full Lodash 4.17.21 API to Executron user scripts as the global `_` identifier and provide matching autocomplete and type checking in the variables Monaco editor.

## Runtime

- Add `lodash@4.17.21` to `apps/executron`.
- Embed Lodash's official browser build into the isolated realm before Executron evaluates user code.
- Keep `_` global; user scripts do not import Lodash.
- Run Lodash inside the realm rather than proxying host functions, preserving synchronous methods, callbacks, wrappers, and chain behavior without crossing the isolation boundary.
- Keep the existing execution timeout and sandbox capabilities unchanged.

## Editor Types

- Add `@types/lodash` to `web` development dependencies at the version matching Lodash 4.17.
- Include the Lodash declarations in the Monaco extra libraries registered by the variables editor.
- Preserve the existing generated declarations for `twir`, storage, secrets, fetch, timers, and sandbox APIs.
- Declare the same global `_` API that exists at runtime.

## Error Handling

- Treat failure to bundle or type-check Lodash as a build failure.
- Do not load Lodash from a CDN or perform network access during script execution.
- Existing Executron runtime error and timeout handling remains responsible for user-code failures.

## Verification

- Add or extend an Executron test that executes user code using representative static and chain APIs, such as `_.chunk` and `_.chain`.
- Verify the global is available without an import.
- Run the Executron type check/tests and build.
- Run the relevant web type check/build to verify the Monaco declaration imports.
