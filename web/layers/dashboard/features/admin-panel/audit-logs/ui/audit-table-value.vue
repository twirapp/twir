<script setup lang="ts">
import formatJson from '@crashmax/json-format-highlight'

import type { AdminAuditLogsQuery } from '@/gql/graphql'

import {
	Accordion,
	AccordionContent,
	AccordionItem,
	AccordionTrigger,
} from '@/components/ui/accordion'
import { ScrollArea } from '@/components/ui/scroll-area'

defineProps<{
	log: AdminAuditLogsQuery['adminAuditLogs']['logs'][0]
}>()

function computeDisplayedText(text?: string | null) {
	if (!text) return 'N/A'

	try {
		const json = JSON.parse(text)
		return formatJson(json, {
			wordWrap: true,
		})
	} catch {
		return text
	}
}
</script>

<template>
	<Accordion type="multiple" collapsible>
		<AccordionItem v-if="log.oldValue" value="oldValue">
			<AccordionTrigger>Old value</AccordionTrigger>
			<AccordionContent>
				<ScrollArea class="max-h-[200px] rounded-md border bg-[#1e1e1e]">
					<pre class="code" v-html="computeDisplayedText(log.oldValue)" />
				</ScrollArea>
			</AccordionContent>
		</AccordionItem>
		<AccordionItem v-if="log.newValue" value="newValue" class="border-none">
			<AccordionTrigger>New value</AccordionTrigger>
			<AccordionContent>
				<ScrollArea class="max-h-[200px] rounded-md border bg-[#1e1e1e]">
					<pre class="code" v-html="computeDisplayedText(log.newValue)" />
				</ScrollArea>
			</AccordionContent>
		</AccordionItem>
	</Accordion>
</template>

<style scoped>
@reference '@/assets/index.css';

.code {
	@apply p-2 bg-[#1e1e1e] select-text max-h-[200px];
}
</style>
