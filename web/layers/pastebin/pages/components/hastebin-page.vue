<script setup lang="ts">
import type { PasteBinOutputDto } from '@twir/api/openapi'

const props = defineProps<{
	item?: PasteBinOutputDto
}>()

const router = useRouter()
const requestUrl = useRequestURL()
const api = useOapi()

const content = ref<string>()

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

async function duplicate(newContent: string) {

}

const topButtons = computed<{
	name: string
	icon: string
	disabled?: boolean
	tooltip: string
	href?: string
	onClick?: () => void
}[]>(() => {
			const isItemProped = !!props.item

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
				},
				{
					name: 'copy',
					icon: 'lucide:copy',
					disabled: !isItemProped,
					tooltip: 'Copy',
					onClick: () => duplicate(props.item!.content),
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
		<Shiki v-if="item" :code="item.content" as="div" class="h-full" />
		<textarea v-else v-model="content" class="h-full w-full p-2 bg-transparent outline-none rounded-md input" />
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
