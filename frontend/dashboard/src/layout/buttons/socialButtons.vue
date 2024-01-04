<script setup lang="ts">
import LinkButton from './linkButton.vue';

import DiscordLogo from '@/assets/icons/integrations/discord.svg?component';
import GithubLogo from '@/assets/icons/integrations/github.svg?component';

const GITHUB_URL = 'https://github.com/twirapp/twir';
const DISCORD_URL = 'https://discord.gg/Q9NBZq3zVV';
</script>

<template>
	<link-button :url="GITHUB_URL">
		<github-logo />
	</link-button>
	<link-button :url="DISCORD_URL">
		<discord-logo />
	</link-button>
</template>
