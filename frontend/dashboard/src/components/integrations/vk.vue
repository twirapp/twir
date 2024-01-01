<script setup lang="ts">
import { useVKIntegration } from '@/api/index.js';
import IconVk from '@/assets/icons/integrations/vk.svg?component';
import SongDescription from '@/components/integrations/helpers/songDescription.vue';
import OauthComponent from '@/components/integrations/variants/oauth.vue';

const manager = useVKIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const { data: authLink } = manager.useAuthLink();
</script>

<template>
	<oauth-component
		title="VK"
		:data="data"
		:logout="() => logout.mutateAsync({})"
		:authLink="authLink?.link"
		:icon="IconVk"
	>
		<template #description>
			<song-description />
		</template>
	</oauth-component>
</template>
