<script setup lang='ts'>
import { NGrid, NGridItem } from 'naive-ui';

import Discord from '@/components/integrations/discord.vue';
import Donatello from '@/components/integrations/donatello.vue';
import Donatepay from '@/components/integrations/donatepay.vue';
import Donatestream from '@/components/integrations/donatestream.vue';
import Donationalerts from '@/components/integrations/donationalerts.vue';
import Faceit from '@/components/integrations/faceit.vue';
import Lastfm from '@/components/integrations/lastfm.vue';
import Spotify from '@/components/integrations/spotify.vue';
import Streamlabs from '@/components/integrations/streamlabs.vue';
import Valorant from '@/components/integrations/valorant.vue';
import Vk from '@/components/integrations/vk.vue';

// const { t } = useI18n();
</script>

<template>
	<!-- <div class="integrations">
		<n-alert type="info" class="info">
			<span v-html="t('integrations.songServicesInfo')" />
		</n-alert>

		<n-table :bordered="false" :single-line="false" size="large">
			<thead style="display: none">
				<tr>
					<th></th>
					<th></th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				<Spotify />
				<Lastfm />
				<Vk />
				<Streamlabs />
				<Donationalerts />
				<Faceit />
			</tbody>
		</n-table>
	</div> -->

	<div class="integrations">
		<n-grid cols="1 m:2" :x-gap="16" :y-gap="16" responsive="screen">
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
		</n-grid>
	</div>
</template>

<style>
.integrations {
	display: flex;
	align-items: center;
	justify-content: center;
	max-width: 1000px;
	margin: 0 auto;
}

.info {
	width: 50%;
	margin-bottom: 15px
}
</style>
