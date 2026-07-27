<script setup lang="ts">
import { intlFormat } from 'date-fns'

import { useQuotes } from '#layers/public/api/use-quotes.ts'

definePageMeta({
	layout: 'public',
})

const { data } = await useQuotes()
</script>

<template>
	<div class="bg-card w-full flex-wrap rounded-md border">
		<UiTable>
			<UiTableHeader>
				<UiTableRow>
					<UiTableHead class="w-[8%]"> # </UiTableHead>
					<UiTableHead class="w-[52%]"> Quote </UiTableHead>
					<UiTableHead class="w-[15%]"> Added by </UiTableHead>
					<UiTableHead class="w-[12%]"> Game </UiTableHead>
					<UiTableHead class="w-[13%]"> Date </UiTableHead>
				</UiTableRow>
			</UiTableHeader>
			<UiTableBody>
				<UiTableRow
					v-for="quote in data?.quotesPublic"
					:key="quote.number"
				>
					<UiTableCell class="font-medium">
						#{{ quote.number }}
					</UiTableCell>
					<UiTableCell class="whitespace-normal wrap-break-word">
						{{ quote.text }}
					</UiTableCell>
					<UiTableCell>
						{{ quote.creatorName ?? '-' }}
					</UiTableCell>
					<UiTableCell>
						{{ quote.gameName ?? '-' }}
					</UiTableCell>
					<UiTableCell>
						{{
							intlFormat(new Date(quote.createdAt), {
								day: 'numeric',
								month: 'numeric',
								year: 'numeric',
							})
						}}
					</UiTableCell>
				</UiTableRow>
			</UiTableBody>
		</UiTable>
	</div>
</template>
