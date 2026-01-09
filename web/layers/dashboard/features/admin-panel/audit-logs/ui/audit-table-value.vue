<script setup lang="ts">
import formatJson from '@crashmax/json-format-highlight'

import type { AdminAuditLogsQuery } from '~/gql/graphql'




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
	<UiAccordion type="multiple" collapsible>
		<UiAccordionItem v-if="log.oldValue" value="oldValue">
			<UiAccordionTrigger>Old value</UiAccordionTrigger>
			<UiAccordionContent>
				<UiScrollArea class="max-h-[200px] rounded-md border bg-[#1e1e1e]">
					<pre class="code" v-html="computeDisplayedText(log.oldValue)" />
				</UiScrollArea>
			</UiAccordionContent>
		</UiAccordionItem>
		<UiAccordionItem v-if="log.newValue" value="newValue" class="border-none">
			<UiAccordionTrigger>New value</UiAccordionTrigger>
			<UiAccordionContent>
				<UiScrollArea class="max-h-[200px] rounded-md border bg-[#1e1e1e]">
					<pre class="code" v-html="computeDisplayedText(log.newValue)" />
				</UiScrollArea>
			</UiAccordionContent>
		</UiAccordionItem>
	</UiAccordion>
</template>

<style scoped>
@reference '~/assets/css/tailwind.css';

.code {
	@apply p-2 bg-[#1e1e1e] select-text max-h-[200px];
}
</style>
