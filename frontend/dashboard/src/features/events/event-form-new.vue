<script setup lang="ts">
import { useForm } from 'vee-validate'
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import EventBasicInfo from './components/event-basic-info.vue'
import OperationsTab from './components/operations-tab.vue'

import { EventType, useEventsApi } from '@/api/events'
import { Button } from '@/components/ui/button'
import { useToast } from '@/components/ui/toast/use-toast'
import { eventFormSchema } from '@/features/events/event-form-schema.ts'
import PageLayout from '@/layout/page-layout.vue'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const { toast } = useToast()
const eventsApi = useEventsApi()
const isNewEvent = computed(() => route.params.id === 'new')
const eventId = computed(() => isNewEvent.value ? '' : String(route.params.id))

// Fetch event data if editing
const { fetching: isLoadingEvent, executeQuery } = eventsApi.useQueryEventById(eventId.value)

// Initialize form
const eventForm = useForm({
	validationSchema: eventFormSchema,
	initialValues: {
		type: EventType.Follow,
		description: '',
		enabled: true,
		onlineOnly: false,
		operations: [],
	},
	keepValuesOnUnmount: true,
})

onMounted(async () => {
	if (isNewEvent.value) return

	const { data } = await executeQuery()
	if (!data.value?.eventById) {
		toast({
			title: t('events.notFound'),
			description: t('events.notFoundDescription'),
			variant: 'destructive',
		})
		router.push('/dashboard/events')
		return
	}

	const event = data.value.eventById

	eventForm.setValues({
		type: event.type,
		description: event.description,
		enabled: event.enabled,
		onlineOnly: event.onlineOnly,
		rewardId: event.rewardId || undefined,
		commandId: event.commandId || undefined,
		keywordId: event.keywordId || undefined,
		operations: event.operations.map(op => ({
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
})

// Form submission
const createEventMutation = eventsApi.useMutationCreateEvent()
const updateEventMutation = eventsApi.useMutationUpdateEvent()

const onSubmit = eventForm.handleSubmit(async (input) => {
	try {
		if (isNewEvent.value) {
			await createEventMutation.executeMutation({
				input,
			})
		} else {
			await updateEventMutation.executeMutation({
				id: eventId.value,
				input,
			})
		}

		toast({
			title: t('events.createSuccess'),
			description: isNewEvent.value ? t('events.createSuccessDescription') : t('events.updateSuccessDescription'),
			duration: 2500,
		})
		router.push('/dashboard/events')
	} catch (error) {
		console.error(error)
		toast({
			title: isNewEvent.value ? t('events.createError') : t('events.updateError'),
			description: isNewEvent.value ? t('events.createErrorDescription') : t('events.updateErrorDescription'),
			variant: 'destructive',
		})
	}
})
</script>

<template>
	<PageLayout sticky-header show-back back-to="/dashboard/events">
		<template #title>
			{{ isNewEvent ? t('sharedTexts.create') : t('sharedTexts.edit') }}
		</template>

		<template #action>
			<Button type="submit" :disabled="eventForm.values.operations?.length === 0">
				{{ t('sharedButtons.save') }}
			</Button>
		</template>

		<template #content>
			<div v-if="isLoadingEvent && !isNewEvent" class="flex justify-center items-center h-64">
				<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
			</div>

			<form v-else class="space-y-6" @submit="onSubmit">
				<EventBasicInfo />
				<OperationsTab />
			</form>
		</template>
	</PageLayout>
</template>
