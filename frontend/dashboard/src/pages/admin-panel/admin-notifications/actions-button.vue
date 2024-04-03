<script setup lang="ts">
import { TrashIcon, PencilIcon } from 'lucide-vue-next';

import { Button } from '@/components/ui/button';

defineProps<{ notificationId: number }>();

</script>

<template>
	<div class="flex items-center gap-2">
		<Button variant="destructive" size="icon">
			<TrashIcon class="h-4 w-4" />
		</Button>
		<Button variant="secondary" size="icon">
			<PencilIcon class="h-4 w-4" />
		</Button>
	</div>
</template>
