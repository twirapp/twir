<script setup lang="ts">
import {
	AccordionContent,
	AccordionHeader,
	AccordionItem,
	AccordionRoot,
	AccordionTrigger,
} from 'radix-vue'
import { ChevronDown } from 'lucide-vue-next'

import EventsTab from '@/features/overlays/kappagen/components/events-tab.vue'
import SettingsTab from '@/features/overlays/kappagen/components/settings-tab.vue'
import AnimationsTab from '@/features/overlays/kappagen/components/animations-tab.vue'

const content = [
	{
		name: 'settings',
		title: 'Settings',
		component: SettingsTab,
	},
	{
		name: 'animations',
		title: 'Animations',
		component: AnimationsTab,
	},
	{
		name: 'events',
		title: 'Events',
		component: EventsTab,
	},
]
</script>

<template>
	<div class="flex flex-col gap-4 mt-2 z-50">
		<AccordionRoot
			class="w-[400px] flex flex-col gap-4"
			default-value="animations"
			type="single"
			:collapsible="true"
		>
			<AccordionItem
				v-for="item in content"
				:key="item.name"
				:value="item.name"
				class="flex flex-col gap-2"
			>
				<AccordionHeader class="flex">
					<AccordionTrigger
						class="flex h-[45px] flex-1 cursor-pointer items-center justify-between bg-stone-700/40 px-5 text-[15px] leading-none outline-none group rounded-md hover"
					>
						<span>{{ item.title }}</span>
						<ChevronDown class="group-data-[state=open]:rotate-180" />
					</AccordionTrigger>
				</AccordionHeader>
				<AccordionContent
					class="data-[state=open]:animate-accordion-down data-[state=closed]:animate-accordion-up data-[state=open]:block data-[state=closed]:hidden overflow-hidden text-[15px]"
					force-mount
				>
					<component :is="item.component" />
				</AccordionContent>
			</AccordionItem>
		</AccordionRoot>
	</div>
</template>
