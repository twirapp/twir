<script setup lang="ts">
import { BarChart } from 'lucide-vue-next'
import { NA, NScrollbar } from 'naive-ui'

import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import {
	Dialog,
	DialogTrigger,
} from '@/components/ui/dialog'
import CommunityEmotesDetailsContent from '@/features/community-emotes-statistic/ui/community-emotes-details-content.vue'

defineEmits<{
	select: []
}>()
</script>

<template>
	<Dialog>
		<DialogTrigger as-child>
			<NA class="flex text-xs cursor-pointer items-center" @click="$emit('select')">
				<span>Usage Stats</span>
				<BarChart class="size-4 ml-1" />
			</NA>
		</DialogTrigger>
		<DialogOrSheet class="p-0">
			<NScrollbar style="max-height: 85vh" trigger="none">
				<CommunityEmotesDetailsContent />
			</NScrollbar>
		</DialogOrSheet>
	</Dialog>
</template>
