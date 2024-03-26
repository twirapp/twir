<script setup lang="ts">
import { NGrid, NGridItem } from 'naive-ui';

import { responsiveCols } from '@/components/consants.js';
import Discord from '@/components/integrations/discord.vue';
import Donatello from '@/components/integrations/donatello.vue';
import Donatepay from '@/components/integrations/donatepay.vue';
import Donatestream from '@/components/integrations/donatestream.vue';
import Donationalerts from '@/components/integrations/donationalerts.vue';
import Faceit from '@/components/integrations/faceit.vue';
import Lastfm from '@/components/integrations/lastfm.vue';
import Nightbot from '@/components/integrations/nightbot.vue';
import SevenTv from '@/components/integrations/seven-tv.vue';
import Spotify from '@/components/integrations/spotify.vue';
import Streamlabs from '@/components/integrations/streamlabs.vue';
import Valorant from '@/components/integrations/valorant.vue';
import Vk from '@/components/integrations/vk.vue';
</script>

<template>
	<div class="integrations">
		<n-grid :cols="responsiveCols" :x-gap="16" :y-gap="16" responsive="screen">
			<n-grid-item :span="1">
				<SevenTv />
			</n-grid-item>
			<n-grid-item :span="1">
				<Discord />
			</n-grid-item>
			<n-grid-item :span="1">
				<Spotify />
			</n-grid-item>
			<n-grid-item :span="1">
				<Lastfm />
			</n-grid-item>
			<n-grid-item :span="1">
				<Vk />
			</n-grid-item>
			<n-grid-item :span="1">
				<Donationalerts />
			</n-grid-item>
			<n-grid-item :span="1">
				<Streamlabs />
			</n-grid-item>
			<n-grid-item :span="1">
				<Donatello />
			</n-grid-item>
			<n-grid-item :span="1">
				<Donatepay />
			</n-grid-item>
			<n-grid-item :span="1">
				<Donatestream />
			</n-grid-item>
			<n-grid-item :span="1">
				<Faceit />
			</n-grid-item>
			<n-grid-item :span="1">
				<Valorant />
			</n-grid-item>
			<n-grid-item :span="1">
				<Nightbot />
			</n-grid-item>
		</n-grid>
	</div>
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
