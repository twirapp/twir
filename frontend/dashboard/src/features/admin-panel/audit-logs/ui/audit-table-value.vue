<script setup lang="ts">
import type { AdminAuditLogsQuery } from '@/gql/graphql'

import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'

defineProps<{
	log: AdminAuditLogsQuery['adminAuditLogs']['logs'][0]
}>()

function computeDisplayedText(text?: string | null) {
	if (!text) return 'N/A'

	try {
		const json = JSON.parse(text)
		return JSON.stringify(json, null, 4)
	} catch (e) {
		return text
	}
}
</script>

<template>
	<Accordion type="multiple" collapsible>
		<AccordionItem value="oldValue" :disabled="!log.oldValue">
			<AccordionTrigger>Old value</AccordionTrigger>
			<AccordionContent>
				<pre class="relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-sm font-semibold">
					{{ computeDisplayedText(log.oldValue) }}
				</pre>
			</AccordionContent>
		</AccordionItem>
		<AccordionItem value="newValue" :disabled="!log.newValue">
			<AccordionTrigger>New value</AccordionTrigger>
			<AccordionContent>
				<pre class="relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-sm font-semibold">
					{{ computeDisplayedText(log.newValue) }}
				</pre>
			</AccordionContent>
		</AccordionItem>
	</Accordion>
</template>
