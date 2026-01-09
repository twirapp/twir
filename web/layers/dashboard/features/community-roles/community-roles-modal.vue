<script setup lang="ts">


import CommunityRoles from './community-roles.vue'

import DialogOrSheet from '#layers/dashboard/components/dialog-or-sheet.vue'


const { t } = useI18n()
</script>

<template>
	<UiDialog>
		<UiDialogTrigger as-child>
			<slot name="trigger" />
		</UiDialogTrigger>

		<DialogOrSheet>
			<UiDialogHeader>
				<UiDialogTitle>
					{{ t('sidebar.roles') }}
				</UiDialogTitle>
			</UiDialogHeader>
			<CommunityRoles />
		</DialogOrSheet>
	</UiDialog>
</template>
