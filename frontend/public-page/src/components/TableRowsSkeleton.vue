<script setup lang="ts">
import { Skeleton } from '@/components/ui/skeleton';
import { TableCell, TableRow } from '@/components/ui/table';

defineProps<{
	rows: number
	colspan: number
}>();

</script>

<template>
	<TableRow v-for="i in rows" :key="i">
		<TableCell :colspan="colspan" class="w-full font-medium">
			<Skeleton class="w-full h-[20px] rounded-full" />
		</TableCell>
	</TableRow>
</template>

<style scoped>

</style>
