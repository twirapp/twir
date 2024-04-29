<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useMutationDropAllAuthSessions } from '@/api/admin/actions'
import ActionConfirm from '@/components/ui/action-confirm.vue'
import { Button } from '@/components/ui/button'

const { t } = useI18n()
const dropAllAuthSessions = useMutationDropAllAuthSessions()

async function onDropSessions() {
	await dropAllAuthSessions.executeMutation({})
}

const confirmOpened = ref(false)
</script>

<template>
	<Button variant="destructive" @click="confirmOpened = true">
		{{ t('adminPanel.adminActions.dropAllSession') }}
	</Button>

	<ActionConfirm
		v-model:open="confirmOpened"
		:confirm-text="t('adminPanel.adminActions.dropAllSessionConfirm')"
		@confirm="onDropSessions"
	/>
</template>

<style scoped>

</style>
