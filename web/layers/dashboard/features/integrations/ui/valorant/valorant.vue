<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

import ValorantIcon from '@/assets/integrations/valorant.svg?use'
import Oauth from '@/components/integrations/variants/oauth.vue'
import {
	useValorantIntegration,
	valorantBroadcaster,
} from '@/features/integrations/composables/valorant/use-valorant-integration.ts'

const { t } = useI18n()

const { userName, avatar, logout, authLink, isDataFetching, refetchData } = useValorantIntegration()

const oauthData = computed(() => {
	if (!userName.value && !avatar.value) {
		return null
	}

	return {
		userName: userName.value,
		avatar: avatar.value,
	}
})

onMounted(async () => {
	valorantBroadcaster.onmessage = async (event) => {
		if (event.data !== 'refresh') return

		await refetchData()
	}
})
</script>

<template>
	<Oauth
		title="Valorant"
		:icon="ValorantIcon"
		icon-width="48px"
		dialog-content-class="w-[600px]"
		:auth-link="authLink"
		:is-loading="isDataFetching"
		:logout="logout"
		:data="oauthData"
	>
		<template #description>
			{{ t('integrations.valorant.description') }}
		</template>
	</Oauth>
</template>
