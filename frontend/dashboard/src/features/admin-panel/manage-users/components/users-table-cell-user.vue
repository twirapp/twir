<script setup lang="ts">
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';

defineProps<{ name: string, avatar: string }>();
</script>

<template>
	<div class="flex items-center gap-4 max-sm:justify-start">
		<Avatar class="h-9 w-9">
			<AvatarImage :src="avatar" :alt="name" loading="lazy" />
			<AvatarFallback>{{ name.charAt(0).toUpperCase() }}</AvatarFallback>
		</Avatar>
		{{ name }}
	</div>
</template>
