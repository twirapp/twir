<script setup lang="ts">
import McpConsent from '~~/layers/dashboard/features/mcp-consent/mcp-consent.vue'
import { ChannelRolePermissionEnum } from '~/gql/graphql.js'

definePageMeta({
	layout: 'dashboard',
	middleware: 'auth',
	noPadding: true,
	neededPermission: ChannelRolePermissionEnum.ManageBotSettings,
})
</script>

<template>
	<McpConsent />
</template>
