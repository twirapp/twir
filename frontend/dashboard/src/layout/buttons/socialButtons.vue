<script setup lang="ts">
import { GITHUB_REPOSITORY_URL, DISCORD_INVITE_URL } from '@twir/brand';

import LinkButton from './linkButton.vue';

import DiscordLogo from '@/assets/integrations/discord.svg?use';
import GithubLogo from '@/assets/integrations/github.svg?use';

</script>

<template>
	<link-button :url="GITHUB_REPOSITORY_URL">
		<github-logo />
	</link-button>
	<link-button :url="DISCORD_INVITE_URL">
		<discord-logo />
	</link-button>
</template>
