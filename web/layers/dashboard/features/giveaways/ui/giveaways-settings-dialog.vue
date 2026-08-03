<script setup lang="ts">
import { useForm } from 'vee-validate'
import { computed, ref, watch } from 'vue'
import { toast } from 'vue-sonner'
import { z } from 'zod'
import { useGiveawaysApi } from '~~/layers/dashboard/api/giveaways.js'

import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogClose,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Textarea } from '@/components/ui/textarea'

const { t } = useI18n()
const api = useGiveawaysApi()

const open = ref(false)

const { data: settingsData, fetching } = api.useGiveawaysSettings()
const settings = computed(() => {
	const data: any = settingsData.value?.giveawaysSettings
	if (!data) return null
	return {
		id: data.id,
		channelId: data.channelId,
		winnerMessage: data.winnerMessage,
	}
})

const formSchema = z.object({
	winnerMessage: z
		.string()
		.min(1, t('giveaways.settings.validation.required'))
		.max(500, t('giveaways.settings.validation.maxLength')),
})

const form = useForm({
	validationSchema: formSchema,
	initialValues: {
		winnerMessage: 'Congratulations {winner}! You won the giveaway!',
	},
	keepValuesOnUnmount: true,
})

watch(
	settings,
	(newSettings) => {
		if (newSettings) {
			form.setValues({
				winnerMessage: newSettings.winnerMessage,
			})
		}
	},
	{ immediate: true }
)

const updateMutation = api.useMutationUpdateSettings()

const buttonText = computed(() => {
	if (
		!settings.value?.winnerMessage ||
		settings.value.winnerMessage === 'Congratulations {winner}! You won the giveaway!'
	) {
		return t('giveaways.settings.create')
	}
	return t('giveaways.settings.update')
})

const handleSubmit = form.handleSubmit(async (values) => {
	const result = await updateMutation.executeMutation({
		opts: {
			winnerMessage: values.winnerMessage,
		},
	})

	if (result.error) {
		toast.error(t('giveaways.settings.error'), {
			description: result.error.message,
		})
		return
	}

	toast.success(t('giveaways.settings.success'), {
		description: t('giveaways.settings.successDescription'),
	})

	open.value = false
})
</script>

<template>
	<Dialog v-model:open="open">
		<DialogTrigger as-child>
			<Button>
				{{ t('giveaways.settings.title') }}
			</Button>
		</DialogTrigger>

		<DialogContent class="sm:max-w-[600px]">
			<DialogHeader>
				<DialogTitle>{{ t('giveaways.settings.dialogTitle') }}</DialogTitle>
				<DialogDescription>
					{{ t('giveaways.settings.dialogDescription') }}
				</DialogDescription>
			</DialogHeader>

			<Alert v-if="!fetching">
				<Icon
					name="lucide:info"
					class="h-4 w-4"
				/>
				<AlertDescription>
					{{ t('giveaways.settings.hint') }}
					<code class="bg-muted rounded px-1 py-0.5 text-sm">{winner}</code>
				</AlertDescription>
			</Alert>

			<form
				@submit="handleSubmit"
				class="space-y-4"
			>
				<FormField
					v-slot="{ componentField }"
					name="winnerMessage"
				>
					<FormItem>
						<FormLabel>{{ t('giveaways.settings.messageLabel') }}</FormLabel>
						<FormControl>
							<Textarea
								v-bind="componentField"
								:placeholder="t('giveaways.settings.messagePlaceholder')"
								class="min-h-[100px]"
								:disabled="fetching"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<DialogFooter>
					<DialogClose as-child>
						<Button
							type="button"
							variant="outline"
						>
							{{ t('sharedButtons.cancel') }}
						</Button>
					</DialogClose>
					<Button
						type="submit"
						:disabled="updateMutation.fetching.value || fetching"
					>
						{{ buttonText }}
					</Button>
				</DialogFooter>
			</form>
		</DialogContent>
	</Dialog>
</template>
