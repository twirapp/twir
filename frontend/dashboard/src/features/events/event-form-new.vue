<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import * as z from 'zod'

import EventBasicInfo from './components/event-basic-info.vue'
import FormActions from './components/form-actions.vue'
import OperationsTab from './components/operations-tab.vue'

import { useEventsApi } from '@/api/events'
import { Button } from '@/components/ui/button'
import { Form } from '@/components/ui/form'
import { useToast } from '@/components/ui/toast/use-toast'
import PageLayout from '@/layout/page-layout.vue'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const { toast } = useToast()
const eventsApi = useEventsApi()
const isNewEvent = computed(() => route.params.id === 'new')
const eventId = computed(() => isNewEvent.value ? '' : String(route.params.id))

// Fetch event data if editing
const { data: eventData, fetching: isLoadingEvent } = eventsApi.useQueryEventById(eventId.value)
const event = computed(() => eventData.value?.eventById)

// Form validation schema
const formSchema = toTypedSchema(z.object({
	type: z.string().min(1, t('events.validation.typeRequired')),
	description: z.string().min(1, t('events.validation.descriptionRequired')),
	enabled: z.boolean().default(true),
	onlineOnly: z.boolean().default(false),
	rewardId: z.string().optional(),
	commandId: z.string().optional(),
	keywordId: z.string().optional(),
	operations: z.array(z.object({
		type: z.string().min(1, t('events.validation.operationTypeRequired')),
		input: z.string().optional(),
		delay: z.number().min(0).default(0),
		repeat: z.number().min(0).default(0),
		useAnnounce: z.boolean().default(false),
		timeoutTime: z.number().min(0).default(0),
		timeoutMessage: z.string().optional(),
		target: z.string().optional(),
		enabled: z.boolean().default(true),
		filters: z.array(z.object({
			type: z.string().min(1, t('events.validation.filterTypeRequired')),
			left: z.string().min(1, t('events.validation.leftRequired')),
			right: z.string().min(1, t('events.validation.rightRequired')),
		})).default([]),
	})).default([]),
}))

// Initialize form
const form = useForm({
	validationSchema: formSchema,
	initialValues: {
		type: '',
		description: '',
		enabled: true,
		onlineOnly: false,
		operations: [],
	},
})

// Update form values when event data is loaded
watch(() => event.value, (newEvent) => {
	if (newEvent) {
		form.setValues({
			type: newEvent.type,
			description: newEvent.description,
			enabled: newEvent.enabled,
			onlineOnly: newEvent.onlineOnly,
			rewardId: newEvent.rewardId || undefined,
			commandId: newEvent.commandId || undefined,
			keywordId: newEvent.keywordId || undefined,
			operations: newEvent.operations.map(op => ({
				type: op.type,
				input: op.input || undefined,
				delay: op.delay,
				repeat: op.repeat,
				useAnnounce: op.useAnnounce,
				timeoutTime: op.timeoutTime,
				timeoutMessage: op.timeoutMessage || undefined,
				target: op.target || undefined,
				enabled: op.enabled,
				filters: op.filters.map(filter => ({
					type: filter.type,
					left: filter.left,
					right: filter.right,
				})),
			})),
		})
	}
}, { immediate: true })

// Form submission
const createEventMutation = eventsApi.useMutationCreateEvent()
const updateEventMutation = eventsApi.useMutationUpdateEvent()
const isSubmitting = ref(false)

async function onSubmit(values: z.infer<typeof formSchema>) {
	isSubmitting.value = true

	try {
		if (isNewEvent.value) {
			await createEventMutation.executeMutation({
				input: {
					type: values.type,
					description: values.description,
					enabled: values.enabled,
					onlineOnly: values.onlineOnly,
					rewardId: values.rewardId,
					commandId: values.commandId,
					keywordId: values.keywordId,
					operations: values.operations.map(op => ({
						type: op.type,
						input: op.input,
						delay: op.delay,
						repeat: op.repeat,
						useAnnounce: op.useAnnounce,
						timeoutTime: op.timeoutTime,
						timeoutMessage: op.timeoutMessage,
						target: op.target,
						enabled: op.enabled,
						filters: op.filters,
					})),
				},
			})

			toast({
				title: t('events.createSuccess'),
				description: t('events.createSuccessDescription'),
			})
		} else {
			await updateEventMutation.executeMutation({
				id: eventId.value,
				input: {
					type: values.type,
					description: values.description,
					enabled: values.enabled,
					onlineOnly: values.onlineOnly,
					rewardId: values.rewardId,
					commandId: values.commandId,
					keywordId: values.keywordId,
					operations: values.operations.map(op => ({
						type: op.type,
						input: op.input,
						delay: op.delay,
						repeat: op.repeat,
						useAnnounce: op.useAnnounce,
						timeoutTime: op.timeoutTime,
						timeoutMessage: op.timeoutMessage,
						target: op.target,
						enabled: op.enabled,
						filters: op.filters,
					})),
				},
			})

			toast({
				title: t('events.updateSuccess'),
				description: t('events.updateSuccessDescription'),
			})
		}

		router.push('/dashboard/events')
	} catch (error) {
		console.error(error)
		toast({
			title: isNewEvent.value ? t('events.createError') : t('events.updateError'),
			description: isNewEvent.value ? t('events.createErrorDescription') : t('events.updateErrorDescription'),
			variant: 'destructive',
		})
	} finally {
		isSubmitting.value = false
	}
}

function goBack() {
	router.push('/dashboard/events')
}
</script>

<template>
	<PageLayout>
		<template #title>
			{{ isNewEvent ? t('events.create') : t('events.edit') }}
		</template>

		<template #action>
			<Button variant="outline" @click="goBack">
				{{ t('sharedTexts.back') }}
			</Button>
		</template>

		<template #content>
			<div v-if="isLoadingEvent && !isNewEvent" class="flex justify-center items-center h-64">
				<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
			</div>

			<Form v-else :form="form" class="space-y-6" @submit="onSubmit">
				<EventBasicInfo :form="form" />
				<OperationsTab :form="form" />
				<FormActions
					:is-new-event="isNewEvent"
					:is-submitting="isSubmitting"
					:on-cancel="goBack"
				/>
			</Form>
		</template>
	</PageLayout>
</template>
