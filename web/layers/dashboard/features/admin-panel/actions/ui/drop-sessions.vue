<script setup lang="ts">
import { ref } from 'vue'
import { useMutationDropAllAuthSessions } from '~~/layers/dashboard/api/admin/actions'

import ActionConfirm from '@/components/ui/action-confirm'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'

const { t } = useI18n()
const dropAllAuthSessions = useMutationDropAllAuthSessions()

async function onDropSessions() {
	await dropAllAuthSessions.executeMutation({})
}

const confirmOpened = ref(false)
</script>

<template>
	<Card
		class="bg-card text-card-foreground flex flex-col gap-4 rounded-lg border border-red-500 p-4 shadow-xs"
	>
		<div class="flex items-center">
			<div class="flex-auto">
				<small class="text-sm leading-none font-medium">
					{{ t('adminPanel.adminActions.dangerZone.revokeSessions') }}
				</small>
				<p class="text-muted-foreground text-sm">
					{{ t('adminPanel.adminActions.dangerZone.revokeAllSessionsDescription') }}
				</p>
			</div>
			<Button
				variant="destructive"
				@click="confirmOpened = true"
			>
				{{ t('adminPanel.adminActions.dangerZone.revoke') }}
			</Button>
		</div>
	</Card>

	<!-- TODO: reusable action confirm -->
	<ActionConfirm
		v-model:open="confirmOpened"
		:confirm-text="t('adminPanel.adminActions.dangerZone.revokeAllSessionsConfirm')"
		@confirm="onDropSessions"
	/>
</template>
