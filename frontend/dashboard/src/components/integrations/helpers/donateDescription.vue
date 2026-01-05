<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>

<template>
	<i18n-t keypath="integrations.donateServicesInfo">
		<RouterLink to="/dashboard/events/custom" class="text-primary hover:underline">
			{{ t('sidebar.events').toLocaleLowerCase() }}
		</RouterLink>

		<RouterLink to="/dashboard/events/chat-alerts" class="text-primary hover:underline">
			{{ t('sidebar.chatAlerts').toLocaleLowerCase() }}
		</RouterLink>
	</i18n-t>
</template>
