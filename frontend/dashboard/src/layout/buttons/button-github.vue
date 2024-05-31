<script setup lang="ts">
import { GITHUB_REPOSITORY_URL } from '@twir/brand'

import ButtonLink from './button-link.vue'

import GithubLogo from '@/assets/integrations/github.svg?use'
</script>

<template>
	<ButtonLink :url="GITHUB_REPOSITORY_URL">
		<GithubLogo />
	</ButtonLink>
</template>
