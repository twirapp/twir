<script setup lang="ts">
import { Trash2, Volume2 } from 'lucide-vue-next'
import { computed, ref } from 'vue'


import { useTTSOverlayApi } from '#layers/dashboard/api/overlays-tts'








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
		<UiCard>
			<UiCardHeader>
				<div class="flex items-center justify-between">
					<div>
						<UiCardTitle>{{ t('overlays.tts.tabs.users') }}</UiCardTitle>
						<UiCardDescription>{{ t('overlays.tts.usersDescription') }}</UiCardDescription>
					</div>
					<div class="flex gap-2">
						<UiButton variant="outline" :disabled="users.length === 0" @click="toggleAll">
							{{
								allSelected
									? t('overlays.tts.users.undoSelection')
									: t('overlays.tts.users.selectAll')
							}}
						</UiButton>
						<UiButton
							variant="destructive"
							:disabled="selectedCount === 0"
							@click="showDeleteDialog = true"
						>
							<Trash2 class="h-4 w-4 mr-2" />
							{{ t('sharedButtons.delete') }} ({{ selectedCount }})
						</UiButton>
					</div>
				</div>
			</UiCardHeader>

			<UiCardContent class="space-y-4">
				<!-- Test Input -->
				<div class="flex gap-2">
					<UiInput
						v-model="testText"
						:placeholder="t('overlays.tts.testTextPlaceholder')"
						class="flex-1"
					/>
				</div>

				<!-- Empty State -->
				<UiAlert v-if="!isLoading && users.length === 0">
					<UiAlertTitle>{{ t('overlays.tts.users.empty') }}</UiAlertTitle>
					<UiAlertDescription>
						{{ t('overlays.tts.users.emptyDescription') }}
					</UiAlertDescription>
				</UiAlert>

				<!-- Users Grid -->
				<div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<UiCard
						v-for="user in users"
						:key="user.userId"
						class="cursor-pointer transition-colors hover:bg-accent"
						:class="{ 'bg-accent': selectedUsers.has(user.userId) }"
						@click="toggleUser(user.userId)"
					>
						<UiCardContent class="p-4">
							<div class="flex items-center justify-between">
								<div class="flex items-center gap-3 flex-1">
									<UiAvatar>
										<UiAvatarImage
											:src="user.twitchProfile.profileImageUrl ?? ''"
											:alt="user.twitchProfile.displayName ?? ''"
										/>
										<UiAvatarFallback>
											{{ user.twitchProfile.displayName?.charAt(0) || '?' }}
										</UiAvatarFallback>
									</UiAvatar>

									<div class="flex-1 min-w-0">
										<div class="flex items-center gap-2">
											<p class="font-semibold truncate">
												{{ user.twitchProfile.displayName }}
											</p>
											<UiBadge v-if="user.isChannelOwner" variant="secondary"> Owner </UiBadge>
										</div>
										<p class="text-sm text-muted-foreground">
											{{ t('overlays.tts.voice') }}: {{ user.voice }} |
											{{ t('overlays.tts.pitch') }}: {{ user.pitch }} |
											{{ t('overlays.tts.rate') }}: {{ user.rate }}
										</p>
									</div>
								</div>

								<div class="flex items-center gap-2" @click.stop>
									<UiButton variant="ghost" size="icon" @click="testUserVoice(user)">
										<Volume2 class="h-4 w-4" />
									</UiButton>
									<UiCheckbox
										:model-value="selectedUsers.has(user.userId)"
										@update:model-value="() => toggleUser(user.userId)"
									/>
								</div>
							</div>
						</UiCardContent>
					</UiCard>
				</div>
			</UiCardContent>
		</UiCard>

		<!-- Delete Confirmation Dialog -->
		<UiAlertDialog v-model:open="showDeleteDialog">
			<UiAlertDialogContent>
				<UiAlertDialogHeader>
					<UiAlertDialogTitle>{{ t('deleteConfirmation.title') }}</UiAlertDialogTitle>
					<UiAlertDialogDescription>
						{{ t('overlays.tts.deleteUsersConfirmation', { count: selectedCount }) }}
					</UiAlertDialogDescription>
				</UiAlertDialogHeader>
				<UiAlertDialogFooter>
					<UiAlertDialogCancel>{{ t('deleteConfirmation.cancel') }}</UiAlertDialogCancel>
					<UiAlertDialogAction @click="handleDelete">
						{{ t('deleteConfirmation.confirm') }}
					</UiAlertDialogAction>
				</UiAlertDialogFooter>
			</UiAlertDialogContent>
		</UiAlertDialog>
	</div>
</template>
