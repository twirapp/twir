<script setup lang="ts">
import { ref } from 'vue'


import { useMutationDropAllAuthSessions } from '#layers/dashboard/api/admin/actions'




const { t } = useI18n()
const dropAllAuthSessions = useMutationDropAllAuthSessions()

async function onDropSessions() {
	await dropAllAuthSessions.executeMutation({})
}

const confirmOpened = ref(false)
</script>

<template>
	<UiCard class="rounded-lg border bg-card text-card-foreground shadow-xs border-red-500 p-4 flex flex-col gap-4">
		<div class="flex items-center">
			<div class="flex-auto">
				<small class="text-sm font-medium leading-none">
					{{ t('adminPanel.adminActions.dangerZone.revokeSessions') }}
				</small>
				<p class="text-sm text-muted-foreground">
					{{ t('adminPanel.adminActions.dangerZone.revokeAllSessionsDescription') }}
				</p>
			</div>
			<UiButton variant="destructive" @click="confirmOpened = true">
				{{ t('adminPanel.adminActions.dangerZone.revoke') }}
			</UiButton>
		</div>
	</UiCard>

	<!-- TODO: reusable action confirm -->
	<UiActionConfirm
		v-model:open="confirmOpened"
		:confirm-text="t('adminPanel.adminActions.dangerZone.revokeAllSessionsConfirm')"
		@confirm="onDropSessions"
	/>
</template>
