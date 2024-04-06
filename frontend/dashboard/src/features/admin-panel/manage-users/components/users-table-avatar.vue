<script setup lang="ts">
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';

defineProps<{ avatar: string, name: string }>();
</script>

<template>
	<Avatar class="h-9 w-9">
		<AvatarImage :src="avatar" :alt="name" loading="lazy" />
		<AvatarFallback>{{ name.charAt(0).toUpperCase() }}</AvatarFallback>
	</Avatar>
</template>
