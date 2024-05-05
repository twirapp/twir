<script setup lang="ts">
import { InfoIcon } from 'lucide-vue-next'

import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogTrigger,
} from '@/components/ui/dialog'
import CommunityEmotesDetailsContent from '@/features/community-emotes-statistic/components/community-emotes-details-content.vue'

defineEmits<{
	select: []
}>()
</script>

<template>
	<Dialog>
		<DialogTrigger as-child>
			<Button variant="secondary" size="icon" @click="$emit('select')">
				<InfoIcon class="size-4" />
			</Button>
		</DialogTrigger>
		<DialogOrSheet class="p-0">
			<CommunityEmotesDetailsContent />
		</DialogOrSheet>
	</Dialog>
</template>

<style scoped>

</style>
