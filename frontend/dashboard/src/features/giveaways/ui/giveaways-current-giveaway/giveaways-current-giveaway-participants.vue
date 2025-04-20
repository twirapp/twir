<script setup lang="ts">
import { useVirtualizer } from '@tanstack/vue-virtual'
import { debouncedRef } from '@vueuse/core'
import { computed, ref, useTemplateRef } from 'vue'

import { useGiveaways } from '@/features/giveaways/composables/giveaways-use-api.ts'

const { participants } = useGiveaways()

const searchTerm = ref('')
const debouncerSearchTerm = debouncedRef(searchTerm, 200)
const filteredParticipants = computed(() => {
	return participants.value.filter(p => p.toString().includes(debouncerSearchTerm.value))
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
	<div class="flex-1">
		<!--		<div class="py-1 pr-2"> -->
		<!--			<span>Total participants: {{ participants.length }}</span> -->
		<!--			<Input v-model:model-value="searchTerm" placeholder="Search..." class="h-8" /> -->
		<!--		</div> -->
		<div ref="participantsRef" class="overflow-auto h-[90%]">
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
					class="border-b border-border pl-2 flex items-center"
					:style="{
						position: 'absolute',
						top: 0,
						left: 0,
						width: '100%',
						height: `${virtualRow.size}px`,
						transform: `translateY(${virtualRow.start}px)`,
					}"
				>
					{{ filteredParticipants[virtualRow.index] }}
				</div>
			</div>
		</div>
	</div>
</template>
