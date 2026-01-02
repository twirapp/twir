# Bug Fix: Proxy Object Clone Error

## Issue

When opening the overlay builder, the page was empty with repeated errors in the console:

```
Uncaught (in promise) DOMException: Proxy object could not be cloned.
    useOverlayBuilder useOverlayBuilder.ts:44
```

## Root Cause

The `useOverlayBuilder` composable was using `structuredClone()` to deep clone the project state for history management. However, `structuredClone()` cannot clone Vue's reactive proxy objects, which are created by Vue 3's `reactive()` function.

### Problem Code

```typescript
const project = reactive<OverlayProject>({ ... })

const history = reactive<HistoryState>({
  past: [],
  present: structuredClone(project), // ❌ ERROR: Can't clone reactive proxy
  future: [],
})
```

## Solution

Replace all `structuredClone()` calls with `JSON.parse(JSON.stringify(toRaw(object)))`:

1. Import `toRaw` from Vue
2. Use `toRaw()` to convert reactive proxies to plain objects
3. Use JSON serialization for deep cloning

### Fixed Code

```typescript
import { toRaw } from 'vue'

const project = reactive<OverlayProject>({ ... })

const history = reactive<HistoryState>({
  past: [],
  present: JSON.parse(JSON.stringify(toRaw(project))), // ✅ Works
  future: [],
})
```

## Changes Made

### File: `composables/useOverlayBuilder.ts`

Replaced all occurrences of `structuredClone()` with `JSON.parse(JSON.stringify(toRaw()))`:

1. **History initialization** (line 44)
2. **saveToHistory function** (lines 68, 71)
3. **undo function** (lines 80, 84)
4. **redo function** (lines 90, 94)
5. **duplicateLayer function** (line 166)
6. **duplicateLayers function** (line 187)
7. **copyToClipboard function** (line 283)
8. **pasteFromClipboard function** (line 300)
9. **loadProject function** (lines 482, 483)
10. **exportProject function** (line 491)

## Why JSON.parse(JSON.stringify()) Works

1. **toRaw()**: Converts Vue reactive proxy to plain JavaScript object
2. **JSON.stringify()**: Serializes the plain object to string
3. **JSON.parse()**: Deserializes back to a new plain object (deep clone)

This creates a true deep copy without any Vue reactivity attached.

## Alternative Considered

We could use a library like `lodash.cloneDeep`, but JSON serialization is:
- Built-in (no extra dependencies)
- Fast enough for our use case
- Simple and well-understood
- Works with all serializable data

## Limitations

JSON serialization doesn't handle:
- Functions (not needed in our state)
- Dates (convert to ISO strings)
- RegExp (not used in state)
- Symbols (not used in state)
- Circular references (not present in our state)

For our overlay state (positions, dimensions, strings), JSON serialization is perfect.

## Testing

After fix:
- ✅ Page loads without errors
- ✅ Overlay data displays correctly
- ✅ Undo/redo works
- ✅ Copy/paste works
- ✅ Duplicate works
- ✅ History management works

## Related Issues

This is a common issue when working with Vue 3's Composition API and trying to clone reactive objects. Always use `toRaw()` before cloning reactive data.

## References

- [Vue 3 toRaw() documentation](https://vuejs.org/api/reactivity-advanced.html#toraw)
- [MDN structuredClone()](https://developer.mozilla.org/en-US/docs/Web/API/structuredClone)
- [Vue Reactivity in Depth](https://vuejs.org/guide/extras/reactivity-in-depth.html)
