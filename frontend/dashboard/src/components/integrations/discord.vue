<script setup lang='ts'>
import '@discord-message-components/vue/dist/style.css';

import {
	DiscordMessages,
	DiscordMessage,
	DiscordEmbed,
	DiscordEmbedFields,
	DiscordEmbedField,
	DiscordMention,
	// eslint-disable-next-line @typescript-eslint/ban-ts-comment
	// @ts-ignore
} from '@discord-message-components/vue';
import { useQueries } from '@tanstack/vue-query';
import {
	ChannelType,
	type GetDataResponse,
} from '@twir/grpc/generated/api/api/integrations_discord';
import {
	NTabs,
	NTabPane,
	NButton,
	NAvatar,
	NDivider,
	NSelect,
	NFormItem,
	NSwitch,
	useMessage,
	NMention,
	NPopconfirm,
	NAlert,
} from 'naive-ui';
import { computed, ref, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import TwirCircle from '@/../public/TwirInCircle.svg?url';
import { useDiscordIntegration, getGuildChannelsFn, useProfile } from '@/api/index.js';
import IconDiscord from '@/assets/icons/integrations/discord.svg?component';
import StreamStarting from '@/assets/images/streamStarting.jpeg?url';
import IntegrationWithSettings from '@/components/integrations/variants/withSettings.vue';

const manager = useDiscordIntegration();
const { data: authLink } = manager.getConnectLink();
const {
	data: discordIntegrationData,
	isLoading: isDataLoading,
} = manager.getData();
const guildDisconnect = manager.disconnectGuild();
const updateSettings = manager.updateData();

const formValue = ref<GetDataResponse>({
	guilds: [],
});

const message = useMessage();
const { t } = useI18n();

async function saveSettings() {
	await updateSettings.mutateAsync(formValue.value);
	message.success(t('sharedTexts.saved'));
}

function connectGuild() {
	if (!authLink.value?.link) return;

	window.open(authLink.value.link, 'Twir connect discord', 'width=800,height=600');

	// window.location.replace(authLink.value.link);
}

const currentTab = ref<string>('');
watch(discordIntegrationData, (v) => {
	if (!v) return;
	currentTab.value = v.guilds?.[0]?.name;
	formValue.value = toRaw(v);
}, { immediate: true });

async function disconnectGuildById(guildId: string) {
	const guild = discordIntegrationData.value?.guilds?.find(g => g.id === guildId);
	if (!guild) return;

	await guildDisconnect.mutateAsync(guildId);
	message.success('Disconnected');
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

		result.push(channelsData.map(c => ({ label: `#${c.name}`, value: c.id })));
	}

	return result;
});

function getRolesMentionsOptions(guildId: string) {
	const guild = discordIntegrationData.value?.guilds?.find((g) => g.id === guildId);
	if (!guild) return [];

	return guild.roles.map((r) => ({
		label: r.name,
		value: r.name.replace('@', ''),
	})) ?? [];
}

function getGuildRoleColorByName(guildId: string, roleName: string) {
	const guild = discordIntegrationData.value?.guilds.find(g => g.id === guildId);
	if (!guild) return;
	const role = guild.roles.find(r => r.name === roleName.replace('@', ''));
	if (!role || role.color === '0') return null;

	const hexColor = Number(role.color).toString(16);

	return `#${hexColor.padStart(6, '0')}`;
}

const { data: currentUser } = useProfile();
</script>

<template>
	<integration-with-settings
		name="Discord"
		:save="saveSettings"
		:isLoading="isDataLoading"
		modal-width="80vw"
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
				v-model:value="currentTab"
				closable
				tab-style="min-width: 80px;"
			>
				<n-tab-pane
					v-for="(guild, guildIndex) in formValue.guilds" :key="guild.id"
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

					<div style="display: flex; flex-direction: column; gap: 12px">
						<div class="block">
							<div style="display: flex; flex-direction: column; gap: 8px; width: 50%">
								<span style="font-size: 16px">
									Alerts
								</span>
								<n-divider style="margin: 0; margin-bottom: 5px" />

								<div class="switch">
									<n-switch v-model:value="guild.liveNotificationEnabled" />
									<span>Enabled</span>
								</div>

								<div class="switch">
									<n-switch v-model:value="guild.liveNotificationShowTitle" />
									<span>Show title of stream</span>
								</div>

								<div class="switch">
									<n-switch v-model:value="guild.liveNotificationShowCategory" />
									<span>Show category of stream</span>
								</div>

								<div class="switch">
									<n-switch v-model:value="guild.liveNotificationShowViewers" />
									<span>Show viewers</span>
								</div>

								<n-form-item label="Target channels" style="margin-top: 4px;">
									<n-select
										v-model:value="guild.liveNotificationChannelsIds"
										multiple
										clearable
										filterable
										:options="liveChannelSelectorOptions.at(guildIndex) ?? []"
										:max-tag-count="5"
									/>
								</n-form-item>

								<n-form-item label="Stream online message" style="margin-top: 4px;">
									<n-mention
										v-model:value="guild.liveNotificationMessage"
										type="textarea"
										:options="getRolesMentionsOptions(guild.id)"
										placeholder="The message twir will send when stream started"
										:maxlength="5"
									/>
								</n-form-item>

								<n-form-item label="Stream offline message" style="margin-top: 4px;">
									<n-mention
										v-model:value="guild.offlineNotificationMessage"
										type="textarea"
										:options="getRolesMentionsOptions(guild.id)"
										placeholder="The message twir will send when stream goes offline"
									/>
								</n-form-item>

								<n-alert type="info">
									Twir will update embed periodically to make sure embed in sync with your stream state (viewers, title, category, e.t.c)
								</n-alert>
							</div>

							<div style="width: 50%">
								<DiscordMessages>
									<DiscordMessage :bot="true" author="TwirApp" :avatar="TwirCircle">
										<template v-for="m, _ of guild.liveNotificationMessage.split(' ')" :key="_">
											<DiscordMention
												v-if="m.startsWith('@')"
												type="role"
												:highlight="true"
												:roleColor="getGuildRoleColorByName(guild.id, m.trim())"
											>
												{{ m.replace('@', '') }}
											</DiscordMention>
											<template v-else>
												{{ m }}
											</template>
											{{ ' ' }}
										</template>
										<template #embeds>
											<DiscordEmbed
												embedTitle="Today we are doing amazing things!"
												:url="`https://twitch.tv/${currentUser?.login}`"
												:timestamp="new Date()"
												:footerIcon="TwirCircle"
												borderColor="#6441a5"
												:thumbnail="currentUser?.avatar"
												:image="StreamStarting"
											>
												<template #fields>
													<DiscordEmbedFields>
														<DiscordEmbedField
															v-if="guild.liveNotificationShowTitle"
															fieldTitle="Title"
															inline
														>
															Today we are doing amazing things!
														</DiscordEmbedField>
														<DiscordEmbedField
															v-if="guild.liveNotificationShowViewers"
															fieldTitle="Viewers"
															inline
														>
															5
														</DiscordEmbedField>
														<DiscordEmbedField
															v-if="guild.liveNotificationShowCategory"
															fieldTitle="Category"
														>
															Software and game development
														</DiscordEmbedField>
													</DiscordEmbedFields>
												</template>

												<template #footer>
													TwirApp
												</template>
											</DiscordEmbed>
										</template>
									</DiscordMessage>
								</DiscordMessages>
							</div>
						</div>

						<div class="block" style="display: flex; flex-direction: column; gap: 8px;">
							<span style="font-size: 16px">
								Danger zone
							</span>
							<n-divider style="margin: 0; margin-bottom: 5px" />

							<n-popconfirm
								:positive-text="t('deleteConfirmation.confirm')"
								:negative-text="t('deleteConfirmation.cancel')"
								@positive-click="() => disconnectGuildById(guild.id)"
							>
								<template #trigger>
									<n-button type="error" secondary>Disconnect guild</n-button>
								</template>
								{{ t('deleteConfirmation.text') }}
							</n-popconfirm>
						</div>
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
		</template>
	</integration-with-settings>
</template>

<style scoped>
.switch {
	display: flex;
	gap: 8px;
	align-items: start;
}
.block {
	background-color: rgba(255, 255, 255, 0.06);
	border-radius: 11px;
	padding: 12px;
	display: flex;
	gap: 16px;
}

.discord-embed-image {
	max-height: 300px;
}

.guild-avatar {
	width: 20px;
	height: 20px;
	display: flex;
	align-items: center;
	justify-content: center;
}
</style>
