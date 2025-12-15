<script setup lang="ts">
import { PencilIcon, TrashIcon } from 'lucide-vue-next'
import { ref } from 'vue'

import KeywordDialog from './keywords-dialog.vue'

import type { KeywordResponse } from '@/api/keywords'

import { useKeywordsApi } from '@/api/keywords'
import ActionConfirm from '@/components/ui/action-confirm.vue'
import { Button } from '@/components/ui/button'
import { Switch } from '@/components/ui/switch'

const props = defineProps<{ keyword: KeywordResponse }>()
const showDelete = ref(false)

const keywordsApi = useKeywordsApi()
const updateMutation = keywordsApi.useMutationUpdateKeyword()
const removeMutation = keywordsApi.useMutationRemoveKeyword()

function toggleEnabledGreetings() {
	updateMutation.executeMutation({
		id: props.keyword.id,
		opts: { enabled: !props.keyword.enabled },
	})
}

function deleteGreetings() {
	removeMutation.executeMutation({ id: props.keyword.id })
}
</script>

<template>
	<div class="flex items-center gap-2">
		<Switch :model-value="keyword.enabled" @update:model-value="toggleEnabledGreetings" />

		<KeywordDialog :keyword="keyword">
			<template #dialog-trigger>
				<Button variant="secondary" size="icon">
					<PencilIcon class="h-4 w-4" />
				</Button>
			</template>
		</KeywordDialog>

		<Button variant="destructive" size="icon" @click="showDelete = true">
			<TrashIcon class="h-4 w-4" />
		</Button>
	</div>

	<ActionConfirm v-model:open="showDelete" @confirm="deleteGreetings" />
</template>
