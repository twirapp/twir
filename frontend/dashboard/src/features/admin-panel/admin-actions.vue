<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import { useMutationDropAllAuthSessions } from '@/api/admin/actions.js'
import { Button } from '@/components/ui/button'

const { t } = useI18n()
const dropAllAuthSessions = useMutationDropAllAuthSessions()

function onDropSessions() {
	// TODO: use confirm dialog from shadcn
	// eslint-disable-next-line no-alert
	if (confirm(t('adminPanel.adminActions.dropAllSessionConfirm'))) {
		dropAllAuthSessions.executeMutation({})
	}
}
</script>

<template>
	<Button variant="destructive" @click="onDropSessions">
		{{ t('adminPanel.adminActions.dropAllSession') }}
	</Button>
</template>
