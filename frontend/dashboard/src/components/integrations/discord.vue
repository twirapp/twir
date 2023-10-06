<script setup lang='ts'>
import { useDiscordIntegration } from '@/api/index.js';
import IconDiscord from '@/assets/icons/integrations/discord.svg?component';
import OauthComponent from '@/components/integrations/variants/oauth.vue';

const manager = useDiscordIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const { data: authLink } = manager.useAuthLink();
</script>

<template>
	<oauth-component
		name="Discord"
		:data="data"
		:logout="() => logout.mutateAsync({})"
		:authLink="authLink?.link"
		:icon="IconDiscord"
		icon-color="#5865F2"
	/>
</template>
