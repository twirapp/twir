<script setup lang="ts">
import StreamersList from './streamers-list.vue'
import { getStreamers } from './streamers.js'

const streamers = await getStreamers()
</script>

<template>
	<section id="streamers" class="container">
		<div class="purple-gradient">
			<div class="container py-24 px-5">
				<h1 class="font-semibold text-[#B0ADFF] text-base uppercase tracking-wider mb-3 text-center">
					Streamers
				</h1>
				<div class="flex justify-between gap-4 flex-col">
					<h3 class="text-white text-4xl font-bold text-center">
						Community around Twir
					</h3>
				</div>
			</div>

			<ClientOnly>
				<StreamersList :streamers />
			</ClientOnly>
		</div>
	</section>
</template>
