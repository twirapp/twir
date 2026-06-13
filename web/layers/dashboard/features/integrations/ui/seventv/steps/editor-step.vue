<script setup lang="ts">
import { useSeventvData } from '@/features/integrations/composables/seventv/use-seventv-data.ts'

const { botProfile, userProfile } = useSeventvData()
</script>

<template>
	Go to <a class="underline cursor-pointer" :href="`https://7tv.app/users/${userProfile?.id}/editors`">https://7tv.app/users/{{ userProfile?.id }}/editors</a>
	and add <b>{{ botProfile?.username }}</b> as editor
</template>
