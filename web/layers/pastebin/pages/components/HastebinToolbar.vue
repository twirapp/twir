<script setup lang="ts">
import { Icon } from '#components'
import { usePasteStore } from '#layers/pastebin/stores/pasteStore'

const emit = defineEmits<{
	save: []
	new: []
	copy: []
}>()

const router = useRouter()

const userStore = useAuth()
await callOnce(UserStoreKey, () => userStore.getUserData())

const pasteStore = usePasteStore()
const { currentPaste, editableContent } = storeToRefs(pasteStore)

const requestUrl = useRequestURL()

const canSave = computed(() => {
	if (currentPaste.value) {
		return false
	}

	return editableContent.value?.length
})

const canCopy = computed(() => {
	return !!currentPaste.value
})

const buttons = computed(() => {
	return [
		{
			name: 'Save',
			icon: 'lucide:save',
			disabled: !canSave.value,
			tooltip: 'Save',
			onClick: () => emit('save'),
		},
		{
			name: 'New',
			icon: 'lucide:file-plus',
			tooltip: 'New',
			onClick: () => emit('new'),
		},
		{
			name: 'Copy',
			icon: 'lucide:copy',
			disabled: !canCopy.value,
			tooltip: 'Copy',
			onClick: () => emit('copy'),
		},
		{
			name: 'Raw',
			icon: 'lucide:text',
			disabled: !currentPaste.value,
			tooltip: 'Text',
			href: currentPaste.value?.id ? `${requestUrl.origin}/h/${currentPaste.value?.id}/raw` : undefined,
		},
		{
			name: 'Profile',
			icon: 'lucide:user',
			tooltip: 'Profile',
			disabled: !userStore.user,
			href: `${requestUrl.origin}/h/profile`,
			onClick: () => router.push('/h/profile'),
		},
		{
			name: 'Api',
			icon: 'lucide:code',
			tooltip: 'Api',
			href: `${requestUrl.origin}/api/docs#tag/pastebin`,
		},
	]
})
</script>

<template>
	<div class="flex flex-row gap-2 items-center fixed top-2 right-4 h-14 p-2 bg-gray-500/50 rounded-md">
		<UiButton
			v-for="button of buttons"
			:key="button.name"
			variant="link"
			size="icon"
			:disabled="button.disabled ?? false"
			:as="button.href ? 'a' : 'button'"
			:href="button.href"
			class="flex flex-col items-center justify-center"
			@click="button.onClick"
		>
			<Icon class="icon" :name="button.icon" />
			<span class="text-xs">{{ button.name }}</span>
		</UiButton>
	</div>
</template>

<style scoped>
.icon {
  @apply size-6 cursor-pointer text-gray-300;
}

.icon:disabled {
  @apply opacity-50;
}

.icon:hover {
  @apply text-gray-100;
}
</style>
