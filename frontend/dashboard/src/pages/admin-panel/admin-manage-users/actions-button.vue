<script setup lang="ts">
import { BanIcon, ShieldIcon } from 'lucide-vue-next';

import { Button } from '@/components/ui/button';

defineProps<{
	userId: string
	isBanned: boolean
	isAdmin: boolean
}>();
</script>

<template>
	<div class="flex items-center gap-2">
		<Button
			:class="{ 'bg-red-600 hover:bg-red-600/90': isBanned }"
			variant="secondary"
			size="icon"
		>
			<BanIcon class="h-4 w-4" />
		</Button>
		<Button
			:class="{ 'bg-green-600 hover:bg-green-600/90': isAdmin }"
			variant="secondary"
			size="icon"
		>
			<ShieldIcon class="h-4 w-4" />
		</Button>
	</div>
</template>
