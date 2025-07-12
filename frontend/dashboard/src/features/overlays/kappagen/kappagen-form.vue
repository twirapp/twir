<script setup lang="ts">
import EventsTab from '@/features/overlays/kappagen/components/events-tab.vue'
import SettingsTab from '@/features/overlays/kappagen/components/settings-tab.vue'
import AnimationsTab from '@/features/overlays/kappagen/components/animations-tab.vue'
</script>

<template>
	<div class="flex flex-col gap-4">
		<SettingsTab />
		<EventsTab />
		<AnimationsTab />
	</div>
</template>
