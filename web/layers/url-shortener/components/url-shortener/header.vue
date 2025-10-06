<script setup lang="ts">
import { WandSparklesIcon } from 'lucide-vue-next'
</script>

<template>
	<div class="flex flex-col gap-6 items-center">
		<div class="flex rounded-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,10%)] p-5">
			<WandSparklesIcon fill="rgba(255,255,255,0.3)" class="size-8" />
		</div>
		<div class="flex flex-col gap-4 items-center justify-center">
			<div class="flex flex-col gap-0.5 items-center justify-center">
				<span class="text-4xl font-bold"> Twir </span>
				<span class="text-xl text-[hsl(240,11%,65%)]"> URL Shortener </span>
			</div>
			<span class="text-xl text-center w-3/4">
				Free lightning-fast URL shortening for modern web sharing
			</span>
		</div>
	</div>
</template>
