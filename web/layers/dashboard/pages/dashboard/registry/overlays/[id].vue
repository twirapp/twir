<script setup lang="ts">
import OverlayRegistryEdit from '~~/layers/dashboard/components/registry/overlays/edit.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<OverlayRegistryEdit />
</template>
