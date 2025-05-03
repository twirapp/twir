<script setup lang="ts">
import { Pencil, Trash2 } from 'lucide-vue-next'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import { useUserAccessFlagChecker } from '@/api'
import { type Event, useEventsApi } from '@/api/events'
import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { useToast } from '@/components/ui/toast/use-toast'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const props = defineProps<{
	event: Event
}>()

const { t } = useI18n()
const router = useRouter()
const { toast } = useToast()
const eventsApi = useEventsApi()
const deleteEventMutation = eventsApi.useMutationDeleteEvent()
const showDeleteDialog = ref(false)
const userCanManageEvents = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageEvents)

function editEvent() {
	if (!userCanManageEvents.value) return
	router.push(`/dashboard/events/${props.event.id}`)
}

async function deleteEvent() {
	if (!userCanManageEvents.value) return

	try {
		await deleteEventMutation.executeMutation({ id: props.event.id })
		toast({
			title: t('events.deleteSuccess'),
			description: t('events.deleteSuccessDescription'),
		})
		showDeleteDialog.value = false
	} catch (error) {
		console.error(error)
		toast({
			title: t('events.deleteError'),
			description: t('events.deleteErrorDescription'),
			variant: 'destructive',
		})
	}
}
</script>

<template>
	<div class="flex items-center gap-2">
		<Button variant="ghost" size="icon" :disabled="!userCanManageEvents" @click="editEvent">
			<Pencil class="h-4 w-4" />
		</Button>
		<Button variant="ghost" size="icon" :disabled="!userCanManageEvents" @click="showDeleteDialog = true">
			<Trash2 class="h-4 w-4" />
		</Button>

		<AlertDialog :open="showDeleteDialog" @update:open="showDeleteDialog = $event">
			<AlertDialogContent>
				<AlertDialogHeader>
					<AlertDialogTitle>{{ t('events.deleteConfirmTitle') }}</AlertDialogTitle>
					<AlertDialogDescription>
						{{ t('events.deleteConfirmDescription') }}
					</AlertDialogDescription>
				</AlertDialogHeader>
				<AlertDialogFooter>
					<AlertDialogCancel>{{ t('sharedTexts.cancel') }}</AlertDialogCancel>
					<AlertDialogAction @click="deleteEvent">
						{{ t('sharedTexts.delete') }}
					</AlertDialogAction>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	</div>
</template>
