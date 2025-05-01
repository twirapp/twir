<script setup lang="ts">
import type { PasteBinOutputDto } from '@twir/api/openapi'

const props = defineProps<{
	item?: PasteBinOutputDto
}>()

const router = useRouter()
const requestUrl = useRequestURL()
const api = useOapi()

const content = ref<string>()
const isEditMode = ref(!props.item)
const textareaRef = ref<HTMLTextAreaElement>()

async function create() {
	if (!content.value?.length) return

	const req = await api.v1.pastebinCreate({
		content: content.value,
	})
	if (req.error) {
		throw req.error
	}

	await router.push(`/h/${req.data?.id}`)
}

async function duplicate() {
	if (!props.item) return

	content.value = props.item.content
	isEditMode.value = true

	// Focus the textarea after it's rendered
	nextTick(() => {
		textareaRef.value?.focus()
	})
}

function newPaste() {
	content.value = ''
	isEditMode.value = true

	// Focus the textarea after it's rendered
	nextTick(() => {
		textareaRef.value?.focus()
	})
}

// Focus textarea when component is mounted if in edit mode
onMounted(() => {
	if (isEditMode.value) {
		textareaRef.value?.focus()
	}
})

// Watch for changes to isEditMode and focus textarea when it becomes true
watch(isEditMode, (newValue) => {
	if (newValue) {
		nextTick(() => {
			textareaRef.value?.focus()
		})
	}
})

const topButtons = computed<{
	name: string
	icon: string
	disabled?: boolean
	tooltip: string
	href?: string
	onClick?: () => void
}[]>(() => {
			const isItemProped = !!props.item && !isEditMode.value

			return [
				{
					name: 'save',
					icon: 'lucide:save',
					disabled: isItemProped || !content.value?.length,
					tooltip: 'Save',
					onClick: create,
				},
				{
					name: 'new',
					icon: 'lucide:file-plus',
					tooltip: 'New',
					onClick: newPaste,
				},
				{
					name: 'copy',
					icon: 'lucide:copy',
					disabled: !isItemProped,
					tooltip: 'Copy',
					onClick: duplicate,
				},
				{
					name: 'text',
					icon: 'lucide:text',
					disabled: !isItemProped,
					tooltip: 'Text',
					href: `${requestUrl.origin}/h/${props.item?.id}/raw`,
				},
			]
		})
</script>

<template>
	<div class="h-full w-full p-4 relative">
		<div
			class="flex flex-row gap-2 items-center absolute top-2 right-4 h-10 p-2 bg-gray-500/50 rounded-md"
		>
			<UiButton
				v-for="button of topButtons"
				:key="button.name"
				variant="link"
				size="icon"
				:disabled="button.disabled"
				:as="button.href ? 'a' : 'button'"
				:href="button.href"
				@click="button.onClick"
			>
				<Icon class="icon" :name="button.icon" />
			</UiButton>
		</div>
		<Shiki v-if="props.item && !isEditMode" :code="props.item.content" as="div" class="h-full" />
		<textarea
			v-else
			ref="textareaRef"
			v-model="content"
			class="h-full w-full p-2 bg-transparent outline-none rounded-md input"
		/>
	</div>
</template>

<style scoped>
:deep(.shiki) {
	@apply h-full w-full
}

:deep(code) {
	counter-reset: step;
	counter-increment: step 0;
}

:deep(pre code) {
	font-family: 'JetBrains Mono';
}

:deep(code .line::before) {
	content: counter(step);
	counter-increment: step;
	width: 1rem;
	margin-right: 1.5rem;
	display: inline-block;
	text-align: right;
	color: rgba(115, 138, 148, .4)
}

.icon {
	@apply size-6 cursor-pointer text-gray-300
}

.icon:disabled {
	@apply opacity-50
}

.icon:hover {
	@apply text-gray-100
}

.input {
	font-family: 'JetBrains Mono';
}
</style>
