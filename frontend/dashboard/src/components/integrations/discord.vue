<script setup lang='ts'>
import { useQueries } from '@tanstack/vue-query';
import { ChannelType, type GetDataResponse } from '@twir/grpc/generated/api/api/integrations_discord';
import {
	NTabs,
	NTabPane,
	NButton,
	NAvatar,
	NDivider,
	NSelect,
	NAlert,
	NFormItem,
	NSwitch,
} from 'naive-ui';
import { computed, onMounted, ref, toRaw, watch } from 'vue';

import { useDiscordIntegration, getGuildChannelsFn } from '@/api/index.js';
import IconDiscord from '@/assets/icons/integrations/discord.svg?component';
import IntegrationWithSettings from '@/components/integrations/variants/withSettings.vue';

const manager = useDiscordIntegration();
const { data: authLink } = manager.getConnectLink();
const { data: discordIntegrationData, isError: isDataError } = manager.getData();
const guildDisconnect = manager.disconnectGuild();

const formValue = ref<GetDataResponse>({
	guilds: [],
});

function connectGuild() {
	if (!authLink.value?.link) return;

	window.location.replace(authLink.value.link);
}

const currentTab = ref<string | undefined>();
watch(discordIntegrationData, (v) => {
	if (!v) return;
	currentTab.value = v.guilds?.[0]?.name;
	formValue.value = toRaw(v);
});
onMounted(() => {
	currentTab.value = discordIntegrationData.value?.guilds?.[0]?.name;
});

async function disconnectGuildByName(guildName: string) {
	const guild = discordIntegrationData.value?.guilds?.find((g: any) => g.name === guildName);
	if (!guild) return;

	await guildDisconnect.mutateAsync(guild.id);
}

const guildsChannelsQueries = computed(() => {
	return discordIntegrationData.value?.guilds?.map((guild: any) => ({
		queryKey: ['discord', 'guild', guild.id, 'channels'],
		queryFn: () => getGuildChannelsFn(guild.id),
	})) ?? [];
});

const guildsChannelsResults = useQueries({
	queries: guildsChannelsQueries,
});

const liveChannelSelectorOptions = computed(() => {
	const result: Array<Array<{ label: string, value: string }>> = [];

	for (let index = 0; index < guildsChannelsResults.length; index++) {
		const channels = guildsChannelsResults.at(index)!;

		const channelsData = channels.data?.channels.filter(c => c.type === ChannelType.TEXT) ?? [];

		result.push(channelsData.map(c => ({ label: c.name, value: c.id })));
	}

	return result;
});
</script>

<template>
	<integration-with-settings
		name="Discord"
	>
		<template #icon>
			<IconDiscord style="width: 30px; fill: #5865F2; display: flex" />
		</template>

		<template #content>
			<div style="display: flex; flex-direction: column">
				<span>Use this integration for setup live alerts into your discord guilds</span>
				<span style="font-size: 11px">
					Connected {{ discordIntegrationData?.guilds?.length }} guilds
				</span>
			</div>
		</template>

		<template #settings>
			<n-tabs
				v-if="!isDataError"
				v-model:value="currentTab"
				closable
				tab-style="min-width: 80px;"
				@close="disconnectGuildByName"
			>
				<n-tab-pane
					v-for="(guild, guildIndex) in discordIntegrationData?.guilds" :key="guild.id"
					:name="guild.name"
				>
					<template #tab>
						<div style="display: flex; gap: 5px; align-items: center; justify-content: center">
							<n-avatar
								round
								:src="`https://cdn.discordapp.com/icons/${guild.id}/${guild.icon}.png`"
								class="guild-avatar"
								:render-fallback="() => guild.name.charAt(0)"
							/>
							<span>
								{{ guild.name }}
							</span>
						</div>
					</template>

					<div class="block">
						<span style="font-size: 16px">
							Alerts
						</span>
						<n-divider style="margin: 0; margin-bottom: 5px" />

						<div style="display: flex; gap: 8px; align-items: start;">
							<n-switch v-model:value="formValue.guilds.at(guildIndex)!.liveNotificationEnabled" />
							<span>Enabled</span>
						</div>

						<n-form-item label="Target channels" style="margin-top: 4px;">
							<n-select
								v-model:value="formValue.guilds.at(guildIndex)!.liveNotificationChannelsIds"
								multiple
								clearable
								filterable
								:options="liveChannelSelectorOptions.at(guildIndex) ?? []"
							/>
						</n-form-item>
					</div>
				</n-tab-pane>
				<template v-if="!discordIntegrationData?.guilds?.length" #prefix>
					No guilds connected
				</template>

				<template #suffix>
					<n-button
						:disabled="(discordIntegrationData?.guilds?.length || 0) >= 2"
						type="success"
						size="small"
						secondary
						@click="connectGuild"
					>
						+ connect
					</n-button>
				</template>
			</n-tabs>

			<n-alert v-else>
				no guilds
			</n-alert>
		</template>
	</integration-with-settings>
</template>

<style scoped>
.block {
	background-color: rgba(255, 255, 255, 0.06);
	border-radius: 11px;
	padding: 12px;
	display: flex;
	flex-direction: column;
	gap: 8px;
}

.guild-avatar {
	width: 20px;
	height: 20px;
	display: flex;
	align-items: center;
	justify-content: center;
}
</style>
