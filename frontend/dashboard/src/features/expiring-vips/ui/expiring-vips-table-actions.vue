<script setup lang="ts">
import { PencilIcon, TrashIcon } from 'lucide-vue-next'

import { Button } from '@/components/ui/button'
</script>

<template>
	<div class="flex items-center gap-2">
		<Button disabled size="icon">
			<PencilIcon class="h-4 w-4" />
		</Button>
		<Button disabled variant="destructive" size="icon">
			<TrashIcon class="h-4 w-4" />
		</Button>
	</div>
</template>
