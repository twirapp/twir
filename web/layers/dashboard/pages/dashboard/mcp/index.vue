<script setup lang="ts">
import McpGuide from '~~/layers/dashboard/features/mcp-guide/mcp-guide.vue'
import { ChannelRolePermissionEnum } from '~/gql/graphql.js'

definePageMeta({
	layout: 'dashboard',
	middleware: 'auth',
	noPadding: true,
	neededPermission: ChannelRolePermissionEnum.ManageBotSettings,
})
</script>

<template>
	<McpGuide />
</template>
