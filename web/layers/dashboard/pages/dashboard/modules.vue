<script setup lang="ts">
import { ChannelRolePermissionEnum } from '~/gql/graphql'

definePageMeta({
	layout: 'default',
	noPadding: true,
	neededPermission: ChannelRolePermissionEnum.ViewModules,
})
</script>

<template>
	<ModulesPage />
</template>
