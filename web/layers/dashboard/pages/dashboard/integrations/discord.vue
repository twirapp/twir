<script setup lang="ts">
import DiscordCallback from '~~/layers/dashboard/features/integrations/pages/discord-callback.vue'
import DiscordSettings from '~~/layers/dashboard/features/integrations/pages/discord-settings.vue'
import { useRoute } from 'vue-router'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })

const route = useRoute<'dashboard-integrations-discord'>()
const isCallback = computed(() => typeof route.query.code === 'string')
</script>

<template>
	<DiscordCallback v-if="isCallback" />
	<DiscordSettings v-else />
</template>
