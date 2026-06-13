<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { useI18n } from 'vue-i18n'
import * as z from 'zod'

import { useMutationEventSubSubscribe } from '@/api/admin/actions'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter } from '@/components/ui/card'
import {
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { toast } from 'vue-sonner'

const { t } = useI18n()

const mutationEventSubSubscribe = useMutationEventSubSubscribe()

const formSchema = toTypedSchema(
	z.object({
		type: z.string(),
		version: z.string(),
	})
)

const { handleSubmit } = useForm({
	validationSchema: formSchema,
	initialValues: {
		version: '1',
	},
})

const onSubmit = handleSubmit(async (values) => {
	const result = await mutationEventSubSubscribe.executeMutation({
		opts: values,
	})

	if (result.error) {
		toast.error(result.error.message, {
			duration: 2500,
		})
	}
})
</script>

<template>
	<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
		{{ t('adminPanel.adminActions.eventsub.title') }}
	</h4>

	<Card>
		<form @submit.prevent="onSubmit">
			<CardContent class="p-4">
				<div class="grid items-center w-full gap-4">
					<FormField v-slot="{ componentField }" name="version">
						<FormItem>
							<Label for="version">
								{{ t('adminPanel.adminActions.eventsub.version') }}
							</Label>
							<FormControl>
								<Input v-bind="componentField" />
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField v-slot="{ componentField }" name="type">
						<FormItem>
							<Label for="type">
								{{ t('adminPanel.adminActions.eventsub.type') }}
							</Label>
							<FormControl>
								<Input v-bind="componentField" />
							</FormControl>
							<FormDescription>
								You can find all available subscription types in
								<a
									class="underline"
									href="https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/"
									target="_blank"
									>subscription types</a
								>.
							</FormDescription>
							<FormMessage />
						</FormItem>
					</FormField>
				</div>
			</CardContent>

			<CardFooter class="flex justify-end p-4">
				<Button>
					{{ t('sharedButtons.send') }}
				</Button>
			</CardFooter>
		</form>
	</Card>
</template>
