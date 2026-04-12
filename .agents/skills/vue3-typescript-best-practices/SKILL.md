---
name: vue3-best-practices
description: Production-quality Vue 3 + TypeScript reference for TwirChat desktop. Covers composables, Pinia, component design, SVG-as-components, Electrobun RPC constraints, and patterns to avoid (any, inline SVG, logic in components).
license: MIT
compatibility: opencode
---

# Vue 3 + TypeScript Best Practices

**Purpose**: Production-quality reference for writing Vue 3 code in TypeScript desktop apps (Electrobun + Vite).

**Scope**: Composables, Pinia state management, component design, TypeScript integration, reactivity patterns, SVG workflow, and TwirChat-specific conventions.

---

## 0. TwirChat-Specific Constraints (READ FIRST)

These rules are **non-negotiable** for TwirChat frontend code.

### The Two-Runtime Rule

The desktop app has two distinct runtime environments. **Never cross them.**

| Environment      | Location     | Access                                |
| ---------------- | ------------ | ------------------------------------- |
| **Main process** | `src/bun/`   | Full Bun/Node.js, SQLite, file system |
| **Webview**      | `src/views/` | Browser context only — NO Bun APIs    |

**Forbidden in `src/views/`:**

- `bun:sqlite` — use RPC `get*` / `set*` methods instead
- `node:fs` — use RPC for file operations
- Direct imports from `src/store/` — SQLite runs on Bun side only

```typescript
// ❌ WRONG — will crash the browser at runtime
import { SettingsStore } from '../../store/settings-store'
const settings = SettingsStore.get()

// ✅ CORRECT — use Electrobun RPC
const settings = await rpc.request.getSettings()
```

### RPC Pattern (How Views Talk to Main Process)

```typescript
// src/shared/rpc.ts — defines the schema
// src/bun/index.ts — implements bun-side handlers
// src/views/main/main.ts — view-side Electroview.defineRPC

// In any Vue component or composable:
import { rpc } from '../main' // views/main/main.ts export

// Request (awaitable, returns data)
const accounts = await rpc.request.getAccounts()

// Message (fire-and-forget RPC call to bun side)
rpc.send.someMessage({ payload: 'value' })

// Listen for messages pushed FROM bun side
rpc.on.chat_message((msg) => {
  /* ... */
})
```

**Always wrap RPC calls in composables** — never call `rpc.request.*` directly in component setup.

### Overlay vs Main Window

The overlay (`src/views/overlay/`) has **no Electrobun RPC** — it uses a plain WebSocket to `overlay-server.ts`. Do not import `rpc` in overlay code.

---

---

## 1. Composables: Extraction, Naming, and Patterns

### When to Extract a Composable

Extract when logic is:

- **Reusable** across 2+ components
- **Stateful** (maintains internal state over time)
- **Side-effect heavy** (subscriptions, timers, cleanup required)
- **Complex** enough to warrant naming

❌ **Don't extract**:

```typescript
// Single-use logic, just put it in the component
const handleClick = () => {
  count.value++
}
```

✅ **Do extract**:

```typescript
// Reusable, stateful, has cleanup
export function useFormValidation(initialData: T) {
  const data = ref(initialData)
  const errors = reactive<Record<string, string>>({})
  const isDirty = computed(() => !isEqual(data.value, initialData))

  const reset = () => {
    data.value = initialData
    errors.value = {}
  }

  onUnmounted(() => {
    // cleanup
  })

  return { data, errors, isDirty, reset }
}
```

### Naming Convention: `use*`

All composables must start with `use`:

- ✅ `useLocalStorage(key)`
- ✅ `useAsyncData(url)`
- ✅ `useFormValidation(schema)`
- ❌ `getLocalStorage()` (function, not composable)
- ❌ `createValidator()` (use `useValidator`)

### Return Patterns

**Pattern 1: Direct refs/computed** (preferred)

```typescript
export function useCounter() {
  const count = ref(0)
  const doubled = computed(() => count.value * 2)
  const increment = () => {
    count.value++
  }

  return { count, doubled, increment }
}
```

**Pattern 2: Object with reactive** (use sparingly)

```typescript
export function useFormState(initial: T) {
  const state = reactive({ ...initial, loading: false })

  // BUT: destructure returns break reactivity
  // const { name, email } = useFormState() // ❌ WRONG
  // Use: const formState = useFormState(); formState.name // ✅ CORRECT
  // Or: const { name, email } = toRefs(useFormState()) // ✅ CORRECT

  return state
}
```

**Pattern 3: Never mix** `reactive()` **with destructuring**

```typescript
// ❌ WRONG: breaks reactivity
export function useCounter() {
  const state = reactive({ count: 0 })
  return state // Caller will destructure and lose reactivity
}

// ✅ RIGHT: return refs or use toRefs
export function useCounter() {
  const count = ref(0)
  return { count }
}
```

### Reactivity Rules Inside Composables

**Rule 1**: `ref()` for primitives, `reactive()` for objects

```typescript
// ✅ CORRECT
const count = ref(0) // primitive → ref
const user = reactive({ name: '', age: 0 }) // object → reactive

// ❌ WRONG
const count = reactive(0) // primitive in reactive
const user = ref({ name: '', age: 0 }) // object in ref (awkward access)
```

**Rule 2**: Always use `toRefs()` when returning reactive objects

```typescript
export function useUser() {
  const user = reactive({ name: '', email: '' })

  // ❌ WRONG: caller destructures and loses reactivity
  return user

  // ✅ CORRECT: caller can destructure safely
  return toRefs(user)
}

// Usage
const { name, email } = useUser() // Reactivity preserved
```

**Rule 3**: Use `readonly()` for exposed state you don't want modified

```typescript
export function useAppConfig() {
  const config = reactive({ apiUrl: 'https://...', timeout: 5000 })

  const updateConfig = (newConfig: Partial<typeof config>) => {
    Object.assign(config, newConfig)
  }

  return { config: readonly(config), updateConfig }
}
```

### Side Effects and Cleanup

```typescript
export function useEventListener(target: string, event: string, handler: Function) {
  const el = ref<HTMLElement | null>(null)

  onMounted(() => {
    el.value = document.querySelector(target)
    el.value?.addEventListener(event, handler as EventListener)
  })

  onUnmounted(() => {
    el.value?.removeEventListener(event, handler as EventListener)
  })

  return { el }
}
```

### Composable TypeScript Patterns

**Generic composables**:

```typescript
export function useAsync<T>(fn: () => Promise<T>) {
  const data = ref<T | null>(null)
  const loading = ref(false)
  const error = ref<Error | null>(null)

  const execute = async () => {
    loading.value = true
    error.value = null
    try {
      data.value = await fn()
    } catch (e) {
      error.value = e as Error
    } finally {
      loading.value = false
    }
  }

  return { data: readonly(data), loading: readonly(loading), error: readonly(error), execute }
}

// Usage
const { data, loading, error, execute } = useAsync(async () => {
  return (await fetch('/api/users').then((r) => r.json())) as User[]
})
```

---

## 2. Pinia: Store Definition, Actions, Getters, TypeScript Typing

### Store Basics: `defineStore`

```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

// Option API store (traditional)
export const useCounterStore = defineStore('counter', {
  state: () => ({
    count: 0,
    name: 'Counter',
  }),

  getters: {
    doubled: (state) => state.count * 2,
    description: (state) => `${state.name}: ${state.count}`,
  },

  actions: {
    increment() {
      this.count++
    },
    setCount(val: number) {
      this.count = val
    },
  },
})
```

### Setup Stores (Recommended with TypeScript)

```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useCounterStore = defineStore('counter', () => {
  const count = ref(0)
  const name = ref('Counter')

  const doubled = computed(() => count.value * 2)
  const description = computed(() => `${name.value}: ${count.value}`)

  function increment() {
    count.value++
  }

  function setCount(val: number) {
    count.value = val
  }

  return { count, name, doubled, description, increment, setCount }
})
```

### TypeScript Store Typing

```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

interface User {
  id: string
  name: string
  email: string
}

interface UserStoreState {
  users: User[]
  currentUser: User | null
  loading: boolean
  error: Error | null
}

export const useUserStore = defineStore('user', () => {
  const users = ref<User[]>([])
  const currentUser = ref<User | null>(null)
  const loading = ref(false)
  const error = ref<Error | null>(null)

  const getUserById = computed(() => (id: string) => users.value.find((u) => u.id === id))

  async function fetchUsers() {
    loading.value = true
    error.value = null
    try {
      const response = await fetch('/api/users')
      users.value = await response.json()
    } catch (e) {
      error.value = e as Error
    } finally {
      loading.value = false
    }
  }

  function setCurrentUser(user: User | null) {
    currentUser.value = user
  }

  return {
    users,
    currentUser,
    loading,
    error,
    getUserById,
    fetchUsers,
    setCurrentUser,
  }
})
```

### Actions: Mutations, Side Effects, Async

```typescript
export const useAccountStore = defineStore('account', () => {
  const account = ref<Account | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Simple action (synchronous)
  function setAccount(acc: Account) {
    account.value = acc
  }

  // Async action with error handling
  async function login(email: string, password: string) {
    loading.value = true
    error.value = null
    try {
      const response = await fetch('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password }),
      })
      if (!response.ok) throw new Error(response.statusText)
      account.value = await response.json()
    } catch (e) {
      error.value = (e as Error).message
    } finally {
      loading.value = false
    }
  }

  // Action that uses other actions
  async function logout() {
    try {
      await fetch('/auth/logout', { method: 'POST' })
      setAccount(null) // call another action
    } catch (e) {
      error.value = (e as Error).message
    }
  }

  return { account, loading, error, setAccount, login, logout }
})
```

### Getters with Arguments

```typescript
export const useProductStore = defineStore('product', () => {
  const products = ref<Product[]>([])

  // Getter that returns a function
  const getProductsByCategory = computed(() => (category: string) => {
    return products.value.filter((p) => p.category === category)
  })

  // Usage
  const electronics = getProductsByCategory.value('electronics')

  return { products, getProductsByCategory }
})
```

### Store Subscriptions

```typescript
export const useSettingsStore = defineStore('settings', () => {
  const theme = ref<'light' | 'dark'>('light')
  const language = ref('en')

  function setTheme(t: typeof theme.value) {
    theme.value = t
  }

  return { theme, language, setTheme }
})

// In component
const settingsStore = useSettingsStore()

// Subscribe to changes
const unsubscribe = settingsStore.$subscribe((mutation, state) => {
  console.log('Settings changed:', mutation.type, state)
  localStorage.setItem('settings', JSON.stringify(state))
})

// Don't forget cleanup
onUnmounted(() => unsubscribe())
```

### Store Access: Using `storeToRefs`

```typescript
import { storeToRefs } from 'pinia'

export default {
  setup() {
    const counterStore = useCounterStore()

    // ❌ WRONG: destructuring breaks reactivity
    // const { count, doubled } = counterStore

    // ✅ CORRECT: use storeToRefs for destructuring
    const { count, doubled } = storeToRefs(counterStore)

    // ✅ CORRECT: direct access (no destructuring)
    // counterStore.count, counterStore.doubled

    return { count, doubled }
  },
}
```

### Composable Stores Pattern

```typescript
export const useCartStore = defineStore('cart', () => {
  // Use other stores
  const productStore = useProductStore()

  const items = ref<CartItem[]>([])

  const total = computed(() => {
    return items.value.reduce((sum, item) => {
      const product = productStore.getProductsByCategory.value(item.id)
      return sum + (product?.price ?? 0) * item.quantity
    }, 0)
  })

  async function addItem(productId: string, quantity: number) {
    const existing = items.value.find((i) => i.id === productId)
    if (existing) {
      existing.quantity += quantity
    } else {
      items.value.push({ id: productId, quantity })
    }
    // Trigger action from another store
    await productStore.fetchProducts()
  }

  return { items, total, addItem }
})
```

---

## 3. Component Design: `<script setup>`, Props, Emits, defineModel

### Basic `<script setup>` Component

```vue
<template>
  <button @click="handleClick">{{ label }}: {{ count }}</button>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Props {
  label?: string
  initialCount?: number
}

withDefaults(defineProps<Props>(), {
  label: 'Count',
  initialCount: 0,
})

const emit = defineEmits<{
  increment: [newCount: number]
}>()

const count = ref(0)

function handleClick() {
  count.value++
  emit('increment', count.value)
}
</script>
```

### Props with TypeScript Interfaces

```vue
<script setup lang="ts">
interface Props {
  modelValue: string
  placeholder?: string
  disabled?: boolean
  maxLength?: number
  type?: 'text' | 'email' | 'password'
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: 'Enter text...',
  disabled: false,
  maxLength: 255,
  type: 'text',
})

// Props are read-only
// props.modelValue = 'test' // ❌ Type error
</script>
```

### Emits with Strong Typing

```vue
<script setup lang="ts">
interface Props {
  value: string
}

defineProps<Props>()

const emit = defineEmits<{
  // Event name → payload tuple
  'update:value': [newValue: string]
  submit: [value: string, timestamp: number]
  error: [error: Error]
}>()

function updateValue(newValue: string) {
  emit('update:value', newValue)
}

function handleSubmit() {
  emit('submit', 'test', Date.now())
}
</script>
```

### defineModel (Vue 3.4+)

```vue
<!-- Parent: <MyInput v-model="message" /> -->

<script setup lang="ts">
// Single model
const model = defineModel<string>()

// ✅ Automatically handles v-model binding
// model.value gets/sets the prop and emits update:modelValue

// Usage in template
// <input v-model="model" />

// With modifiers (trim, number, lazy)
const modelWithTrim = defineModel({ default: '', get: (v) => v.trim() })
</script>
```

### defineExpose for Template Refs

```vue
<script setup lang="ts">
import { ref } from 'vue'

const inputEl = ref<HTMLInputElement | null>(null)
const internalValue = ref('')

function focus() {
  inputEl.value?.focus()
}

function clear() {
  internalValue.value = ''
}

// Only expose these methods
defineExpose({ focus, clear })
</script>

<!-- Parent -->
<template>
  <MyInput ref="inputRef" />
  <button @click="inputRef?.focus()">Focus</button>
</template>

<script setup>
const inputRef = ref()
</script>
```

### Computed vs Methods

**Use `computed` for**:

- Derived state that depends on reactive data
- Expensive calculations (Vue caches automatically)
- Properties accessed multiple times in template

```typescript
const fullName = computed(() => `${firstName.value} ${lastName.value}`)
```

**Use `methods` for**:

- Event handlers
- Side effects
- Actions unrelated to template rendering

```typescript
function handleClick() {
  /* do something */
}
function logMessage() {
  /* logging */
}
```

### Conditional Rendering with TypeScript

```vue
<script setup lang="ts">
import { ref } from 'vue'

type View = 'list' | 'detail' | 'edit'
const view = ref<View>('list')

const users = ref<User[]>([])
const selectedUser = ref<User | null>(null)

function selectUser(user: User) {
  selectedUser.value = user
  view.value = 'detail'
}
</script>

<template>
  <div>
    <!-- TypeScript guards type inside v-if block -->
    <ListView v-if="view === 'list'" @select="selectUser" />
    <DetailView v-else-if="view === 'detail' && selectedUser" :user="selectedUser" />
    <EditView v-else-if="view === 'edit' && selectedUser" :user="selectedUser" />
  </div>
</template>
```

### Full Component Example

```vue
<template>
  <form @submit.prevent="handleSubmit" class="form">
    <input
      v-model="formData.email"
      type="email"
      placeholder="Email"
      @blur="validateField('email')"
    />
    <span v-if="errors.email" class="error">{{ errors.email }}</span>

    <input
      v-model="formData.password"
      type="password"
      placeholder="Password"
      @blur="validateField('password')"
    />
    <span v-if="errors.password" class="error">{{ errors.password }}</span>

    <button :disabled="loading">{{ loading ? 'Logging in...' : 'Login' }}</button>
  </form>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'

interface LoginForm {
  email: string
  password: string
}

const emit = defineEmits<{
  'login-success': [user: User]
  'login-error': [error: Error]
}>()

const formData = reactive<LoginForm>({ email: '', password: '' })
const errors = reactive<Record<keyof LoginForm, string>>({ email: '', password: '' })
const loading = ref(false)

function validateField(field: keyof LoginForm) {
  if (formData[field].length === 0) {
    errors[field] = 'This field is required'
  } else {
    errors[field] = ''
  }
}

async function handleSubmit() {
  // Validate all
  Object.keys(formData).forEach((field) => validateField(field as keyof LoginForm))

  if (Object.values(errors).some((e) => e)) return

  loading.value = true
  try {
    const response = await fetch('/auth/login', {
      method: 'POST',
      body: JSON.stringify(formData),
    })
    const user = await response.json()
    emit('login-success', user)
  } catch (error) {
    emit('login-error', error as Error)
  } finally {
    loading.value = false
  }
}
</script>
```

---

## 4. Avoiding `any`: Type-Safe Props, Emits, Refs, and Provide/Inject

### Props Typing (Never `any`)

```typescript
// ❌ WRONG
const props = defineProps<{ user: any }>()

// ✅ CORRECT
interface User {
  id: string
  name: string
  email: string
}

const props = defineProps<{ user: User }>()
```

### Event Handler Typing

```vue
<script setup lang="ts">
interface Props {
  items: string[]
}

defineProps<Props>()

// ❌ WRONG
// function handleClick(e: any) { }

// ✅ CORRECT
function handleClick(e: MouseEvent) {
  console.log(e.clientX, e.clientY)
}

function handleInputChange(e: Event) {
  const target = e.target as HTMLInputElement
  console.log(target.value)
}

function handleKeyDown(e: KeyboardEvent) {
  if (e.key === 'Enter') {
    /* ... */
  }
}
</script>

<template>
  <button @click="handleClick">Click</button>
  <input @change="handleInputChange" />
  <input @keydown="handleKeyDown" />
</template>
```

### Template Refs Typing

```typescript
// ❌ WRONG
const inputRef = ref() // type is `any`

// ✅ CORRECT
const inputRef = ref<HTMLInputElement | null>(null)

const componentRef = ref<InstanceType<typeof MyComponent> | null>(null)

onMounted(() => {
  inputRef.value?.focus()
  componentRef.value?.triggerAction()
})
```

### Composable Return Typing

```typescript
// ❌ WRONG
export function useForm(initial: any) {
  return { data: ref(initial), errors: ref({}) }
}

// ✅ CORRECT
export function useForm<T extends Record<string, any>>(initial: T) {
  const data = ref<T>(initial)
  const errors = ref<Record<keyof T, string>>({} as any)

  const reset = () => {
    data.value = structuredClone(initial)
    Object.keys(errors.value).forEach((key) => {
      errors.value[key as keyof T] = ''
    })
  }

  return { data, errors, reset }
}
```

### Provide/Inject Typing

```typescript
// ❌ WRONG
const userKey = Symbol()
provide(userKey, {
  /* any user object */
})
const user = inject(userKey) // type is `any`

// ✅ CORRECT: Use strong typing
interface User {
  id: string
  name: string
}

// Create a typed inject key
type UserProvideKey = InjectionKey<User>
const userKey = Symbol() as UserProvideKey

// Parent provides with type
provide<User>(userKey, { id: '1', name: 'John' })

// Child injects with type
const user = inject<User>(userKey)
// user is now typed as User | undefined
```

### Emits Typing

```typescript
// ❌ WRONG
const emit = defineEmits(['click', 'submit'])

// ✅ CORRECT
const emit = defineEmits<{
  click: [x: number, y: number]
  submit: [data: FormData]
  error: [error: Error]
}>()

// Usage
emit('click', 10, 20)
emit('error', new Error('Something went wrong'))
```

---

## 5. Component Composition: Slots, Provide/Inject, v-model

### Slots Basics

**Named slots**:

```vue
<!-- Parent -->
<template>
  <Card>
    <template #header>
      <h1>Title</h1>
    </template>
    <template #default>
      <p>Content</p>
    </template>
    <template #footer>
      <button>Close</button>
    </template>
  </Card>
</template>

<!-- Card.vue -->
<template>
  <div class="card">
    <div class="card-header">
      <slot name="header"></slot>
    </div>
    <div class="card-body">
      <slot></slot>
      <!-- default slot -->
    </div>
    <div class="card-footer">
      <slot name="footer"></slot>
    </div>
  </div>
</template>
```

### Scoped Slots

```vue
<!-- Parent passes data to child slot -->
<template>
  <UserList v-slot="{ user, index }">
    <div :key="index">{{ index + 1 }}. {{ user.name }}</div>
  </UserList>
</template>

<!-- UserList.vue -->
<script setup lang="ts">
interface User {
  name: string
}

defineProps<{ users: User[] }>()
</script>

<template>
  <div>
    <slot v-for="(user, index) in users" :key="index" :user="user" :index="index"></slot>
  </div>
</template>
```

### Renderless Components (composition without markup)

```vue
<!-- ClickOutside.vue - handles click outside logic -->
<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

interface Props {
  element?: string
}

defineProps<Props>()

const emit = defineEmits<{
  'click-outside': [event: MouseEvent]
}>()

const rootEl = ref<HTMLElement | null>(null)

function handleClickOutside(e: MouseEvent) {
  if (rootEl.value && !rootEl.value.contains(e.target as Node)) {
    emit('click-outside', e)
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div ref="rootEl">
    <slot></slot>
  </div>
</template>

<!-- Usage -->
<template>
  <ClickOutside @click-outside="closeMenu">
    <button @click="isOpen = true">Menu</button>
    <div v-if="isOpen" class="menu">...</div>
  </ClickOutside>
</template>
```

### Provide/Inject Pattern

```typescript
// useDialogProvider.ts
import { provide, inject, ref, type InjectionKey } from 'vue'

export interface DialogContext {
  isOpen: boolean
  title: string
  open: (title: string) => void
  close: () => void
}

const dialogKey = Symbol() as InjectionKey<DialogContext>

export function useDialogProvider() {
  const isOpen = ref(false)
  const title = ref('')

  const context: DialogContext = {
    isOpen,
    title,
    open: (t: string) => {
      title.value = t
      isOpen.value = true
    },
    close: () => {
      isOpen.value = false
    }
  }

  provide(dialogKey, context)
  return context
}

export function useDialog() {
  const dialog = inject<DialogContext>(dialogKey)
  if (!dialog) throw new Error('useDialog must be used within DialogProvider')
  return dialog
}

// DialogProvider.vue
<template>
  <div>
    <slot></slot>
    <Teleport to="body" v-if="dialog.isOpen">
      <div class="dialog-overlay" @click="dialog.close">
        <div class="dialog" @click.stop>
          <h2>{{ dialog.title }}</h2>
          <slot name="dialog-content"></slot>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { useDialogProvider } from './useDialogProvider'
const dialog = useDialogProvider()
</script>

// Usage
<template>
  <DialogProvider>
    <ChildComponent />
  </DialogProvider>
</template>

<script setup lang="ts">
import { useDialog } from './useDialogProvider'

const dialog = useDialog()

function openDialog() {
  dialog.open('Confirm Action')
}
</script>
```

### Advanced v-model with Modifiers

```vue
<!-- Parent -->
<template>
  <MyInput v-model.trim.number="value" />
</template>

<!-- MyInput.vue -->
<script setup lang="ts">
interface Props {
  modelValue: string | number
  modelModifiers?: Record<string, boolean>
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: string | number]
}>()

function updateValue(e: Event) {
  let value: any = (e.target as HTMLInputElement).value

  if (props.modelModifiers?.trim) {
    value = value.trim()
  }

  if (props.modelModifiers?.number) {
    value = parseInt(value, 10) || value
  }

  emit('update:modelValue', value)
}
</script>

<template>
  <input :value="modelValue" @input="updateValue" />
</template>
```

---

## 6. Reactivity System: ref, reactive, toRef, toRefs, watch, watchEffect

### ref vs reactive

```typescript
// ref: for primitives, best for simple state
const count = ref(0) // Ref<number>
count.value++ // Must use .value in script

const name = ref('John') // Ref<string>
name.value = 'Jane'

// reactive: for objects, no .value needed in templates
const user = reactive({ name: 'John', age: 30 })
user.name = 'Jane' // Direct access, no .value

// ❌ DON'T mix: reactive with object in ref
const userRef = ref({ name: 'John' }) // awkward: userRef.value.name

// ✅ DO: Use reactive for objects
const user = reactive({ name: 'John' })
```

### toRef / toRefs: Extract Reactivity

```typescript
const user = reactive({ name: 'John', email: 'john@example.com' })

// ❌ WRONG: Destructuring loses reactivity
const { name, email } = user
// name and email are plain strings, not reactive

// ✅ CORRECT: Use toRefs
const { name, email } = toRefs(user)
// name and email are Ref<string>, remain reactive

// ✅ CORRECT: Use toRef for single property
const name = toRef(user, 'name')
// name is Ref<string>, tracks user.name changes
```

### watch vs watchEffect

**watch: explicit dependency tracking**

```typescript
const count = ref(0)
const multiplied = ref(0)

// Watch specific source
watch(count, (newVal, oldVal) => {
  multiplied.value = newVal * 2
})

// Watch multiple sources
watch([count, multiplied], ([c, m], [oldC, oldM]) => {
  console.log(`Count: ${c}, Multiplied: ${m}`)
})

// Watch object property with deep
watch(
  () => user.profile.name,
  (name) => {
    console.log(name)
  },
)

// Watch with options
watch(
  count,
  (newVal) => {
    console.log(newVal)
  },
  { immediate: true, deep: true, flush: 'post' },
)
```

**watchEffect: automatic dependency collection**

```typescript
const count = ref(0)
const multiplied = ref(0)

// Runs whenever reactive dependencies change
watchEffect(() => {
  // Automatically tracks count and multiplied
  console.log(`${count.value} * 2 = ${multiplied.value}`)
})

// No explicit dependency list needed
// More concise for simple watchers
```

### Shallow Reactivity

```typescript
// Regular reactive: deep (expensive for large objects)
const user = reactive({ name: 'John', profile: { bio: 'Dev' } })

// ✅ Shallow: only top-level properties tracked
const user = shallowReactive({ name: 'John', profile: { bio: 'Dev' } })
user.name = 'Jane' // triggers watcher
user.profile.bio = 'Designer' // DOES NOT trigger (shallow)

// ✅ Shallow ref: only .value change triggers
const data = shallowRef({ items: [1, 2, 3] })
data.value = { items: [4, 5, 6] } // triggers
data.value.items.push(7) // DOES NOT trigger
```

### Computed with Dependencies

```typescript
const firstName = ref('John')
const lastName = ref('Doe')

// Auto-dependency tracking
const fullName = computed(() => {
  return `${firstName.value} ${lastName.value}`
})

// With setter (for v-model)
const fullName = computed({
  get: () => `${firstName.value} ${lastName.value}`,
  set: (value: string) => {
    const [first, last] = value.split(' ')
    firstName.value = first
    lastName.value = last
  },
})
```

### Unref / isRef / isReactive / isReadonly

```typescript
import { unref, isRef, isReactive, isReadonly } from 'vue'

const count = ref(0)
const user = reactive({ name: 'John' })

unref(count) // Returns 0 (not the Ref object)
isRef(count) // true
isReactive(user) // true
isReadonly(readonly(user)) // true
```

---

## 7. Advanced Patterns: VueUse, Error Handling, Async State

### VueUse Common Composables

**useLocalStorage**:

```typescript
import { useLocalStorage } from '@vueuse/core'

const theme = useLocalStorage('theme', 'light')

theme.value = 'dark' // Auto-syncs to localStorage

// With custom serialization
const user = useLocalStorage('user', null, {
  serializer: {
    read: (v) => JSON.parse(v ?? 'null'),
    write: (v) => JSON.stringify(v),
  },
})
```

**useAsyncState**:

```typescript
import { useAsyncState } from '@vueuse/core'

const { state: users, isReady, isLoading, error, execute } = useAsyncState(
  async () => {
    const response = await fetch('/api/users')
    return response.json() as User[]
  },
  [] // initial state
)

// Execute manually
await execute()

// Usage
if (isLoading.value) return <div>Loading...</div>
if (error.value) return <div>Error: {error.value.message}</div>
return <div>{state.value.map(...)}</div>
```

**useEventListener**:

```typescript
import { useEventListener } from '@vueuse/core'

const count = ref(0)

// Automatically manages addEventListener/removeEventListener
useEventListener('click', () => {
  count.value++
})

// On specific element
const el = ref<HTMLElement>()
useEventListener(el, 'scroll', () => {
  console.log('scrolled')
})
```

**useDebounceFn**:

```typescript
import { useDebounceFn } from '@vueuse/core'

const search = ref('')

const debouncedSearch = useDebounceFn(async (query: string) => {
  const results = await fetch(`/api/search?q=${query}`).then((r) => r.json())
  return results
}, 300)

watch(search, (query) => {
  debouncedSearch(query)
})
```

### Error Handling Composable

```typescript
export function useAsyncTask<T>(fn: () => Promise<T>) {
  const data = ref<T | null>(null)
  const loading = ref(false)
  const error = ref<Error | null>(null)

  async function execute() {
    loading.value = true
    error.value = null
    try {
      data.value = await fn()
    } catch (e) {
      error.value = e instanceof Error ? e : new Error(String(e))
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    execute()
  })

  return {
    data: readonly(data),
    loading: readonly(loading),
    error: readonly(error),
    retry: execute,
  }
}

// Usage
const { data, loading, error, retry } = useAsyncTask(async () => {
  const response = await fetch('/api/users')
  if (!response.ok) throw new Error(`HTTP ${response.status}`)
  return response.json() as User[]
})
```

### Async State in Components

```vue
<script setup lang="ts">
import { ref, computed } from 'vue'

type Status = 'idle' | 'loading' | 'success' | 'error'

const formData = reactive({ email: '', password: '' })
const status = ref<Status>('idle')
const error = ref<string | null>(null)

async function handleLogin() {
  status.value = 'loading'
  error.value = null

  try {
    const response = await fetch('/auth/login', {
      method: 'POST',
      body: JSON.stringify(formData),
    })

    if (!response.ok) {
      throw new Error(`Login failed: ${response.statusText}`)
    }

    const user = await response.json()
    status.value = 'success'
    // Redirect or emit success
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Unknown error'
    status.value = 'error'
  }
}

const canSubmit = computed(() => {
  return formData.email && formData.password && status.value !== 'loading'
})
</script>

<template>
  <form @submit.prevent="handleLogin">
    <div v-if="status === 'error'" class="error-banner">{{ error }}</div>
    <input v-model="formData.email" type="email" />
    <input v-model="formData.password" type="password" />
    <button :disabled="!canSubmit">
      {{ status === 'loading' ? 'Logging in...' : 'Login' }}
    </button>
  </form>
</template>
```

---

## 8. TypeScript Strictness Configuration

### tsconfig.json for Desktop App

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": false,
    "esModuleInterop": true,

    // Strict mode (REQUIRED for production)
    "strict": true,
    "noImplicitAny": true,
    "strictNullChecks": true,
    "strictFunctionTypes": true,
    "strictBindCallApply": true,
    "strictPropertyInitialization": true,
    "noImplicitThis": true,
    "alwaysStrict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true,

    "resolveJsonModule": true,
    "isolatedModules": true,
    "moduleResolution": "bundler",
    "allowSyntheticDefaultImports": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true,
    "types": ["bun-types"],
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"]
    }
  },
  "include": ["src/**/*.ts", "src/**/*.vue", "src/**/*.d.ts"],
  "exclude": ["node_modules", "dist"]
}
```

---

## 9. SVG Workflow: vite-svg-loader (No Inline SVG)

**Rule**: Never embed raw `<svg>` markup in `.vue` files. Never use `v-html` to inject SVG strings. Put SVG files in `src/assets/icons/` and import them as Vue components.

### Setup (already configured in this project)

```typescript
// vite.main.config.ts
import svgLoader from 'vite-svg-loader'

export default defineConfig({
  plugins: [vue(), svgLoader({ defaultImport: 'component' })],
})
```

### TypeScript declaration (env.d.ts)

```typescript
/// <reference types="vite/client" />
/// <reference types="vite-svg-loader" />

declare module '*.svg' {
  import type { FunctionalComponent, SVGAttributes } from 'vue'
  const src: FunctionalComponent<SVGAttributes>
  export default src
}

declare module '*.svg?url' {
  const src: string
  export default src
}

declare module '*.svg?raw' {
  const src: string
  export default src
}
```

### Usage in components

```vue
<script setup lang="ts">
// ✅ CORRECT — import SVG as Vue component
import CloseIcon from '@desktop/assets/icons/close.svg'
import SettingsIcon from '@desktop/assets/icons/settings.svg'
</script>

<template>
  <!-- Use like any Vue component; pass class/style props -->
  <CloseIcon class="h-4 w-4 text-current" />
  <SettingsIcon :style="{ color: iconColor, fontSize: '1.25rem' }" />
</template>
```

```typescript
// ❌ WRONG — inline SVG string in template
// <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">...</svg>

// ❌ WRONG — v-html with SVG string
// <span v-html="platformIconSvg" />

// ❌ WRONG — returning SVG string from function
// function getPlatformIcon(platform: string): string { return '<svg>...' }
```

### Icon organization

```
src/assets/icons/
├── platforms/
│   ├── twitch.svg
│   ├── youtube.svg
│   └── kick.svg
├── ui/
│   ├── close.svg
│   ├── settings.svg
│   ├── add.svg
│   └── chevron-down.svg
└── actions/
    ├── send.svg
    └── copy.svg
```

### Reusable Icon wrapper (optional, for size/color control)

```vue
<!-- src/views/main/components/ui/Icon.vue -->
<script setup lang="ts">
import type { FunctionalComponent, SVGAttributes } from 'vue'

interface Props {
  icon: FunctionalComponent<SVGAttributes>
  size?: string
  color?: string
}

const props = withDefaults(defineProps<Props>(), {
  size: '1em',
  color: 'currentColor',
})
</script>

<template>
  <component
    :is="props.icon"
    :style="{ width: size, height: size, color }"
    class="inline-block flex-shrink-0"
  />
</template>

<!-- Usage -->
<!-- <Icon :icon="CloseIcon" size="1.25rem" /> -->
```

---

## 10. Pinia: Setup for TwirChat Desktop

Pinia is installed in `packages/desktop`. Register it once in `src/views/main/main.ts`.

### Registration (main.ts)

```typescript
// src/views/main/main.ts
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { Electroview } from 'electrobun/view'
import App from './App.vue'

export const rpc = Electroview.defineRPC<TwirChatRPCSchema>({ ... })
new Electroview({ rpc })

const app = createApp(App)
app.use(createPinia())
app.mount('#app')
```

### Store conventions for TwirChat

All stores live in `src/views/main/stores/`. Use **setup stores** (not Options API stores) for TypeScript compatibility.

```typescript
// src/views/main/stores/accounts.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Account } from '@twirchat/shared/types'
import { rpc } from '../main'

export const useAccountsStore = defineStore('accounts', () => {
  const accounts = ref<Account[]>([])
  const loading = ref(false)

  const twitchAccounts = computed(() => accounts.value.filter((a) => a.platform === 'twitch'))

  async function loadAccounts() {
    loading.value = true
    try {
      accounts.value = await rpc.request.getAccounts()
    } finally {
      loading.value = false
    }
  }

  function addAccount(account: Account) {
    accounts.value.push(account)
  }

  return { accounts, loading, twitchAccounts, loadAccounts, addAccount }
})
```

### Accessing stores correctly

```typescript
// ✅ CORRECT — direct property access
const accountsStore = useAccountsStore()
accountsStore.loadAccounts()
console.log(accountsStore.accounts)

// ✅ CORRECT — storeToRefs for destructured reactive access
import { storeToRefs } from 'pinia'
const { accounts, loading } = storeToRefs(accountsStore)

// ❌ WRONG — destructuring breaks reactivity
const { accounts } = useAccountsStore() // NOT reactive!
```

### What belongs in a Pinia store vs a composable

| Use Pinia store when                         | Use a composable when                     |
| -------------------------------------------- | ----------------------------------------- |
| State shared across multiple components      | State local to one component tree         |
| State needs to persist through navigation    | State is component lifecycle-scoped       |
| Multiple components read/write the same data | Logic is reusable but not globally shared |
| Devtools debugging needed                    | Simple utility logic                      |

**TwirChat Pinia stores to create:**

- `useAccountsStore` — loaded accounts from RPC
- `useSettingsStore` — app settings (from RPC, locally cached)
- `useChannelStatusStore` — live/offline status per channel
- `useChatAppearanceStore` — font size, colors, chat preferences

---

## 11. Composable Patterns for TwirChat

Common composables to extract from components.

### useRpcListener — subscribe to bun → view messages

```typescript
// src/views/main/composables/useRpcListener.ts
import { onUnmounted } from 'vue'
import { rpc } from '../main'
import type { TwirChatRPCSchema } from '@desktop/shared/rpc'

type BunMessages = TwirChatRPCSchema['bun']['messages']

export function useRpcListener<K extends keyof BunMessages>(
  event: K,
  handler: (payload: BunMessages[K]) => void,
) {
  rpc.on[event](handler)

  onUnmounted(() => {
    rpc.off?.[event]?.(handler)
  })
}

// Usage
useRpcListener('chat_message', (msg) => {
  messages.value.push(msg)
})
```

### usePolling — generic polling with cleanup

```typescript
// src/views/main/composables/usePolling.ts
import { ref, onUnmounted } from 'vue'

export function usePolling(fn: () => Promise<void>, intervalMs: number) {
  const isRunning = ref(false)
  let timer: ReturnType<typeof setInterval> | null = null

  async function tick() {
    if (isRunning.value) return
    isRunning.value = true
    try {
      await fn()
    } finally {
      isRunning.value = false
    }
  }

  function start() {
    void tick()
    timer = setInterval(() => void tick(), intervalMs)
  }

  function stop() {
    if (timer !== null) {
      clearInterval(timer)
      timer = null
    }
  }

  onUnmounted(stop)

  return { start, stop, isRunning }
}

// Usage (in ChatList.vue)
const { start } = usePolling(async () => {
  const statuses = await rpc.request.getChannelsStatus()
  channelStatuses.value = statuses
}, 30_000)

onMounted(start)
```

### usePlatformMeta — centralized platform colors/icons

```typescript
// src/views/main/composables/usePlatformMeta.ts
import type { Platform } from '@twirchat/shared/types'
import TwitchIcon from '@desktop/assets/icons/platforms/twitch.svg'
import YoutubeIcon from '@desktop/assets/icons/platforms/youtube.svg'
import KickIcon from '@desktop/assets/icons/platforms/kick.svg'
import type { FunctionalComponent, SVGAttributes } from 'vue'

interface PlatformMeta {
  color: string
  icon: FunctionalComponent<SVGAttributes>
  label: string
}

const META: Record<Platform, PlatformMeta> = {
  twitch: { color: '#9147ff', icon: TwitchIcon, label: 'Twitch' },
  youtube: { color: '#ff0000', icon: YoutubeIcon, label: 'YouTube' },
  kick: { color: '#53fc18', icon: KickIcon, label: 'Kick' },
}

export function usePlatformMeta(platform: Platform): PlatformMeta {
  return META[platform]
}

// Usage
const { color, icon: PlatformIcon, label } = usePlatformMeta(props.platform)
// <PlatformIcon class="h-4 w-4" />
```

---

## 12. Anti-Patterns in This Codebase (Do Not Repeat)

Patterns found in existing code that **must not be repeated** in new code:

| Anti-Pattern                                     | Why It's Wrong                                         | Correct Approach                            |
| ------------------------------------------------ | ------------------------------------------------------ | ------------------------------------------- |
| Inline `<svg>` in templates                      | Hard to reuse, bloats components, no SVGO optimization | Import `.svg` file via vite-svg-loader      |
| `v-html` with SVG strings                        | XSS risk, no TS types, not composable                  | SVG as Vue component                        |
| `platformColor()` duplicated in every file       | Copy-paste debt, inconsistent                          | `usePlatformMeta()` composable              |
| `rpc.request.*` called directly in `setup()`     | No cleanup, no error handling, can't be tested         | Wrap in composable with error/loading state |
| `setInterval` in component `onMounted`           | Memory leak if not cleaned up                          | `usePolling()` composable with auto cleanup |
| Logic in component that spans >50 lines          | Untestable, hard to read                               | Extract to `use*` composable                |
| Direct mutation of layout state                  | Bypasses reactivity contract                           | Actions in Pinia store                      |
| `const { x } = useStore()` without `storeToRefs` | Breaks reactivity silently                             | `const { x } = storeToRefs(useStore())`     |
| Importing from `src/store/` in views             | Crashes at runtime (SQLite in browser)                 | Use `rpc.request.*`                         |

---

## 13. Checklists

### Component Shipping Checklist

- ✅ All props have type definitions (no `any`)
- ✅ All emits have event signatures
- ✅ Event handlers typed (MouseEvent, KeyboardEvent, etc.)
- ✅ Template refs typed with `HTMLElement | null`
- ✅ Composables used have typed return values
- ✅ `withDefaults` used for optional props
- ✅ `storeToRefs` used when destructuring store
- ✅ `onUnmounted` cleanup for subscriptions/listeners
- ✅ Error states handled (loading, error, empty states)
- ✅ No inline `<svg>` — use imported SVG component
- ✅ No `v-html` for icons or SVG strings
- ✅ `bun run fix` (lint + format) passed
- ✅ No `as any` casts (only `as SpecificType` when necessary)
- ✅ No direct import from `src/store/` in views
- ✅ No `rpc.request.*` calls in component body — wrap in composable

### Composable Audit Checklist

- ✅ Starts with `use*` naming convention
- ✅ Returns object, not reactive root
- ✅ Generic types for reusability (if applicable)
- ✅ `onUnmounted` cleanup (events, timers, RPC listeners)
- ✅ `readonly()` for exposed state user shouldn't modify
- ✅ TypeScript: all parameters and returns typed
- ✅ Side effects isolated (fetch, RPC) in `onMounted`/`onUnmounted`
- ✅ Uses `toRefs()` when returning reactive objects
- ✅ Error handling (try/catch, error ref)
- ✅ `bun run fix` passed

---

## 14. Resources

- **Official Vue 3 Docs**: https://vuejs.org
- **Official Pinia Docs**: https://pinia.vuejs.org
- **VueUse Library**: https://vueuse.org (composable patterns)
- **vite-svg-loader**: https://github.com/jpkleemans/vite-svg-loader
- **Vue Language Tools**: https://github.com/johnsoncodehk/volar (TypeScript support)

---

## Summary: Mental Model for TwirChat Frontend

**State Management**:

- Global app state → Pinia store (`useAccountsStore()`, `useSettingsStore()`)
- Component local state → `ref()` or `reactive()`
- Derived state → `computed()`
- Bun-side persistence → always via `rpc.request.*`

**Side Effects**:

- Polling → `usePolling()` composable (auto-cleanup)
- RPC listeners → `useRpcListener()` composable (auto-cleanup)
- Data fetching → composable with loading/error state
- Cleanup → `onUnmounted()` always

**Typing**:

- Props → interface + `defineProps<Props>()`
- Emits → type def + `defineEmits<{ event: [args] }>()`
- Refs → `ref<Type>()`
- Stores → access directly or `storeToRefs()` for destructuring

**Composables**:

- Name: `use*`
- Return: object with refs/computed, NOT reactive root
- Cleanup: `onUnmounted()`

**SVGs**:

- Place in `src/assets/icons/`
- Import as Vue component via vite-svg-loader
- Never inline `<svg>` or use `v-html` for icons

**Components**:

- `<script setup lang="ts">`
- Type everything (no `any`)
- `withDefaults()` for optional props
- Handle loading/error states
- No direct store imports from `src/store/`
