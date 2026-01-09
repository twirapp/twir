<script setup lang="ts">
import { ChannelRolePermissionEnum } from '~/gql/graphql'

definePageMeta({
	layout: 'default',
	neededPermission: ChannelRolePermissionEnum.ViewOverlays,
})
</script>

<template>
	<OverlaysPage />
</template>
