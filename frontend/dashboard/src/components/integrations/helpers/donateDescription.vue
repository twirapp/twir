<script setup lang="ts">
import { NA } from 'naive-ui';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
</script>

<template>
	<i18n-t
		keypath="integrations.donateServicesInfo"
	>
		<router-link to="/dashboard/events/custom" #="{ navigate, href }" custom>
			<n-a :href="href" @click="navigate">
				{{ t('sidebar.events').toLocaleLowerCase() }}
			</n-a>
		</router-link>

		<router-link to="/dashboard/events/chat-alerts" #="{ navigate, href }" custom>
			<n-a :href="href" @click="navigate">
				{{ t('sidebar.chatAlerts').toLocaleLowerCase() }}
			</n-a>
		</router-link>
	</i18n-t>
</template>
