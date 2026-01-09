<script setup lang="ts">
import { toast } from 'vue-sonner'

import RewardsSelector from '../rewardsSelector.vue'

import type { SongRequestsSettingsOpts } from '@/gql/graphql'

import { useCommandsApi } from '#layers/dashboard/api/commands/commands'
import { useSongRequestsApi } from '#layers/dashboard/api/song-requests'
import TwitchSearchUsers from '#layers/dashboard/components/twitchUsers/twitch-users-select.vue'
import CommandsList from '#layers/dashboard/features/commands/ui/list.vue'
import { SongRequestsSearchChannelOrVideoOptsType } from '@/gql/graphql'

const props = defineProps<{
	open?: boolean
}>()

const emit = defineEmits<{
	'update:open': [value: boolean]
}>()

const { t } = useI18n()

const isOpen = ref(props.open ?? false)

watch(
	() => props.open,
	(v) => {
		if (v !== undefined) isOpen.value = v
	}
)

watch(isOpen, (v) => {
	emit('update:open', v)
})

const youtubeModuleManager = useSongRequestsApi()
const youtubeModuleData = youtubeModuleManager.useSongRequestQuery()
const youtubeModuleSettings = computed(() => {
	return youtubeModuleData.data?.value?.songRequests
})
const youtubeModuleUpdater = youtubeModuleManager.useSongRequestMutation()

const formValue = ref<SongRequestsSettingsOpts>({
	enabled: false,
	acceptOnlyWhenOnline: true,
	takeSongFromDonationMessages: false,
	playerNoCookieMode: false,
	channelPointsRewardId: '',
	maxRequests: 500,
	announcePlay: true,
	neededVotesForSkip: 30,
	denyList: {
		artistsNames: [],
		songs: [],
		users: [],
		channels: [],
		words: [],
	},
	song: {
		maxLength: 10,
		minLength: 0,
		minViews: 50000,
		acceptedCategories: [],
	},
	user: {
		maxRequests: 20,
		minWatchTime: 0,
		minFollowTime: 0,
		minMessages: 0,
	},
	translations: {
		notEnabled: 'Song requests not enabled.',
		nowPlaying: 'Now playing "{{songTitle}}" {{songLink}} requested by @{{orderedByDisplayName}}',
		noText: 'You should provide text for song request.',
		acceptOnlyWhenOnline: 'Requests accepted only on online streams.',
		song: {
			notFound: 'Song not found.',
			alreadyInQueue: 'Song already in queue.',
			ageRestrictions: 'Age restriction on that song.',
			cannotGetInformation: 'Cannot get information about song.',
			live: 'Seems like that song is live, which is disallowed.',
			denied: 'That song is denied for requests.',
			requestedMessage:
				'Song "{{songTitle}}" requested, queue position {{position}}. Estimated wait time before your track will be played is {{waitTime}}.',
			maximumOrdered: 'Maximum number of songs is queued ({{maximum}}).',
			minViews:
				"Song {{songTitle}} ({{songViews}} views) haven't {{neededViews}} views for being ordered",
			maxLength: 'Maximum length of song is {{maxLength}}',
			minLength: 'Minimum length of song is {{minLength}}',
		},
		user: {
			denied: 'You are denied to request any song.',
			maxRequests: 'Maximum number of songs ordered by you ({{count}})',
			minMessages:
				'You have only {{userMessages}} messages, but needed {{neededMessages}} for requesting song',
			minWatched:
				"You've only watched {{userWatched}} but needed {{neededWatched}} to request a song.",
			minFollow:
				'You are followed for {{userFollow}} minutes, but needed {{neededFollow}} for requesting song',
		},
		channel: {
			denied: 'That channel is denied for requests.',
		},
	},
})

watch(
	youtubeModuleSettings,
	async (v) => {
		if (!v) return
		formValue.value = toRaw(v)
	},
	{ immediate: true }
)

async function save() {
	const data = unref(formValue)
	await youtubeModuleUpdater.executeMutation({ opts: data })
	toast.success(t('sharedTexts.saved'))
	isOpen.value = false
}

const channelsSearchValue = ref('')
const channelsSearchDebounced = useDebounce(channelsSearchValue, 500)

const channelsSearchOpts = computed(() => {
	return {
		query: [...formValue!.value!.denyList!.channels, channelsSearchDebounced.value],
		type: SongRequestsSearchChannelOrVideoOptsType.Channel,
	}
})
const channelsSearch = youtubeModuleManager.useYoutubeVideoOrChannelSearch(channelsSearchOpts)

const channelsOptions = computed(() => {
	return (
		channelsSearch?.data.value?.songRequestsSearchChannelOrVideo.items.map((channel) => {
			return {
				label: channel.title,
				value: channel.id,
				image: channel.thumbnail,
			}
		}) ?? []
	)
})

const songsSearchValue = ref('')
const songsSearchDebounced = useDebounce(songsSearchValue, 500)

const songsSearchOpts = computed(() => {
	return {
		query: [...formValue!.value!.denyList!.songs, songsSearchDebounced.value],
		type: SongRequestsSearchChannelOrVideoOptsType.Video,
	}
})

const songsSearch = youtubeModuleManager.useYoutubeVideoOrChannelSearch(songsSearchOpts)
const songsSearchOptions = computed(() => {
	return (
		songsSearch?.data.value?.songRequestsSearchChannelOrVideo.items.map((channel) => {
			return {
				label: channel.title,
				value: channel.id,
				image: channel.thumbnail,
			}
		}) ?? []
	)
})

const commandsManager = useCommandsApi()
const { data: commands } = commandsManager.useQueryCommands()
const srCommands = computed(() => {
	return (
		commands.value?.commands.filter((c) => c.module === 'SONGS' && c.defaultName !== 'song') ?? []
	)
})

function findChannelLabel(id: string) {
	return channelsOptions.value.find((o) => o.value === id)?.label ?? id
}

function findChannelImage(id: string): string {
	return channelsOptions.value.find((o) => o.value === id)?.image ?? ''
}

function findSongLabel(id: string) {
	return songsSearchOptions.value.find((o) => o.value === id)?.label ?? id
}

function findSongImage(id: string): string {
	return songsSearchOptions.value.find((o) => o.value === id)?.image ?? ''
}
</script>

<template>
	<UiDialog v-model:open="isOpen">
		<UiDialogContent class="sm:max-w-[700px] max-h-[85vh] p-0">
			<UiDialogHeader class="p-6 pb-0">
				<DialogTitle>{{ t('sharedTexts.settings') }}</DialogTitle>
			</UiDialogHeader>

			<UiScrollArea class="max-h-[calc(85vh-140px)]">
				<div class="p-6 pt-4">
					<UiTabs default-value="general" class="w-full">
						<UiTabsList class="w-full grid grid-cols-5">
							<UiTabsTrigger value="general">{{ t('songRequests.tabs.general') }}</UiTabsTrigger>
							<UiTabsTrigger value="commands">{{ t('commands.name') }}</UiTabsTrigger>
							<UiTabsTrigger value="users">{{ t('songRequests.tabs.users') }}</UiTabsTrigger>
							<UiTabsTrigger value="songs">{{ t('songRequests.tabs.songs') }}</UiTabsTrigger>
							<UiTabsTrigger value="translations">{{
								t('songRequests.tabs.translations')
							}}</UiTabsTrigger>
						</UiTabsList>

						<UiTabsContent value="general" class="mt-4 space-y-4">
							<div class="space-y-3">
								<div class="flex flex-row items-center gap-4 rounded-lg border p-4">
									<UiSwitch v-model="formValue.enabled" />
									<div class="space-y-0.5 flex-1">
										<Label class="text-base font-medium">{{ t('sharedTexts.enabled') }}</Label>
									</div>
								</div>

								<div class="flex flex-row items-center gap-4 rounded-lg border p-4">
									<UiSwitch v-model="formValue.takeSongFromDonationMessages" />
									<div class="space-y-0.5 flex-1">
										<UiLabel class="text-base font-medium">{{
											t('songRequests.settings.takeSongFromDonationMessage')
										}}</UiLabel>
										<p class="text-sm text-muted-foreground">
											{{ t('songRequests.settings.takeSongFromDonationMessageDescription') }}
										</p>
									</div>
								</div>

								<div class="flex flex-row items-center gap-4 rounded-lg border p-4">
									<UiSwitch v-model="formValue.acceptOnlyWhenOnline" />
									<div class="space-y-0.5 flex-1">
										<UiLabel class="text-base font-medium">{{
											t('songRequests.settings.onlineOnly')
										}}</UiLabel>
										<p class="text-sm text-muted-foreground">
											{{ t('songRequests.settings.onlineOnlyDescription') }}
										</p>
									</div>
								</div>

								<div class="flex flex-row items-center gap-4 rounded-lg border p-4">
									<UiSwitch v-model="formValue.announcePlay" />
									<div class="space-y-0.5 flex-1">
										<UiLabel class="text-base font-medium">{{
											t('songRequests.settings.announcePlay')
										}}</UiLabel>
									</div>
								</div>

								<div class="flex flex-row items-center gap-4 rounded-lg border p-4">
									<UiSwitch v-model="formValue.playerNoCookieMode" />
									<div class="space-y-0.5 flex-1">
										<UiLabel class="text-base font-medium">{{
											t('songRequests.settings.playerNoCookieMode')
										}}</UiLabel>
										<p class="text-sm text-muted-foreground">
											{{ t('songRequests.settings.playerNoCookieModeDescription') }}
										</p>
									</div>
								</div>
							</div>

							<div class="space-y-2">
								<UiLabel>{{ t('songRequests.settings.neededPercentageForskip') }}</UiLabel>
								<UiInput v-model="formValue.neededVotesForSkip" type="number" :min="0" :max="100" />
							</div>

							<div class="space-y-2">
								<UiLabel>{{ t('songRequests.settings.channelReward') }}</UiLabel>
								<RewardsSelector
									v-model="formValue.channelPointsRewardId"
									only-with-input
									clearable
								/>
							</div>

							<div class="space-y-2">
								<UiLabel>{{ t('songRequests.settings.deniedChannels') }}</UiLabel>
								<div v-if="formValue.denyList!.channels.length" class="flex flex-wrap gap-2 mb-2">
									<div
										v-for="channelId of formValue.denyList!.channels"
										:key="channelId"
										class="flex items-center gap-2 bg-muted rounded-md px-2 py-1"
									>
										<UiAvatar class="size-5">
											<UiAvatarImage :src="findChannelImage(channelId)" />
											<UiAvatarFallback>{{ findChannelLabel(channelId).charAt(0) }}</UiAvatarFallback>
										</UiAvatar>
										<span class="text-sm">{{ findChannelLabel(channelId) }}</span>
										<button
											type="button"
											class="text-muted-foreground hover:text-foreground cursor-pointer"
											@click="
												formValue.denyList!.channels = formValue.denyList!.channels.filter(
													(c) => c !== channelId
												)
											"
										>
											×
										</button>
									</div>
								</div>
								<UiPopover>
									<UiPopoverTrigger as-child>
										<UiButton variant="outline" class="w-full justify-start">
											{{ 'Search channels...' }}
										</UiButton>
									</UiPopoverTrigger>
									<UiPopoverContent class="w-full p-0" align="start">
										<UiCommand>
											<UiCommandInput
												v-model="channelsSearchValue"
												placeholder="Search channels..."
											/>
											<UiCommandList>
												<UiCommandEmpty>No channels found.</UiCommandEmpty>
												<UiCommandGroup>
													<UiCommandItem
														v-for="option of channelsOptions.filter(
															(o) => !formValue.denyList!.channels.includes(o.value)
														)"
														:key="option.value"
														:value="option.value"
														@select="
															() => {
																if (!formValue.denyList!.channels.includes(option.value)) {
																	formValue.denyList!.channels.push(option.value)
																}
															}
														"
														class="cursor-pointer"
													>
														<UiAvatar class="size-5 mr-2">
															<UiAvatarImage :src="option.image" />
															<UiAvatarFallback>{{ option.label.charAt(0) }}</UiAvatarFallback>
														</UiAvatar>
														<span>{{ option.label }}</span>
													</UiCommandItem>
												</UiCommandGroup>
											</UiCommandList>
										</UiCommand>
									</UiPopoverContent>
								</UiPopover>
							</div>

							<div class="space-y-2">
								<UiLabel>{{ t('songRequests.settings.deniedWords') }}</UiLabel>
								<UiTagsInput v-model="formValue.denyList!.words">
									<UiTagsInputItem
										v-for="item in formValue.denyList!.words"
										:key="item"
										:value="item"
									>
										<UiTagsInputItemText />
										<UiTagsInputItemDelete />
									</UiTagsInputItem>
									<TagsInputInput placeholder="Add word..." />
								</UiTagsInput>
							</div>
						</UiTabsContent>

						<UiTabsContent value="commands" class="mt-4">
							<CommandsList class="mb-2" :commands="srCommands" />
						</UiTabsContent>

						<UiTabsContent value="users" class="mt-4 space-y-4">
							<div class="grid grid-cols-2 gap-4">
								<div class="space-y-2">
									<UiLabel>{{ t('songRequests.settings.users.maxRequests') }}</UiLabel>
									<UiInput v-model="formValue.user!.maxRequests" type="number" :min="0" :max="1000" />
								</div>
								<div class="space-y-2">
									<UiLabel>{{ t('songRequests.settings.users.minimalWatchTime') }}</UiLabel>
									<UiInput
										v-model="formValue.user!.minWatchTime"
										type="number"
										:min="0"
										:max="999999999"
									/>
								</div>
								<div class="space-y-2">
									<UiLabel>{{ t('songRequests.settings.users.minimalMessages') }}</UiLabel>
									<UiInput
										v-model="formValue.user!.minMessages"
										type="number"
										:min="0"
										:max="999999999"
									/>
								</div>
								<div class="space-y-2">
									<UiLabel>{{ t('songRequests.settings.users.minimalFollowTime') }}</UiLabel>
									<UiInput
										v-model="formValue.user!.minFollowTime"
										type="number"
										:min="0"
										:max="99999999999999"
									/>
								</div>
							</div>

							<div class="space-y-2">
								<UiLabel>{{ t('songRequests.settings.deniedUsers') }}</UiLabel>
								<TwitchSearchUsers v-model="formValue.denyList!.users" />
							</div>
						</UiTabsContent>

						<UiTabsContent value="songs" class="mt-4 space-y-4">
							<div class="grid grid-cols-2 gap-4">
								<div class="space-y-2">
									<UiLabel>{{ t('songRequests.settings.songs.maxRequests') }}</UiLabel>
									<UiInput
										v-model="formValue.maxRequests"
										type="number"
										:min="0"
										:max="99999999999999"
									/>
								</div>
								<div class="space-y-2">
									<UiLabel>{{ t('songRequests.settings.songs.minLength') }}</UiLabel>
									<UiInput v-model="formValue.song!.minLength" type="number" :min="0" :max="999999" />
								</div>
								<div class="space-y-2">
									<UiLabel>{{ t('songRequests.settings.songs.maxLength') }}</UiLabel>
									<UiInput v-model="formValue.song!.maxLength" type="number" :min="0" :max="999999" />
								</div>
								<div class="space-y-2">
									<UiLabel>{{ t('songRequests.settings.songs.minViews') }}</UiLabel>
									<UiInput
										v-model="formValue.song!.minViews"
										type="number"
										:min="0"
										:max="99999999999999"
									/>
								</div>
							</div>

							<div class="space-y-2">
								<UiLabel>{{ t('songRequests.settings.deniedSongs') }}</UiLabel>
								<div v-if="formValue.denyList!.songs.length" class="flex flex-wrap gap-2 mb-2">
									<div
										v-for="songId of formValue.denyList!.songs"
										:key="songId"
										class="flex items-center gap-2 bg-muted rounded-md px-2 py-1"
									>
										<UiAvatar class="size-5 rounded">
											<UiAvatarImage :src="findSongImage(songId)" />
											<UiAvatarFallback>{{ findSongLabel(songId).charAt(0) }}</UiAvatarFallback>
										</UiAvatar>
										<span class="text-sm">{{ findSongLabel(songId) }}</span>
										<button
											type="button"
											class="text-muted-foreground hover:text-foreground"
											@click="
												formValue.denyList!.songs = formValue.denyList!.songs.filter(
													(s) => s !== songId
												)
											"
										>
											×
										</button>
									</div>
								</div>
								<UiPopover>
									<UiPopoverTrigger as-child>
										<UiButton variant="outline" class="w-full justify-start">
											{{ 'Search songs...' }}
										</UiButton>
									</UiPopoverTrigger>
									<UiPopoverContent class="w-full p-0" align="start">
										<UiCommand>
											<UiCommandInput v-model="songsSearchValue" placeholder="Search songs..." />
											<UiCommandList>
												<UiCommandEmpty>No songs found.</UiCommandEmpty>
												<UiCommandGroup>
													<UiCommandItem
														v-for="option of songsSearchOptions.filter(
															(o) => !formValue.denyList!.songs.includes(o.value)
														)"
														:key="option.value"
														:value="option.value"
														@select="
															() => {
																if (!formValue.denyList!.songs.includes(option.value)) {
																	formValue.denyList!.songs.push(option.value)
																}
															}
														"
													>
														<UiAvatar class="size-5 rounded mr-2">
															<UiAvatarImage :src="option.image" />
															<UiAvatarFallback>{{ option.label.charAt(0) }}</UiAvatarFallback>
														</UiAvatar>
														<span>{{ option.label }}</span>
													</UiCommandItem>
												</UiCommandGroup>
											</UiCommandList>
										</UiCommand>
									</UiPopoverContent>
								</UiPopover>
							</div>
						</UiTabsContent>

						<UiTabsContent value="translations" class="mt-4">
							<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">notEnabled</UiLabel>
									<UiTextarea v-model="formValue.translations.notEnabled" class="min-h-20" />
								</div>

								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">nowPlaying</UiLabel>
									<UiTextarea v-model="formValue.translations.nowPlaying" class="min-h-20" />
								</div>

								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">noText</UiLabel>
									<UiTextarea v-model="formValue.translations.noText" class="min-h-20" />
								</div>

								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">acceptOnlyWhenOnline</UiLabel>
									<UiTextarea
										v-model="formValue.translations.acceptOnlyWhenOnline"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">song.notFound</UiLabel>
									<UiTextarea v-model="formValue.translations.song.notFound" class="min-h-20" />
								</div>

								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">song.alreadyInQueue</UiLabel>
									<UiTextarea v-model="formValue.translations.song.alreadyInQueue" class="min-h-20" />
								</div>

								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">song.ageRestrictions</UiLabel>
									<UiTextarea
										v-model="formValue.translations.song.ageRestrictions"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">song.cannotGetInformation</UiLabel>
									<UiTextarea
										v-model="formValue.translations.song.cannotGetInformation"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">song.live</UiLabel>
									<UiTextarea v-model="formValue.translations.song.live" class="min-h-20" />
								</div>

								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">song.denied</UiLabel>
									<UiTextarea v-model="formValue.translations.song.denied" class="min-h-20" />
								</div>

								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">song.requestedMessage</UiLabel>
									<UiTextarea
										v-model="formValue.translations.song.requestedMessage"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">song.maximumOrdered</UiLabel>
									<UiTextarea v-model="formValue.translations.song.maximumOrdered" class="min-h-20" />
								</div>

								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">song.minViews</UiLabel>
									<UiTextarea v-model="formValue.translations.song.minViews" class="min-h-20" />
								</div>

								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">song.maxLength</UiLabel>
									<UiTextarea v-model="formValue.translations.song.maxLength" class="min-h-20" />
								</div>

								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">song.minLength</UiLabel>
									<UiTextarea v-model="formValue.translations.song.minLength" class="min-h-20" />
								</div>

								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">user.denied</UiLabel>
									<UiTextarea v-model="formValue.translations.user.denied" class="min-h-20" />
								</div>

								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">user.maxRequests</UiLabel>
									<UiTextarea v-model="formValue.translations.user.maxRequests" class="min-h-20" />
								</div>

								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">user.minMessages</UiLabel>
									<UiTextarea v-model="formValue.translations.user.minMessages" class="min-h-20" />
								</div>

								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">user.minWatched</UiLabel>
									<UiTextarea v-model="formValue.translations.user.minWatched" class="min-h-20" />
								</div>

								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">user.minFollow</UiLabel>
									<UiTextarea v-model="formValue.translations.user.minFollow" class="min-h-20" />
								</div>

								<div class="space-y-1">
									<UiLabel class="text-xs text-muted-foreground">channel.denied</UiLabel>
									<UiTextarea v-model="formValue.translations.channel.denied" class="min-h-20" />
								</div>
							</div>
						</UiTabsContent>
					</UiTabs>
				</div>
			</UiScrollArea>

			<UiDialogFooter class="p-6 pt-0">
				<UiButton variant="outline" @click="isOpen = false">
					{{ t('sharedButtons.close') }}
				</UiButton>
				<UiButton @click="save">
					{{ t('sharedButtons.save') }}
				</UiButton>
			</UiDialogFooter>
		</UiDialogContent>
	</UiDialog>
</template>
