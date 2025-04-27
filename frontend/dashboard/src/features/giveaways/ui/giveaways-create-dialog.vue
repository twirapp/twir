<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { PlusIcon } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { z } from 'zod'

import { Button } from '@/components/ui/button'
import { Dialog, DialogContent ,	DialogDescription,	DialogFooter,	DialogHeader,	DialogTitle,	DialogTrigger } from '@/components/ui/dialog'
import {
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { useGiveaways } from '@/features/giveaways/composables/giveaways-use-giveaways.ts'

const { t } = useI18n()
const open = ref(false)
const { createGiveaway } = useGiveaways()

// Form validation schema
const formSchema = toTypedSchema(z.object({
	keyword: z.string()
		.min(3, 'Keyword must be at least 3 characters')
		.max(100, 'Keyword must be at most 100 characters'),
}))

// Form setup
const giveawayCreateForm = useForm({
	validationSchema: formSchema,
	initialValues: {
		keyword: '',
	},
	validateOnMount: false,
})

const handleSubmit = giveawayCreateForm.handleSubmit(async (values) => {
	try {
		const result = await createGiveaway(values.keyword)
		if (result) {
			giveawayCreateForm.resetForm()
		}
	} catch (error) {
		console.error(error)
	}
})
</script>

<template>
	<Dialog v-model:open="open">
		<DialogTrigger as-child>
			<Button size="sm" class="flex gap-2 items-center">
				<PlusIcon class="size-4" />
				{{ t('giveaways.createNew') }}
			</Button>
		</DialogTrigger>

		<DialogContent class="sm:max-w-[425px]">
			<DialogHeader>
				<DialogTitle>Create New Giveaway</DialogTitle>
				<DialogDescription>
					Enter a keyword for your giveaway. Users will use this keyword to participate.
				</DialogDescription>
			</DialogHeader>

			<form class="space-y-4" @submit.prevent="handleSubmit">
				<FormField
					v-slot="{ componentField, errorMessage }"
					name="keyword"
				>
					<FormItem>
						<FormLabel>Keyword</FormLabel>
						<FormControl>
							<Input
								placeholder="Enter keyword (e.g. '!giveaway' or 'raffle')"
								v-bind="componentField"
							/>
						</FormControl>
						<FormMessage>{{ errorMessage }}</FormMessage>
					</FormItem>
				</FormField>

				<DialogFooter>
					<Button
						type="button"
						variant="outline"
						@click="open = false"
					>
						Cancel
					</Button>
					<Button
						type="submit"
					>
						Create
					</Button>
				</DialogFooter>
			</form>
		</dialogcontent>
	</Dialog>
</template>
