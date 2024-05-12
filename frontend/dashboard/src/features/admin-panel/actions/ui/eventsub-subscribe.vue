<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { useI18n } from 'vue-i18n'
import * as z from 'zod'

import { useMutationEventSubSubscribe } from '@/api/admin/actions'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter } from '@/components/ui/card'
import { FormControl, FormDescription, FormField, FormItem, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import { useToast } from '@/components/ui/toast'
import { EventsubSubscribeConditionInput } from '@/gql/graphql'

const { t } = useI18n()

const mutationEventSubSubscribe = useMutationEventSubSubscribe()

const conditionTranslations = {
	[EventsubSubscribeConditionInput.User]: {
		label: t('adminPanel.adminActions.eventsub.user'),
		examplePayload: '{ "user_id": "1234" }',
	},
	[EventsubSubscribeConditionInput.Channel]: {
		label: t('adminPanel.adminActions.eventsub.channel'),
		examplePayload: '{ "broadcaster_user_id": "1234" }',
	},
	[EventsubSubscribeConditionInput.ChannelWithBotId]: {
		label: t('adminPanel.adminActions.eventsub.channelWithBotId'),
		examplePayload: '{ "broadcaster_user_id": "1234", "user_id": "1234" }',
	},
	[EventsubSubscribeConditionInput.ChannelWithModeratorId]: {
		label: t('adminPanel.adminActions.eventsub.channelWithModeratorId'),
		examplePayload: '{ "broadcaster_user_id": "1234", "moderator_user_id": "1234" }',
	},
}

const formSchema = toTypedSchema(z.object({
	condition: z.enum([
		EventsubSubscribeConditionInput.User,
		EventsubSubscribeConditionInput.Channel,
		EventsubSubscribeConditionInput.ChannelWithBotId,
		EventsubSubscribeConditionInput.ChannelWithModeratorId,
	]),
	type: z.string(),
	version: z.string(),
}))

const { handleSubmit } = useForm({
	validationSchema: formSchema,
	initialValues: {
		condition: EventsubSubscribeConditionInput.User,
		version: '1',
	},
})

const toast = useToast()
const onSubmit = handleSubmit(async (values) => {
	const result = await mutationEventSubSubscribe.executeMutation({
		opts: values,
	})

	if (result.error) {
		toast.toast({
			duration: 2500,
			variant: 'destructive',
			title: result.error.message,
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
					<FormField v-slot="{ componentField }" name="condition">
						<FormItem>
							<Label for="condition">
								{{ t('adminPanel.adminActions.eventsub.condition') }}
							</Label>

							<Select v-bind="componentField">
								<FormControl>
									<SelectTrigger>
										<SelectValue />
									</SelectTrigger>
								</FormControl>

								<SelectContent>
									<SelectGroup>
										<SelectItem
											v-for="condition in Object.values(EventsubSubscribeConditionInput)"
											:key="condition"
											:value="condition"
										>
											<div class="flex gap-2 items-center">
												<span>
													{{ conditionTranslations[condition].label }}
												</span>
												<span class="text-xs text-zinc-500">
													{{ conditionTranslations[condition].examplePayload }}
												</span>
											</div>
										</SelectItem>
									</SelectGroup>
								</SelectContent>
							</Select>
							<FormMessage />
						</FormItem>
					</FormField>

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
									href="https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/" target="_blank"
								>subscription types</a>.
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
