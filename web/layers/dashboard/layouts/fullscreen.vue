<script setup lang="ts">
import { Toaster } from '@/components/ui/sonner'
</script>

<template>
	<div class="w-full h-full">
		<slot />
		<Toaster />
	</div>
</template>
