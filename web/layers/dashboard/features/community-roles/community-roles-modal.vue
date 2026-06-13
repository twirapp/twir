<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import CommunityRoles from './community-roles.vue'

import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import { Dialog, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'

const { t } = useI18n()
</script>

<template>
	<Dialog>
		<DialogTrigger as-child>
			<slot name="trigger" />
		</DialogTrigger>

		<DialogOrSheet>
			<DialogHeader>
				<DialogTitle>
					{{ t('sidebar.roles') }}
				</DialogTitle>
			</DialogHeader>
			<CommunityRoles />
		</DialogOrSheet>
	</Dialog>
</template>
