<script setup lang="ts">
import { BarChart } from 'lucide-vue-next'

import DialogOrSheet from '#layers/dashboard/components/dialog-or-sheet.vue'



import CommunityEmotesDetailsContent from '~/features/community-emotes-statistic/ui/community-emotes-details-content.vue'

defineEmits<{
	select: []
}>()
</script>

<template>
	<UiDialog>
		<UiDialogTrigger as-child>
			<UiButton variant="link" class="flex text-xs h-auto p-0 items-center" @click="$emit('select')">
				<span>Usage Stats</span>
				<BarChart class="size-4 ml-1" />
			</UiButton>
		</UiDialogTrigger>
		<DialogOrSheet class="p-0">
			<UiScrollArea class="max-h-[85vh]">
				<CommunityEmotesDetailsContent />
			</UiScrollArea>
		</DialogOrSheet>
	</UiDialog>
</template>
