<script setup lang="ts">
import { ref } from 'vue'

import { Button } from '@/components/ui/button'

import SecretDialog from './secret-dialog.vue'

const isOpen = ref(false)
</script>

<template>
	<Button @click="isOpen = true">
		<Icon
			name="lucide:plus"
			class="mr-2 h-4 w-4"
		/>
		Create Secret
	</Button>

	<SecretDialog
		v-model:open="isOpen"
		:secret="null"
	/>
</template>
