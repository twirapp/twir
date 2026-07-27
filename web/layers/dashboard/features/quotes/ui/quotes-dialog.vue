<script setup lang="ts">
import type { Quote } from '~~/layers/dashboard/api/quotes.js'

import { useForm } from 'vee-validate'
import { ref } from 'vue'
import { toast } from 'vue-sonner'
import * as z from 'zod'

import { useQuotesApi } from '~~/layers/dashboard/api/quotes.js'
import DialogOrSheet from '~~/layers/dashboard/components/dialog-or-sheet.vue'

import { Button } from '@/components/ui/button'
import { Dialog, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Textarea } from '@/components/ui/textarea'

const props = defineProps<{
	quote?: Quote | null
}>()

const emit = defineEmits<{
	close: []
}>()

const { t } = useI18n()
const open = ref(false)
const quotesApi = useQuotesApi()
const createMutation = quotesApi.useMutationCreateQuote()
const updateMutation = quotesApi.useMutationUpdateQuote()
const quoteForm = useForm({
	validationSchema: z.object({
		text: z.string().min(1).max(5000),
	}),
	initialValues: {
		text: props.quote?.text ?? '',
	},
})

const save = quoteForm.handleSubmit(async (values) => {
	const result = props.quote
		? await updateMutation.executeMutation({ id: props.quote.id, opts: values })
		: await createMutation.executeMutation({ opts: values })

	if (result.error) {
		toast.error(t('quotes.errorOnSave'))
		return
	}

	toast.success(t(props.quote ? 'quotes.updated' : 'quotes.created'))
	emit('close')
	open.value = false
	quoteForm.resetForm()
})
</script>

<template>
	<Dialog v-model:open="open">
		<DialogTrigger as-child>
			<slot name="dialog-trigger" />
		</DialogTrigger>
		<DialogOrSheet class="sm:max-w-[424px]">
			<DialogHeader>
				<DialogTitle>
					{{ quote ? t('quotes.edit') : t('quotes.create') }}
				</DialogTitle>
			</DialogHeader>
			<form class="flex flex-col gap-4" @submit="save">
				<FormField v-slot="{ componentField }" name="text">
					<FormItem>
						<FormLabel>{{ t('quotes.text') }}</FormLabel>
						<FormControl>
							<Textarea v-bind="componentField" :maxlength="5000" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>
				<DialogFooter>
					<Button type="submit">
						{{ t('sharedButtons.save') }}
					</Button>
				</DialogFooter>
			</form>
		</DialogOrSheet>
	</Dialog>
</template>
