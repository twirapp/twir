<script setup lang="ts">
import { CheckIcon, ListFilterPlusIcon } from 'lucide-vue-next'
import { ref } from 'vue'

import type { EventType } from '~/gql/graphql.ts'




import { useEventsTable } from '~/features/events/composables/use-events-table.ts'
import { EventsOptions } from '~/features/events/constants/events.ts'

const table = useEventsTable()

const typeSelectOptions = Object.values(EventsOptions).map<{
	isGroup: boolean
	name: string
	value?: EventType
	childrens: Array<{ name: string; value: EventType }>
}>((value) => {
	if (value.childrens) {
		return {
			isGroup: true,
			name: value.name,
			childrens: Object.values(value.childrens!).map((c) => ({
				name: c.name,
				value: c.enumValue!,
			})),
		}
	}

	return {
		isGroup: false,
		name: value.name,
		value: value.enumValue,
		childrens: [],
	}
})

function handleSelect(type: EventType) {
	if (!table.selectedTypes.value.includes(type)) {
		table.selectedTypes.value.push(type)
	} else {
		table.selectedTypes.value = table.selectedTypes.value.filter((t) => t !== type)
	}
}

const open = ref(false)
</script>

<template>
	<UiPopover v-model:open="open">
		<UiPopoverTrigger as-child>
			<UiButton
				variant="outline"
				size="sm"
				class="min-w-[200px] flex gap-2 items-center justify-start"
			>
				<ListFilterPlusIcon class="size-4" />
				<span v-if="table.selectedTypes.value.length">
					{{ table.selectedTypes.value.length }} type(s) selected
				</span>
				<span v-else> Filter by type </span>
			</UiButton>
		</UiPopoverTrigger>
		<UiPopoverContent class="p-0" side="right" align="start">
			<UiCommand>
				<UiCommandInput placeholder="Search type..." />
				<UiCommandList>
					<template v-for="selectOption in typeSelectOptions">
						<UiCommandGroup
							v-if="selectOption.isGroup"
							:key="`${selectOption.name}.group`"
							:heading="selectOption.name"
						>
							<UiCommandItem
								v-for="child in selectOption.childrens"
								:key="child.value"
								:value="child.value"
								@select="handleSelect(child.value)"
							>
								{{ child.name }}
								<CheckIcon
									v-if="table.selectedTypes.value.includes(child.value)"
									class="ml-auto text-xs text-muted-foreground size-4"
								/>
							</UiCommandItem>
						</UiCommandGroup>
						<UiCommandGroup v-else>
							<UiCommandItem
								:key="selectOption.name"
								:value="selectOption.value!"
								@select="handleSelect(selectOption.value!)"
							>
								{{ selectOption.name }}
								<CheckIcon
									v-if="table.selectedTypes.value.includes(selectOption.value!)"
									class="ml-auto text-xs text-muted-foreground size-4"
								/>
							</UiCommandItem>
						</UiCommandGroup>
					</template>
				</UiCommandList>
			</UiCommand>
		</UiPopoverContent>
	</UiPopover>
</template>
