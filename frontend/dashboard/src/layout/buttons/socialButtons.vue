<script setup lang="ts">
import { GITHUB_REPOSITORY_URL, DISCORD_INVITE_URL } from '@twir/brand';

import LinkButton from './linkButton.vue';

import DiscordLogo from '@/assets/icons/integrations/discord.svg?component';
import GithubLogo from '@/assets/icons/integrations/github.svg?component';

</script>

<template>
	<link-button :url="GITHUB_REPOSITORY_URL">
		<github-logo />
	</link-button>
	<link-button :url="DISCORD_INVITE_URL">
		<discord-logo />
	</link-button>
</template>
