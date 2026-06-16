<script setup lang="ts">
import { computed, onMounted } from 'vue'
import Oauth from '~~/layers/dashboard/components/integrations/variants/oauth.vue'
import {
	useValorantIntegration,
	valorantBroadcaster,
} from '~~/layers/dashboard/features/integrations/composables/valorant/use-valorant-integration.js'

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
		icon="twir-integrations:valorant"
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
