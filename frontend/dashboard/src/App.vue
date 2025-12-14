<script setup lang="ts">
import { provideClient } from '@urql/vue'
import { ref } from 'vue'
import { RouterView, useRouter } from 'vue-router'
import { Loader } from 'lucide-vue-next'

import { urqlClient } from './plugins/urql.js'

const isRouterReady = ref(false)
const router = useRouter()
router.isReady().finally(() => (isRouterReady.value = true))

provideClient(urqlClient)
</script>

<template>
	<div v-if="!isRouterReady" class="flex justify-center items-center h-full">
		<Loader class="h-10 w-10 sonner-spinner" />
	</div>
	<router-view v-else />
</template>
