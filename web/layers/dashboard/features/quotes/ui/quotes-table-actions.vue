<script setup lang="ts">
import type { Quote } from '~~/layers/dashboard/api/quotes.js'

import { ref } from 'vue'
import { toast } from 'vue-sonner'

import { useQuotesApi } from '~~/layers/dashboard/api/quotes.js'

import QuotesDialog from './quotes-dialog.vue'

import ActionConfirm from '@/components/ui/action-confirm/ActionConfirm.vue'
import { Button } from '@/components/ui/button'

const props = defineProps<{
	quote: Quote
}>()

const { t } = useI18n()
const showDelete = ref(false)
const quotesApi = useQuotesApi()
const removeMutation = quotesApi.useMutationRemoveQuote()

async function deleteQuote() {
	const result = await removeMutation.executeMutation({ id: props.quote.id })
	if (result.error) {
		toast.error(t('quotes.errorOnDelete'))
		return
	}

	toast.success(t('quotes.deleted'))
}
</script>

<template>
	<div class="flex items-center gap-2">
		<QuotesDialog :quote="quote">
			<template #dialog-trigger>
				<Button variant="secondary" size="icon">
					<Icon name="lucide:pencil" class="size-4" />
				</Button>
			</template>
		</QuotesDialog>
		<Button variant="destructive" size="icon" @click="showDelete = true">
			<Icon name="lucide:trash" class="size-4" />
		</Button>
	</div>

	<ActionConfirm
		v-model:open="showDelete"
		:confirm-text="t('quotes.deleteConfirmation')"
		@confirm="deleteQuote"
	/>
</template>
