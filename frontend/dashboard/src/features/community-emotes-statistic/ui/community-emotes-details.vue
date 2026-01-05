<script setup lang="ts">
import { BarChart } from 'lucide-vue-next'

import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogTrigger,
} from '@/components/ui/dialog'
import { ScrollArea } from '@/components/ui/scroll-area'
import CommunityEmotesDetailsContent from '@/features/community-emotes-statistic/ui/community-emotes-details-content.vue'

defineEmits<{
	select: []
}>()
</script>

<template>
	<Dialog>
		<DialogTrigger as-child>
			<Button variant="link" class="flex text-xs h-auto p-0 items-center" @click="$emit('select')">
				<span>Usage Stats</span>
				<BarChart class="size-4 ml-1" />
			</Button>
		</DialogTrigger>
		<DialogOrSheet class="p-0">
			<ScrollArea class="max-h-[85vh]">
				<CommunityEmotesDetailsContent />
			</ScrollArea>
		</DialogOrSheet>
	</Dialog>
</template>
