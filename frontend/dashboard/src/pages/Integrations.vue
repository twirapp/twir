<script setup lang="ts">
import Discord from '@/components/integrations/discord.vue'
import Donatello from '@/components/integrations/donatello.vue'
import Donatepay from '@/components/integrations/donatepay.vue'
import Donatestream from '@/components/integrations/donatestream.vue'
import Donationalerts from '@/components/integrations/donationalerts.vue'
import Faceit from '@/components/integrations/faceit.vue'
import Lastfm from '@/components/integrations/lastfm.vue'
import SevenTv from '@/components/integrations/seven-tv.vue'
import Spotify from '@/components/integrations/spotify.vue'
import Streamlabs from '@/components/integrations/streamlabs.vue'
import Valorant from '@/components/integrations/valorant.vue'
import Vk from '@/components/integrations/vk.vue'
import PageLayout from '@/layout/page-layout.vue'
</script>

<template>
	<PageLayout>
		<template #title>
			Integrations
		</template>

		<template #content>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<SevenTv />
				<Discord />
				<Spotify />
				<Lastfm />
				<Vk />
				<Donationalerts />
				<Streamlabs />
				<Donatello />
				<Donatepay />
				<Donatestream />
				<Faceit />
				<Valorant />
			</div>
		</template>
	</PageLayout>
</template>

<style>
.integrations {
	display: flex;
	align-items: center;
	justify-content: center;
	max-width: 60vw;
	margin: 0 auto;
}

@media screen and (max-width: 1024px) {
	.integrations {
		max-width: 80vw;
	}
}
</style>
