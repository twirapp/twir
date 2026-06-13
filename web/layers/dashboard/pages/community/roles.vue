<script setup lang="ts">
import CommunityRoles from '~~/layers/dashboard/features/community-roles/community-roles.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', neededPermission: 'ViewRoles' })
</script>

<template>
	<CommunityRoles />
</template>
