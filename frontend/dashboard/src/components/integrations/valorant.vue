<script setup lang="ts">
import { useValorantIntegration } from '@/api/index.js'
import IconValorant from '@/assets/integrations/valorant.svg?use'
import OauthComponent from '@/components/integrations/variants/oauth.vue'

const manager = useValorantIntegration()
const { data } = manager.useData()
const logout = manager.useLogout()
const { data: authLink } = manager.useAuthLink()
</script>

<template>
	<OauthComponent
		title="Valorant"
		:data="data"
		:logout="() => logout.mutateAsync({})"
		:authLink="authLink?.link"
		:icon="IconValorant"
	>
		<template #description>
			<i18n-t keypath="integrations.valorant.info">
				<b class="variable">$(valorant.profile.elo)</b>
				<b class="variable">$(valorant.profile.tier)</b>
			</i18n-t>
		</template>
	</OauthComponent>
</template>

<style scoped>
.variable {
	@apply text-[color:var(--n-color-target)];
}
</style>
