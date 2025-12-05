### **Project Development Guidelines for AI Assistants (GitHub Copilot)**

This document outlines the core conventions, technologies, and patterns used in this project. Please
adhere to these guidelines strictly to maintain code consistency and quality.

### **1. General Project Context**

* **Structure:** This is a monorepo.
	* **`frontend/dashboard`**: The main dashboard application (Vue 3 + Vite).
	* **`web`**: The public-facing website (Nuxt 3).
	* **`apps/api-gql`**: The main backend service (Go) serving GraphQL and HTTP APIs.
	* **`libs`**: Shared Go libraries.
* **Package Manager & Runtime:** We use **Bun** for all JavaScript/TypeScript package management,
	script execution, and as the runtime. Use `bun install`, `bun add`, and `bun run` commands.
* **Primary Technologies:**
	* **Frontend (Dashboard):** Vue 3, TypeScript, Vite, Tailwind CSS, vee-validate, zod,
		lucide-vue-next, shadcn-vue, TanStack Query, Urql.
	* **Web (Public Site):** Nuxt 3, TypeScript, Tailwind CSS, shadcn-nuxt, Pinia, Urql.
	* **Backend:** Go (Golang), pgx (PostgreSQL driver), gqlgen (GraphQL).
	* **Tooling:** Bun, Docker.

---

### **2. Vue.js Frontend Development (Dashboard & Web)**

#### **2.1. Component Structure & Syntax**

* **Composition API:** Always use the Composition API.
* **Script Setup:** All Single File Components (SFCs) **must** use the `<script setup lang="ts">`
	syntax. Do not use the `setup()` function within the `export default` block.
* **Type Definitions:** Use `defineProps`, `defineEmits`, and `defineSlots` with explicit TypeScript
	types for clear, type-safe component interfaces.

**Example:**

```vue

<script setup lang="ts">
	import { computed } from 'vue';

	// Use interface or type for props definition
	interface Props {
		title: string;
		isActive?: boolean;
		items: string[];
	}

	const props = withDefaults(defineProps<Props>(), {
		isActive: false,
	});

	const emit = defineEmits<{
		(e: 'itemSelected', item: string): void;
		(e: 'closed'): void;
	}>();

	const handleItemClick = (item: string) => {
		emit('itemSelected', item);
	};

	const titleDisplay = computed(() => props.title.toUpperCase());
</script>

<template>
	<!-- Component template here -->
</template>
```

#### **2.2. Component Imports**

* **File Extension:** Always include the `.vue` file extension when importing Vue components. This
	improves clarity and avoids potential bundler configuration issues.

**Correct:**

```typescript
import UserProfile from './components/UserProfile.vue';
import AppHeader from '@/components/layout/AppHeader.vue';
```

**Incorrect:**

```typescript
import UserProfile from './components/UserProfile';
```

#### **2.3. Existing Components**

* **Delete Confirmation:** The project already has a standardized delete confirmation component. *
	*Do not create a new one.** When you need to confirm a delete action, import and use the existing
	component.
	* Assume its import path is something like:
		`import DeleteConfirmationDialog from '@/components/shared/DeleteConfirmationDialog.vue'`
	* Utilize its props and events as defined in its implementation.

---

### **3. Forms with `vee-validate` and `zod`**

We use `vee-validate` with `zod` for schema-based validation. Follow this pattern precisely.

* **Schema Library:** Use **zod** to define validation schemas.
* **Adapter:** Use the `@vee-validate/zod` library to connect `zod` schemas to `vee-validate`.
* **Hook, not Component:** Use the `useForm` hook from `vee-validate`.
* **Native `<form>` Element:** Bind your submission logic to a native HTML `<form>` element's
	`@submit` event. **DO NOT use the `<Form>` component provided by `vee-validate`**.
* **Fields:**
	* For standard text inputs, textareas, etc., use `v-slot="{ componentField }"` on your field
		wrapper and bind `v-bind="componentField"` to the input element.
	* For switches, checkboxes, and custom toggle components, use `v-slot="{ value, handleChange }"`
		to manage state.
* **Error Messages:** Always include the `FormMessage` component immediately after a form field to
	display validation errors.

**Example Form Structure:**

```vue

<script setup lang="ts">
	import { useForm } from 'vee-validate';
	import { toTypedSchema } from '@vee-validate/zod';
	import { z } from 'zod';
	import FormMessage from '@/components/ui/form/FormMessage.vue';
	import Input from '@/components/ui/input/Input.vue'; // Example custom input
	import Switch from '@/components/ui/switch/Switch.vue'; // Example custom switch

	const formSchema = z.object({
		name: z.string().min(2, 'Name must be at least 2 characters.'),
		email: z.string().email('Must be a valid email.'),
		subscribe: z.boolean().default(false),
	});

	const { handleSubmit, defineField } = useForm({
		validationSchema: toTypedSchema(formSchema),
	});

	const [name, nameAttrs] = defineField('name');
	const [email, emailAttrs] = defineField('email');

	const onSubmit = handleSubmit(values => {
		console.log('Form submitted:', values);
		// API call logic here
	});
</script>

<template>
	<form @submit="onSubmit" class="space-y-4">
		<!-- Standard Input Field using defineField -->
		<div>
			<label for="name">Name</label>
			<Input id="name" v-model="name" v-bind="nameAttrs" type="text" placeholder="Your Name"/>
			<FormMessage name="name"/>
		</div>

		<!-- Standard Input Field using defineField -->
		<div>
			<label for="email">Email</label>
			<Input id="email" v-model="email" v-bind="emailAttrs" type="email"
						 placeholder="email@example.com"
			/>
			<FormMessage name="email"/>
		</div>

		<!-- v-slot usage for a Switch/Checkbox -->
		<div class="flex items-center space-x-2">
			<VeeField name="subscribe" v-slot="{ value, handleChange }">
				<Switch id="subscribe" :checked="value" @update:checked="handleChange"/>
				<label for="subscribe">Subscribe to newsletter</label>
			</VeeField>
			<FormMessage name="subscribe"/>
		</div>

		<button type="submit">Submit</button>
	</form>
</template>
```

---

### **4. Iconography**

* **Primary Library:** **Always use `lucide-vue-next` for icons.** It is the project's standard.
* **Fallback:** Only if a specific icon is absolutely not available in Lucide should you consider
	using another library or a local SVG file. This should be a rare exception.
* **Usage:**
	* **Dashboard (`frontend`):** Import icons by name from the library.
	* **Web (`web`):** Use the `<Icon />` Nuxt component, and pass the name of the icon as
		`name="lucide:icon-name"`.

**Example Dashboard:**

```typescript
import { User, Mail, CheckCircle2 } from 'lucide-vue-next';
```

```vue

<template>
	<button class="btn">
		<User class="h-4 w-4 mr-2"/>
		Profile
	</button>
</template>
```

**Example Web:**

```vue

<template>
	<Icon name="lucide:user" class="h-4 w-4 mr-2"/>
</template>
```

---

### **5. Styling with Tailwind CSS**

* **Utility-First:** All styling must be done using Tailwind CSS utility classes directly in the
	`<template>` block. Avoid writing custom CSS in `<style>` blocks unless absolutely necessary for a
	complex, non-reusable scenario.
* **Project Configuration:** Adhere strictly to the project's `tailwind.config.js`.
	* **Colors:** Use the defined theme colors (e.g., `bg-primary`, `text-accent`,
		`border-destructive`). Do not use arbitrary hex codes or default Tailwind colors if custom ones
		are defined.
	* **Spacing & Sizing:** Use the defined spacing scale (e.g., `p-4`, `m-8`, `w-32`) instead of
		arbitrary values like `p-[15px]`.
	* **Component Classes:** If we use a library like `shadcn-vue` or have our own `@apply` directives
		for component base styles, be aware of and use them.

---

### **6. Dashboard Integrations Page Architecture**

The integrations page uses a **unified GraphQL query** pattern to fetch all integration data in a
single request, optimizing network usage and improving user experience.

#### **6.1. Unified Query Pattern**

* **Single Query File:** All integrations page data is fetched via a unified query in
	`frontend/dashboard/src/api/integrations/integrations-page.ts`.
* **Why:** This approach allows fetching data for all integrations (Discord, Spotify, LastFM,
	Valorant, DonationAlerts, etc.) in a single GraphQL request, which is significantly more efficient
	than making separate requests per integration.
* **Composable:** Use `useIntegrationsPageData()` to access the unified data. It provides computed
	refs for each integration's data.

#### **6.2. Adding or Refactoring Integrations**

When creating a new integration or refactoring an existing one to use GraphQL:

1. **Add fields to the unified query** in `integrations-page.ts`:
	 ```typescript
	 const IntegrationsPageQuery = graphql(`
		 query IntegrationsPageData {
			 # ... existing fields ...

			 # New integration
			 myNewIntegrationData {
				 enabled
				 userName
				 avatar
			 }
			 myNewIntegrationAuthLink
		 }
	 `)
	 ```

2. **Add computed refs** for the new integration data:
	 ```typescript
	 // MyNewIntegration
	 const myNewIntegrationData = computed(() => query.data.value?.myNewIntegrationData ?? null)
	 const myNewIntegrationAuthLink = computed(() => query.data.value?.myNewIntegrationAuthLink ?? null)
	 ```

3. **Export the new computed refs** in the return statement.

4. **Use the unified data in components** instead of creating separate queries:
	 ```typescript
	 const integrationsPage = useIntegrationsPageData()
	 // Access via integrationsPage.myNewIntegrationData
	 ```

#### **6.3. Mutations**

* Mutations (login, logout, update, etc.) should still be defined separately in dedicated files or
	in `integrations.ts`.
* Use `integrationsPageCacheKey` to invalidate the unified query cache after mutations:
	```typescript
	const myMutation = () =>
		useMutation(
			graphql(`mutation MyMutation { ... }`),
			[integrationsPageCacheKey]
		)
	```

---

### **7. Go (Golang) Backend**

* **Code Style:** Follow standard Go formatting (`gofmt`/`goimports`).
* **Project Structure:**
	* **`apps/api-gql`**: Main API service.
		* `internal/delivery/gql`: GraphQL resolvers.
		* `internal/delivery/http`: HTTP handlers.
		* `internal/services`: Business logic layer.
		* `internal/entity`: Domain models.
	* **`libs/repositories`**: Data access layer.
* **Repositories:**
	* Always use **pgx** implementations.
	* Located in `libs/repositories/{repository_name}/pgx/pgx.go`.
	* **NEVER** use GORM or other ORMs.
	* Repository should return entity written in `libs/entities/{entity_name}/entity.go` file.
		* For new models created in repositories, or when editing some repository, you should
			create/update model inside repository, do not create separate file for model. to include
			`isNil` property, and `IsNil` method to check if the model is
			empty, also create `var Nil = &Model{}` to represent an empty model. Example:

```go
type SomeModel struct {
ID            string
ChannelID     string

isNil bool
}

func (c SomeModel) IsNil() bool {
return c.isNil
}

var Nil = SomeModel{
isNil: true,
}

```

* **Mappers:**
	* When creating new services (e.g., in `api-gql`), always create an entity mapper.
	* Data flow: `Model (DB)` -> `Entity (Domain)` -> `DTO (GraphQL/HTTP)`.
* **GraphQL Generation:**
	* After updating GraphQL schemas (`.graphql` files), run `bun cli build gql` to regenerate
		resolvers.
	* After regeneration, refresh your data (re-read Golang files) to pick up changes.
* **Error Handling:**
	* Use `fmt.Errorf` with `%w` for wrapping errors.
	* Create custom error types if needed for specific domain error handling.
* **Entities**
	* Write entities in `libs/entities/{entity_name}/entity.go` file.
	* Entities should contain only domain logic and validation.
		* Avoid dependencies on other layers (e.g., repositories, services).
	* use Nil thing
