<script setup lang="ts">
import { useVirtualizer } from '@tanstack/vue-virtual'
import { debouncedRef } from '@vueuse/core'
import { computed, ref, useTemplateRef } from 'vue'

import { useGiveaways } from '@/features/giveaways/composables/giveaways-use-giveaways.ts'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'

const { participants } = useGiveaways()

const searchTerm = ref('')
const debouncerSearchTerm = debouncedRef(searchTerm, 200)
const filteredParticipants = computed(() => {
	return participants.value.filter(p =>
		p.displayName.toLowerCase().includes(debouncerSearchTerm.value.toLowerCase())
	)
})

const participantsRef = useTemplateRef('participantsRef')
const totalParticipants = computed(() => filteredParticipants.value.length)

const rowVirtualizer = useVirtualizer({
	get count() {
		return totalParticipants.value
	},
	getScrollElement: () => participantsRef.value,
	estimateSize: () => 42,
	overscan: 5,
})
const virtualRows = computed(() => rowVirtualizer.value.getVirtualItems())
const totalSize = computed(() => rowVirtualizer.value.getTotalSize())
</script>

<template>
	<div class="flex-1 flex flex-col">
		<div class="py-2 px-2 border-b border-border">
			<div class="flex items-center justify-between mb-2">
				<span class="text-sm font-medium">Total participants: {{ participants.length }}</span>
			</div>
			<Input v-model="searchTerm" placeholder="Search participants..." class="h-8" />
		</div>
		<div ref="participantsRef" class="overflow-auto flex-1">
			<div
				:style="{
					height: `${totalSize}px`,
					width: '100%',
					position: 'relative',
				}"
				class="flex-1 overflow-auto"
			>
				<div
					v-for="virtualRow in virtualRows"
					:key="virtualRow.index"
					class="border-b border-border pl-2 flex items-center justify-between pr-2"
					:style="{
						position: 'absolute',
						top: 0,
						left: 0,
						width: '100%',
						height: `${virtualRow.size}px`,
						transform: `translateY(${virtualRow.start}px)`,
					}"
				>
					<span>{{ filteredParticipants[virtualRow.index].displayName }}</span>
					<Badge v-if="filteredParticipants[virtualRow.index].isWinner" variant="success">Winner</Badge>
				</div>
			</div>
		</div>
	</div>
</template>
