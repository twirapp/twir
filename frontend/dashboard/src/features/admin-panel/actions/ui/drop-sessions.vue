<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useMutationDropAllAuthSessions } from '@/api/admin/actions'
import ActionConfirm from '@/components/ui/action-confirm.vue'
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
	<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
		{{ t('adminPanel.adminActions.dangerZone.title') }}
	</h4>

	<Card class="rounded-lg border bg-card text-card-foreground shadow-sm border-red-500 p-4 flex flex-col gap-4">
		<div class="flex items-center">
			<div class="flex-auto">
				<small class="text-sm font-medium leading-none">
					{{ t('adminPanel.adminActions.dangerZone.revokeSessions') }}
				</small>
				<p class="text-sm text-muted-foreground">
					{{ t('adminPanel.adminActions.dangerZone.revokeAllSessionsDescription') }}
				</p>
			</div>
			<Button variant="destructive" @click="confirmOpened = true">
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
