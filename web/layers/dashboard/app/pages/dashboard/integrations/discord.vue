<script setup lang="ts">
import DiscordCallback from '~~/layers/dashboard/app/features/integrations/pages/discord-callback.vue'
import DiscordSettings from '~~/layers/dashboard/app/features/integrations/pages/discord-settings.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })

const route = useRoute()
const isCallback = computed(() => typeof route.query.code === 'string')
</script>

<template>
	<DiscordCallback v-if="isCallback" />
	<DiscordSettings v-else />
</template>
