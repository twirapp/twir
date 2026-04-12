# Vue 3 + TypeScript Best Practices Skill

This skill provides a comprehensive, production-quality reference for writing Vue 3 code in TypeScript desktop applications built with Electrobun and Vite.

## Contents

**SKILL.md** (1,518 lines, ~34 KB) covers:

### 1. Composables

- When to extract (reusable, stateful, side-effect heavy)
- Naming convention (`use*`)
- Return patterns (refs vs reactive)
- Reactivity rules (ref vs reactive, toRefs)
- Side effects and cleanup (onUnmounted)
- TypeScript patterns (generics, return types)

### 2. Pinia State Management

- Store definition (Option vs Setup stores)
- TypeScript store typing
- Actions (sync, async, error handling)
- Getters with arguments
- Store subscriptions
- Correct access patterns (storeToRefs for destructuring)
- Composable stores (using other stores)

### 3. Component Design

- `<script setup>` syntax
- Props with TypeScript interfaces
- Emits with strong typing
- defineModel (Vue 3.4+)
- defineExpose for template refs
- Computed vs methods
- Conditional rendering with TypeScript guards
- Full component examples

### 4. Avoiding `any`

- Props typing (never any)
- Event handler typing (MouseEvent, KeyboardEvent, etc.)
- Template refs (HTMLElement | null)
- Composable return typing (generics)
- Provide/Inject typing (InjectionKey)
- Emits typing (with payloads)

### 5. Component Composition

- Named slots
- Scoped slots
- Renderless components
- Provide/Inject pattern with context
- Advanced v-model with modifiers

### 6. Reactivity System

- ref vs reactive comparison
- toRef / toRefs (extract reactivity)
- watch vs watchEffect (explicit vs automatic)
- Shallow reactivity (shallowReactive, shallowRef)
- Computed with dependencies
- Utility functions (unref, isRef, isReactive, isReadonly)

### 7. Advanced Patterns

- VueUse common composables (useLocalStorage, useAsyncState, useEventListener, useDebounceFn)
- Error handling composables
- Async state in components
- Status management (idle/loading/success/error)

### 8. TypeScript Configuration

- Recommended tsconfig.json with strict mode enabled
- All compiler options explained

### 9-10. Checklists & Resources

- Component shipping checklist (12 items)
- Composable audit checklist (11 items)
- Links to official docs (Vue 3, Pinia, VueUse, Vue Language Tools)
- Mental model summary for TwirChat frontend

## Usage

Load this skill before writing Vue 3 code:

```bash
# Agent usage (pseudo-code)
load_skill("vue3-typescript-best-practices")
```

Then reference patterns and examples from SKILL.md when implementing:

- New Vue components
- Composables for state/logic
- Pinia stores
- TypeScript types for props/emits

## Key Principles

1. **Type Everything**: No `any` - use interfaces, generics, and type guards
2. **Composable Discipline**: Name with `use*`, return object of refs, manage cleanup
3. **Reactivity Fundamentals**: Use `ref()` for primitives, `reactive()` for objects, `toRefs()` for destructuring
4. **Pinia Patterns**: Setup stores with TypeScript, access via `storeToRefs()` when destructuring
5. **Component Safety**: `<script setup>`, typed props/emits, proper event handlers

## Document Stats

- **Lines**: 1,518
- **Size**: ~34 KB
- **Code Blocks**: 52 (104 fence markers)
- **Sections**: 10
- **Subsections**: 45
- **Concrete Examples**: 50+
- **Checklists**: 2

## Target Audience

Frontend developers building:

- Electrobun + Vue 3 desktop applications
- TypeScript-first codebases
- Production-quality components with best practices
- Scalable state management with Pinia
- Reusable composables with proper typing

## Notes

- All examples use modern Vue 3.4+ syntax (defineModel, etc.)
- Zero beginner content - production-quality only
- Every pattern has ✅ correct and ❌ incorrect examples
- Community patterns from VueUse documented
- Focused on desktop app context (Electrobun + Vite)

## Integration with Project

This skill is designed specifically for **TwirChat** desktop application development:

- **Target**: `packages/desktop/src/views/` Vue components
- **Context**: Electrobun v1.16.0 + Vue 3 + TypeScript
- **Patterns**: Composition API with `<script setup>`
- **State**: Pinia stores + Electrobun RPC + Backend WebSocket
- **Build**: Vite (HMR for main window, static for overlay)

### Example Application

When building a chat message component for TwirChat:

1. Load this skill
2. Check section 3 (Component Design) for prop/emit patterns
3. Reference section 4 (Avoiding `any`) for event typing
4. Use composable checklist before extracting `useChatState`
5. Follow Pinia patterns (section 2) for global chat store
6. Use component shipping checklist before merging

## Related Skills

- `7tv-events-api` — 7TV EventAPI v3 implementation reference
- (Future) Backend API patterns — Bun REST/WebSocket
- (Future) Electron/Electrobun best practices
