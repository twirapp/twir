<script setup lang="ts">
import { Toaster as Sonner } from '@/components/ui/sonner'
import { TooltipProvider } from '@/components/ui/tooltip'
</script>

<template>
	<TooltipProvider :delay-duration="100">
		<RouterView />
		<Sonner />
	</TooltipProvider>
</template>
