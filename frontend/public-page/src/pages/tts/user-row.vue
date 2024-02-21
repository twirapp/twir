<script setup lang="ts">
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { TableCell, TableRow } from '@/components/ui/table';

defineProps<{
	name: string,
	avatar: string,
	pitch: number,
	rate: number,
	voice: string
}>();
</script>

<template>
	<TableRow>
		<TableCell>
			<Avatar>
				<AvatarImage :src="avatar" alt="user-avatar" />
				<AvatarFallback>{{ name.slice(0, 2) }}</AvatarFallback>
			</Avatar>
		</TableCell>
		<TableCell>
			<a :href="`https://twitch.tv/${name}`" target="_blank" class="hover:underline">
				{{ name }}
			</a>
		</TableCell>
		<TableCell class="text-center">
			{{ voice }}
		</TableCell>
		<TableCell class="text-center">
			{{ rate }}
		</TableCell>
		<TableCell class="text-center">
			{{ pitch }}
		</TableCell>
	</TableRow>
</template>

<style scoped>

</style>
