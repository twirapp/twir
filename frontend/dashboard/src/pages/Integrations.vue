<script setup lang='ts'>
import { NTable, NAlert } from 'naive-ui';
import { useI18n } from 'vue-i18n';

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

const { t } = useI18n();
</script>

<template>
	<div class="integrations">
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
				<Discord />
				<Spotify />
				<Lastfm />
				<Vk />
				<Streamlabs />
				<Donationalerts />
				<Faceit />
				<Donatello />
				<Donatepay />
				<Donatestream />
				<Valorant />
			</tbody>
		</n-table>
	</div>
</template>

<style>
.integrations {
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	height: 100%;
	width: 100%;
}

.integrations table {
	width: 50%;
}

.integrations table tr td:first-of-type {
	width: 70px;
}

.integrations table tr td:last-of-type {
	width: 70px;
}

.integrations table tr td {
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	border-right: 0;
}

.info {
	width: 50%;
	margin-bottom: 15px
}
</style>
