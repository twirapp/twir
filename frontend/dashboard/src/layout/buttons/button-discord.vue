<script setup lang="ts">
import { DISCORD_INVITE_URL } from '@twir/brand'

import ButtonLink from './button-link.vue'

import DiscordLogo from '@/assets/integrations/discord.svg?use'
</script>

<template>
	<ButtonLink :url="DISCORD_INVITE_URL">
		<DiscordLogo />
	</ButtonLink>
</template>
