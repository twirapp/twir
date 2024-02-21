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
} from '@twir/api/messages/integrations_discord/integrations_discord';
import {
	NTabs,
	NTabPane,
	NButton,
	NAvatar,
	NDivider,
	NSelect,
	NSwitch,
	useMessage,
	NMention,
	NPopconfirm,
	NAlert,
	useThemeVars,
	NSpin,
} from 'naive-ui';
import type { SelectBaseOption, SelectOption } from 'naive-ui/es/select/src/interface';
import { type VNode, computed, h, ref, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import WithSettings from './variants/withSettings.vue';
import TwitchMultipleUsersSelector from '../twitchUsers/multiple.vue';

import { useDiscordIntegration, getGuildChannelsFn, useProfile } from '@/api/index.js';
import StreamStarting from '@/assets/images/streamStarting.jpeg?url';
import IconDiscord from '@/assets/integrations/discord.svg?use';

const themeVars = useThemeVars();

const manager = useDiscordIntegration();
const { data: authLink } = manager.getConnectLink();
const {
	data: discordIntegrationData,
	isLoading: isDataLoading,
} = manager.useData();
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

const isConnectDisabled = computed(() => {
	return (discordIntegrationData?.value?.guilds?.length || 0) >= 2;
});

const currentTab = ref<string>();
watch(discordIntegrationData, (v) => {
	if (!v) return;

	if (!currentTab.value) currentTab.value = v.guilds?.at(0)?.name;
	if (!v.guilds.some(g => g.name === currentTab.value)) {
		currentTab.value = v.guilds?.at(0)?.name;
	}

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
	const result: Array<Array<SelectBaseOption>> = [];

	for (let index = 0; index < guildsChannelsResults.length; index++) {
		const channels = guildsChannelsResults.at(index)!;

		const channelsData = channels.data?.channels.filter(c => c.type === ChannelType.TEXT) ?? [];

		result.push(channelsData.map(c => ({
			label: `#${c.name}`,
			value: c.id,
			disabled: !c.canSendMessages,
		})));
	}

	return result;
});

const liveChannelSelectorRenderOption = ({ node, option }: { node: VNode; option: SelectOption }): VNode => {
	if (!option.disabled) return node;

	return h(
		'div',
		{ style: 'display: flex; justify-content: space-between; align-items: center' },
		{
			default: () => [
				node,
				t('integrations.discord.cannotSendMessage'),
			],
		},
	);
};

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
	<with-settings
		title="Discord"
		:save="saveSettings"
		:isLoading="isDataLoading"
		modal-width="80vw"
		:icon="IconDiscord"
		icon-fill="#5865F2"
	>
		<template #description>
			{{ t('integrations.discord.description') }}
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
								:src="`https://cdn.discordapp.com/${guild.id}/${guild.icon}.png`"
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
									{{ t('integrations.discord.alerts.label') }}
								</span>
								<n-divider style="margin: 0; margin-bottom: 5px" />

								<div class="switch">
									<n-switch v-model:value="guild.liveNotificationEnabled" />
									<span>{{ t('sharedTexts.enabled') }}</span>
								</div>

								<div class="switch">
									<n-switch v-model:value="guild.liveNotificationShowTitle" />
									<span>{{ t('integrations.discord.alerts.showTitle') }}</span>
								</div>

								<div class="switch">
									<n-switch v-model:value="guild.liveNotificationShowCategory" />
									<span>{{ t('integrations.discord.alerts.showCategory') }}</span>
								</div>

								<div class="switch">
									<n-switch v-model:value="guild.liveNotificationShowPreview" />
									<span>{{ t('integrations.discord.alerts.showPreview') }}</span>
								</div>

								<div class="switch">
									<n-switch v-model:value="guild.liveNotificationShowProfileImage" />
									<span>{{ t('integrations.discord.alerts.showProfileImage') }}</span>
								</div>

								<div class="switch">
									<n-switch v-model:value="guild.liveNotificationShowViewers" />
									<span>{{ t('integrations.discord.alerts.showViewers') }}</span>
								</div>

								<div class="form-item">
									<span>{{ t('integrations.discord.alerts.channelsSelect') }}</span>
									<n-select
										v-model:value="guild.liveNotificationChannelsIds"
										multiple
										clearable
										filterable
										:options="liveChannelSelectorOptions.at(guildIndex) ?? []"
										:max-tag-count="5"
										:render-option="liveChannelSelectorRenderOption"
										:loading="guildsChannelsResults.some(q => q.isLoading)"
									/>
								</div>

								<n-divider style="margin: 4px;" />

								<div class="form-item">
									<span>{{ t('integrations.discord.alerts.additionalUsersIdsForLiveCheck') }}</span>
									<TwitchMultipleUsersSelector
										v-model="guild.additionalUsersIdsForLiveCheck"
										:max="10"
									/>
									<span v-if="guild.additionalUsersIdsForLiveCheck.length >= 10" class="text-orange-500">
										Maximum additional users selected
									</span>
								</div>

								<div class="form-item">
									<span>{{ t('integrations.discord.alerts.streamOnlineLabel') }}</span>
									<n-mention
										v-model:value="guild.liveNotificationMessage"
										type="textarea"
										:options="getRolesMentionsOptions(guild.id)"
										:placeholder="t('integrations.discord.alerts.streamOnlinePlaceholder')"
										:maxlength="5"
									/>
									<span class="description">{userName}, {displayName}, {title}, {categoryName} – supported variables</span>
								</div>

								<div style="display: flex; flex-direction: column; gap: 8px;">
									<div class="form-item">
										<span>{{ t('integrations.discord.alerts.streamOfflineLabel') }}</span>
										<n-mention
											v-model:value="guild.offlineNotificationMessage"
											type="textarea"
											:disabled="guild.shouldDeleteMessageOnOffline"
											:options="getRolesMentionsOptions(guild.id)"
											:placeholder="t('integrations.discord.alerts.streamOfflinePlaceholder')"
										/>
										<span class="description">{userName}, {displayName} – supported variables</span>
									</div>


									<div class="switch">
										<n-switch v-model:value="guild.shouldDeleteMessageOnOffline" />
										<span>{{ t('integrations.discord.alerts.shouldDeleteMessageOnOffline') }}</span>
									</div>
								</div>

								<n-divider />

								<n-alert type="info">
									{{ t('integrations.discord.alerts.updateAlert') }}
								</n-alert>
							</div>

							<div style="width: 50%">
								<DiscordMessages>
									<DiscordMessage :bot="true" author="TwirApp" avatar="/twir.svg">
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
												footerIcon="/twir.svg"
												borderColor="#6441a5"
												:thumbnail="guild.liveNotificationShowProfileImage ? currentUser?.avatar : null"
												:image="guild.liveNotificationShowPreview ? StreamStarting : null"
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
								{{ t('sharedTexts.dangerZone') }}
							</span>
							<n-divider style="margin: 0; margin-bottom: 5px" />

							<n-popconfirm
								:positive-text="t('deleteConfirmation.confirm')"
								:negative-text="t('deleteConfirmation.cancel')"
								@positive-click="() => disconnectGuildById(guild.id)"
							>
								<template #trigger>
									<n-button type="error" secondary>
										{{ t('integrations.discord.disconnectGuild') }}
									</n-button>
								</template>
								{{ t('deleteConfirmation.text') }}
							</n-popconfirm>
						</div>
					</div>
				</n-tab-pane>

				<template v-if="!discordIntegrationData?.guilds?.length" #prefix>
					{{ t('integrations.discord.noGuilds') }}
				</template>

				<template #suffix>
					<n-button
						:disabled="isConnectDisabled"
						type="success"
						size="small"
						secondary
						@click="connectGuild"
					>
						{{ t('integrations.discord.connectGuild') }}
					</n-button>
				</template>
			</n-tabs>
		</template>

		<template #additionalFooter>
			<div class="profile">
				<n-spin v-if="isDataLoading" style="height: 18px;" />

				<template v-else>
					{{
						t(
							'integrations.discord.connectedGuilds',
							{
								guilds: t(
									'integrations.discord.guildPluralization',
									discordIntegrationData?.guilds?.length ?? 0
								)
							}
						)
					}}
				</template>
			</div>
		</template>
	</with-settings>
</template>

<style scoped>
.form-item {
	display: flex;
	flex-direction: column;
	gap: 4px;
}

.form-item .description {
	font-size: 11px;
}

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

.profile {
	display: flex;
	align-items: center;
	padding: 10px;
	background-color: v-bind('themeVars.buttonColor2');
	border-radius: 4px;
	gap: 8px;
}
</style>
