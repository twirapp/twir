<script setup lang="ts">
import Donatello from '~~/layers/dashboard/components/integrations/donatello.vue'
import Donatepay from '~~/layers/dashboard/components/integrations/donatepay.vue'
import Donatestream from '~~/layers/dashboard/components/integrations/donatestream.vue'
import Donationalerts from '~~/layers/dashboard/components/integrations/donationalerts.vue'
import Faceit from '~~/layers/dashboard/components/integrations/faceit.vue'
import Spotify from '~~/layers/dashboard/components/integrations/spotify.vue'
import Lastfm from '~~/layers/dashboard/features/integrations/ui/lastfm/lastfm.vue'
import Streamlabs from '~~/layers/dashboard/components/integrations/streamlabs.vue'
import Vk from '~~/layers/dashboard/features/integrations/ui/vk/vk.vue'
import Discord from '~~/layers/dashboard/features/integrations/ui/discord/card.vue'
import SevenTv from '~~/layers/dashboard/features/integrations/ui/seventv/seven-tv.vue'
import Valorant from '~~/layers/dashboard/features/integrations/ui/valorant/valorant.vue'
import PageLayout from '~~/layers/dashboard/layout/page-layout.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<PageLayout>
		<template #title> Integrations </template>

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
