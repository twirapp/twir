<script setup lang="ts">
import { ChannelRolePermissionEnum } from '~/gql/graphql'

definePageMeta({
	layout: 'default',
	noPadding: true,
	neededPermission: ChannelRolePermissionEnum.ViewTimers,
})
</script>

<template>
	<TimersPage />
</template>
