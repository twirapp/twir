<script setup lang="ts">
import { ChannelRolePermissionEnum } from '~/gql/graphql'

definePageMeta({
	layout: 'default',
	neededPermission: ChannelRolePermissionEnum.ManageOverlays,
})
</script>

<template>
	<DudesOverlayPage />
</template>
