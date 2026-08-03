<script setup lang="ts">
import CommunityRoles from '~~/layers/dashboard/features/community-roles/community-roles.vue'
import { ChannelRolePermissionEnum } from '~/gql/graphql.js'

definePageMeta({
	layout: 'dashboard',
	middleware: 'auth',
	neededPermission: ChannelRolePermissionEnum.ViewRoles,
	noPadding: true,
})
</script>

<template>
	<CommunityRoles />
</template>
