<script setup lang="ts">
import OverlayRegistryEdit from '~/layers/dashboard/features/overlay-builder/overlay-registry-edit.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<OverlayRegistryEdit />
</template>
