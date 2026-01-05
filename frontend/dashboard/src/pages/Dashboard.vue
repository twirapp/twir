<script setup lang="ts">
import { GridItem, GridLayout } from 'grid-layout-plus'
import { SquarePen } from 'lucide-vue-next'
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'

import AuditLogs from '@/components/dashboard/audit-logs.vue'
import Chat from '@/components/dashboard/chat.vue'
import Events from '@/components/dashboard/events.vue'
import Stream from '@/components/dashboard/stream.vue'
import { useWidgets } from '@/components/dashboard/widgets.js'
import { Button } from '@/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useIsMobile } from '@/composables/use-is-mobile'

const { isMobile } = useIsMobile()
const widgets = useWidgets()
const visibleWidgets = computed(() => widgets.value.filter((v) => v.visible))
const invisibleWidgets = computed(() => widgets.value.filter((v) => !v.visible))

function addWidget(key: string | number) {
	const item = widgets.value.find((v) => v.i === key)
	if (!item) return

	const widgetsLength = visibleWidgets.value.length

	item.visible = true
	item.x = (widgetsLength * 2) % 12
	item.y = widgetsLength + 12
}

const showEmptyItem = ref(false)

function onMouseUp() {
	showEmptyItem.value = false
}

onMounted(async () => {
	await nextTick()

	document.querySelectorAll('.vgl-item__resizer').forEach((el) => {
		el.addEventListener('mousedown', () => {
			showEmptyItem.value = true
		})

		window.addEventListener('mouseup', onMouseUp)
	})
})

onBeforeUnmount(() => {
	window.removeEventListener('mouseup', onMouseUp)
})
</script>

<template>
	<div class="w-full h-full pl-1">
		<GridLayout v-model:layout="widgets" :row-height="30" :use-css-transforms="false">
			<GridItem
				v-for="item in visibleWidgets"
				:key="item.i"
				:x="item.x"
				:y="item.y"
				:w="item.w"
				:h="item.h"
				:i="item.i"
				:min-w="item.minW"
				:min-h="item.minH"
				drag-allow-from=".widgets-draggable-handle"
			>
				<div v-if="showEmptyItem" class="w-full h-full absolute z-50"></div>
				<Chat v-if="item.i === 'chat'" :item="item" class="h-full" />
				<Stream v-if="item.i === 'stream'" :item="item" class="h-full" />
				<Events v-if="item.i === 'events'" :item="item" class="h-full" />
				<AuditLogs v-if="item.i === 'audit-logs'" :item="item" class="h-full" />
			</GridItem>
		</GridLayout>

		<div
			v-if="invisibleWidgets.length"
			class="fixed right-8 bottom-8 z-50"
			:class="[{ 'right-24!': isMobile }]"
		>
			<DropdownMenu>
				<DropdownMenuTrigger as-child>
					<Button variant="secondary" class="h-14 w-14" size="icon">
						<SquarePen class="size-8" />
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent align="end">
					<DropdownMenuItem
						v-for="widget in invisibleWidgets"
						:key="widget.i"
						@click="addWidget(widget.i)"
					>
						{{ String(widget.i) }}
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
		</div>
	</div>
</template>

<style scoped>
@reference '@/assets/index.css';

.vgl-layout {
	@apply w-full;
}

:deep(.vgl-item__resizer) {
	z-index: 51;
}
</style>
