<script setup lang="ts">
import { NButton, useThemeVars } from 'naive-ui';
import { computed } from 'vue';

import DiscordLogo from '@/assets/icons/integrations/discord.svg?component';

const themeVars = useThemeVars();
const discordIconColor = computed(() => themeVars.value.textColor2);

const TWIR_APP_DISCORD_INVITE = 'https://discord.gg/Q9NBZq3zVV';
const openDiscordInvite = () => window.open(TWIR_APP_DISCORD_INVITE, '_blank');
</script>

<template>
	<n-button
		quaternary
		circle
		style="padding: 5px"
		@click="openDiscordInvite"
	>
		<DiscordLogo :style="{ width: '20px', fill: discordIconColor }" />
	</n-button>
</template>
