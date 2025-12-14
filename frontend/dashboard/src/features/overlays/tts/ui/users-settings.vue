<script setup lang="ts">
import { Trash2, Volume2 } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useTTSOverlayApi } from '@/api/overlays-tts'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
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
import { toast } from 'vue-sonner'

const { t } = useI18n()

const api = useTTSOverlayApi()
const { data: usersData, fetching: isLoading } = api.useQueryTTSUsersSettings()
const deleteMutation = api.useMutationDeleteTTSUsersSettings()

const testText = ref('Hello world, this is a test message')
const selectedUsers = ref<Set<string>>(new Set())
const showDeleteDialog = ref(false)

const users = computed(() => usersData.value?.overlaysTTSUsersSettings || [])

const selectedCount = computed(() => selectedUsers.value.size)

const allSelected = computed(() => {
	return users.value.length > 0 && selectedUsers.value.size === users.value.length
})

function toggleAll() {
	if (allSelected.value) {
		selectedUsers.value.clear()
	} else {
		selectedUsers.value = new Set(users.value.map((u) => u.userId))
	}
}

function toggleUser(userId: string) {
	if (selectedUsers.value.has(userId)) {
		selectedUsers.value.delete(userId)
	} else {
		selectedUsers.value.add(userId)
	}
}

async function handleDelete() {
	const userIds = Array.from(selectedUsers.value)

	const result = await deleteMutation.executeMutation({ userIds })

	if (!result.error) {
		toast.success(t('sharedTexts.deleted'), {
			description: t('overlays.tts.usersDeleted', { count: userIds.length }),
		})
		selectedUsers.value.clear()
		showDeleteDialog.value = false
	} else {
		toast.error(t('sharedTexts.error'))
	}
}

function testUserVoice(user: (typeof users.value)[0]) {
	// TODO: Implement actual preview using the TTS say mutation
	console.log('Testing voice for user:', user.userId, 'with text:', testText.value)
	toast({
		title: t('overlays.tts.playingPreview'),
	})
}
</script>

<template>
	<div class="flex flex-col gap-4 p-4">
		<Card>
			<CardHeader>
				<div class="flex items-center justify-between">
					<div>
						<CardTitle>{{ t('overlays.tts.tabs.users') }}</CardTitle>
						<CardDescription>{{ t('overlays.tts.usersDescription') }}</CardDescription>
					</div>
					<div class="flex gap-2">
						<Button variant="outline" :disabled="users.length === 0" @click="toggleAll">
							{{
								allSelected
									? t('overlays.tts.users.undoSelection')
									: t('overlays.tts.users.selectAll')
							}}
						</Button>
						<Button
							variant="destructive"
							:disabled="selectedCount === 0"
							@click="showDeleteDialog = true"
						>
							<Trash2 class="h-4 w-4 mr-2" />
							{{ t('sharedButtons.delete') }} ({{ selectedCount }})
						</Button>
					</div>
				</div>
			</CardHeader>

			<CardContent class="space-y-4">
				<!-- Test Input -->
				<div class="flex gap-2">
					<Input
						v-model="testText"
						:placeholder="t('overlays.tts.testTextPlaceholder')"
						class="flex-1"
					/>
				</div>

				<!-- Empty State -->
				<Alert v-if="!isLoading && users.length === 0">
					<AlertTitle>{{ t('overlays.tts.users.empty') }}</AlertTitle>
					<AlertDescription>
						{{ t('overlays.tts.users.emptyDescription') }}
					</AlertDescription>
				</Alert>

				<!-- Users Grid -->
				<div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<Card
						v-for="user in users"
						:key="user.userId"
						class="cursor-pointer transition-colors hover:bg-accent"
						:class="{ 'bg-accent': selectedUsers.has(user.userId) }"
						@click="toggleUser(user.userId)"
					>
						<CardContent class="p-4">
							<div class="flex items-center justify-between">
								<div class="flex items-center gap-3 flex-1">
									<Avatar>
										<AvatarImage
											:src="user.twitchProfile.profileImageUrl ?? ''"
											:alt="user.twitchProfile.displayName ?? ''"
										/>
										<AvatarFallback>
											{{ user.twitchProfile.displayName?.charAt(0) || '?' }}
										</AvatarFallback>
									</Avatar>

									<div class="flex-1 min-w-0">
										<div class="flex items-center gap-2">
											<p class="font-semibold truncate">
												{{ user.twitchProfile.displayName }}
											</p>
											<Badge v-if="user.isChannelOwner" variant="secondary"> Owner </Badge>
										</div>
										<p class="text-sm text-muted-foreground">
											{{ t('overlays.tts.voice') }}: {{ user.voice }} |
											{{ t('overlays.tts.pitch') }}: {{ user.pitch }} |
											{{ t('overlays.tts.rate') }}: {{ user.rate }}
										</p>
									</div>
								</div>

								<div class="flex items-center gap-2" @click.stop>
									<Button variant="ghost" size="icon" @click="testUserVoice(user)">
										<Volume2 class="h-4 w-4" />
									</Button>
									<Checkbox
										:checked="selectedUsers.has(user.userId)"
										@update:checked="() => toggleUser(user.userId)"
									/>
								</div>
							</div>
						</CardContent>
					</Card>
				</div>
			</CardContent>
		</Card>

		<!-- Delete Confirmation Dialog -->
		<AlertDialog v-model:open="showDeleteDialog">
			<AlertDialogContent>
				<AlertDialogHeader>
					<AlertDialogTitle>{{ t('deleteConfirmation.title') }}</AlertDialogTitle>
					<AlertDialogDescription>
						{{ t('overlays.tts.deleteUsersConfirmation', { count: selectedCount }) }}
					</AlertDialogDescription>
				</AlertDialogHeader>
				<AlertDialogFooter>
					<AlertDialogCancel>{{ t('deleteConfirmation.cancel') }}</AlertDialogCancel>
					<AlertDialogAction @click="handleDelete">
						{{ t('deleteConfirmation.confirm') }}
					</AlertDialogAction>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	</div>
</template>
