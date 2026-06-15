<script setup lang="ts">
import OverlayRegistryEdit from '~~/layers/dashboard/app/components/registry/overlays/edit.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<OverlayRegistryEdit />
</template>
