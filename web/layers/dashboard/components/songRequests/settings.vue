<script setup lang="ts">
import { useDebounce } from '@vueuse/core'
import { computed, ref, toRaw, unref, watch } from 'vue'
import { toast } from 'vue-sonner'
import { useCommandsApi } from '~~/layers/dashboard/api/commands/commands'
import { useSongRequestsApi } from '~~/layers/dashboard/api/song-requests'
import TwitchSearchUsers from '~~/layers/dashboard/components/twitchUsers/twitch-users-select.vue'
import CommandsList from '~~/layers/dashboard/features/commands/ui/list.vue'

import type { SongRequestsSettingsOpts } from '~/gql/graphql.js'

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import {
	Command,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
} from '@/components/ui/command'
import {
	Dialog,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Switch } from '@/components/ui/switch'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
	TagsInput,
	TagsInputInput,
	TagsInputItem,
	TagsInputItemDelete,
	TagsInputItemText,
} from '@/components/ui/tags-input'
import { Textarea } from '@/components/ui/textarea'
import { SongRequestsSearchChannelOrVideoOptsType } from '~/gql/graphql.js'

import RewardsSelector from '../rewardsSelector.vue'

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
	hideOnPause: true,
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
	const { channelApiKey, __typename, ...data } = unref(formValue) as Record<string, unknown>
	await youtubeModuleUpdater.executeMutation({ opts: data as SongRequestsSettingsOpts })
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
	<Dialog v-model:open="isOpen">
		<DialogContent class="max-h-[85vh] p-0 sm:max-w-[700px]">
			<DialogHeader class="p-6 pb-0">
				<DialogTitle>{{ t('sharedTexts.settings') }}</DialogTitle>
			</DialogHeader>

			<ScrollArea class="max-h-[calc(85vh-140px)]">
				<div class="p-6 pt-4">
					<Tabs
						default-value="general"
						class="w-full"
					>
						<TabsList class="grid w-full grid-cols-5">
							<TabsTrigger value="general">{{ t('songRequests.tabs.general') }}</TabsTrigger>
							<TabsTrigger value="commands">{{ t('commands.name') }}</TabsTrigger>
							<TabsTrigger value="users">{{ t('songRequests.tabs.users') }}</TabsTrigger>
							<TabsTrigger value="songs">{{ t('songRequests.tabs.songs') }}</TabsTrigger>
							<TabsTrigger value="translations">{{
								t('songRequests.tabs.translations')
							}}</TabsTrigger>
						</TabsList>

						<TabsContent
							value="general"
							class="mt-4 space-y-4"
						>
							<div class="space-y-3">
								<div class="flex flex-row items-center gap-4 rounded-lg border p-4">
									<Switch v-model="formValue.enabled" />
									<div class="flex-1 space-y-0.5">
										<Label class="text-base font-medium">{{ t('sharedTexts.enabled') }}</Label>
									</div>
								</div>

								<div class="flex flex-row items-center gap-4 rounded-lg border p-4">
									<Switch v-model="formValue.takeSongFromDonationMessages" />
									<div class="flex-1 space-y-0.5">
										<Label class="text-base font-medium">{{
											t('songRequests.settings.takeSongFromDonationMessage')
										}}</Label>
										<p class="text-muted-foreground text-sm">
											{{ t('songRequests.settings.takeSongFromDonationMessageDescription') }}
										</p>
									</div>
								</div>

								<div class="flex flex-row items-center gap-4 rounded-lg border p-4">
									<Switch v-model="formValue.acceptOnlyWhenOnline" />
									<div class="flex-1 space-y-0.5">
										<Label class="text-base font-medium">{{
											t('songRequests.settings.onlineOnly')
										}}</Label>
										<p class="text-muted-foreground text-sm">
											{{ t('songRequests.settings.onlineOnlyDescription') }}
										</p>
									</div>
								</div>

								<div class="flex flex-row items-center gap-4 rounded-lg border p-4">
									<Switch v-model="formValue.announcePlay" />
									<div class="flex-1 space-y-0.5">
										<Label class="text-base font-medium">{{
											t('songRequests.settings.announcePlay')
										}}</Label>
									</div>
								</div>

								<div class="flex flex-row items-center gap-4 rounded-lg border p-4">
									<Switch v-model="formValue.playerNoCookieMode" />
									<div class="flex-1 space-y-0.5">
										<Label class="text-base font-medium">{{
											t('songRequests.settings.playerNoCookieMode')
										}}</Label>
										<p class="text-muted-foreground text-sm">
											{{ t('songRequests.settings.playerNoCookieModeDescription') }}
										</p>
									</div>
								</div>

								<div class="flex flex-row items-center gap-4 rounded-lg border p-4">
									<Switch v-model="formValue.hideOnPause" />
									<div class="flex-1 space-y-0.5">
										<Label class="text-base font-medium">{{
											t('songRequests.settings.hideOnPause')
										}}</Label>
										<p class="text-muted-foreground text-sm">
											{{ t('songRequests.settings.hideOnPauseDescription') }}
										</p>
									</div>
								</div>
							</div>

							<div class="space-y-2">
								<Label>{{ t('songRequests.settings.neededPercentageForskip') }}</Label>
								<Input
									v-model="formValue.neededVotesForSkip"
									type="number"
									:min="0"
									:max="100"
								/>
							</div>

							<div class="space-y-2">
								<Label>{{ t('songRequests.settings.channelReward') }}</Label>
								<RewardsSelector
									v-model="formValue.channelPointsRewardId"
									only-with-input
									clearable
								/>
							</div>

							<div class="space-y-2">
								<Label>{{ t('songRequests.settings.deniedChannels') }}</Label>
								<div
									v-if="formValue.denyList!.channels.length"
									class="mb-2 flex flex-wrap gap-2"
								>
									<div
										v-for="channelId of formValue.denyList!.channels"
										:key="channelId"
										class="bg-muted flex items-center gap-2 rounded-md px-2 py-1"
									>
										<Avatar class="size-5">
											<AvatarImage :src="findChannelImage(channelId)" />
											<AvatarFallback>{{ findChannelLabel(channelId).charAt(0) }}</AvatarFallback>
										</Avatar>
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
								<Popover>
									<PopoverTrigger as-child>
										<Button
											variant="outline"
											class="w-full justify-start"
										>
											{{ 'Search channels...' }}
										</Button>
									</PopoverTrigger>
									<PopoverContent
										class="w-full p-0"
										align="start"
									>
										<Command>
											<CommandInput
												v-model="channelsSearchValue"
												placeholder="Search channels..."
											/>
											<CommandList>
												<CommandEmpty>No channels found.</CommandEmpty>
												<CommandGroup>
													<CommandItem
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
														<Avatar class="mr-2 size-5">
															<AvatarImage :src="option.image" />
															<AvatarFallback>{{ option.label.charAt(0) }}</AvatarFallback>
														</Avatar>
														<span>{{ option.label }}</span>
													</CommandItem>
												</CommandGroup>
											</CommandList>
										</Command>
									</PopoverContent>
								</Popover>
							</div>

							<div class="space-y-2">
								<Label>{{ t('songRequests.settings.deniedWords') }}</Label>
								<TagsInput v-model="formValue.denyList!.words">
									<TagsInputItem
										v-for="item in formValue.denyList!.words"
										:key="item"
										:value="item"
									>
										<TagsInputItemText />
										<TagsInputItemDelete />
									</TagsInputItem>
									<TagsInputInput placeholder="Add word..." />
								</TagsInput>
							</div>
						</TabsContent>

						<TabsContent
							value="commands"
							class="mt-4"
						>
							<CommandsList
								class="mb-2"
								:commands="srCommands"
							/>
						</TabsContent>

						<TabsContent
							value="users"
							class="mt-4 space-y-4"
						>
							<div class="grid grid-cols-2 gap-4">
								<div class="space-y-2">
									<Label>{{ t('songRequests.settings.users.maxRequests') }}</Label>
									<Input
										v-model="formValue.user!.maxRequests"
										type="number"
										:min="0"
										:max="1000"
									/>
								</div>
								<div class="space-y-2">
									<Label>{{ t('songRequests.settings.users.minimalWatchTime') }}</Label>
									<Input
										v-model="formValue.user!.minWatchTime"
										type="number"
										:min="0"
										:max="999999999"
									/>
								</div>
								<div class="space-y-2">
									<Label>{{ t('songRequests.settings.users.minimalMessages') }}</Label>
									<Input
										v-model="formValue.user!.minMessages"
										type="number"
										:min="0"
										:max="999999999"
									/>
								</div>
								<div class="space-y-2">
									<Label>{{ t('songRequests.settings.users.minimalFollowTime') }}</Label>
									<Input
										v-model="formValue.user!.minFollowTime"
										type="number"
										:min="0"
										:max="99999999999999"
									/>
								</div>
							</div>

							<div class="space-y-2">
								<Label>{{ t('songRequests.settings.deniedUsers') }}</Label>
								<TwitchSearchUsers v-model="formValue.denyList!.users" />
							</div>
						</TabsContent>

						<TabsContent
							value="songs"
							class="mt-4 space-y-4"
						>
							<div class="grid grid-cols-2 gap-4">
								<div class="space-y-2">
									<Label>{{ t('songRequests.settings.songs.maxRequests') }}</Label>
									<Input
										v-model="formValue.maxRequests"
										type="number"
										:min="0"
										:max="99999999999999"
									/>
								</div>
								<div class="space-y-2">
									<Label>{{ t('songRequests.settings.songs.minLength') }}</Label>
									<Input
										v-model="formValue.song!.minLength"
										type="number"
										:min="0"
										:max="999999"
									/>
								</div>
								<div class="space-y-2">
									<Label>{{ t('songRequests.settings.songs.maxLength') }}</Label>
									<Input
										v-model="formValue.song!.maxLength"
										type="number"
										:min="0"
										:max="999999"
									/>
								</div>
								<div class="space-y-2">
									<Label>{{ t('songRequests.settings.songs.minViews') }}</Label>
									<Input
										v-model="formValue.song!.minViews"
										type="number"
										:min="0"
										:max="99999999999999"
									/>
								</div>
							</div>

							<div class="space-y-2">
								<Label>{{ t('songRequests.settings.deniedSongs') }}</Label>
								<div
									v-if="formValue.denyList!.songs.length"
									class="mb-2 flex flex-wrap gap-2"
								>
									<div
										v-for="songId of formValue.denyList!.songs"
										:key="songId"
										class="bg-muted flex items-center gap-2 rounded-md px-2 py-1"
									>
										<Avatar class="size-5 rounded">
											<AvatarImage :src="findSongImage(songId)" />
											<AvatarFallback>{{ findSongLabel(songId).charAt(0) }}</AvatarFallback>
										</Avatar>
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
								<Popover>
									<PopoverTrigger as-child>
										<Button
											variant="outline"
											class="w-full justify-start"
										>
											{{ 'Search songs...' }}
										</Button>
									</PopoverTrigger>
									<PopoverContent
										class="w-full p-0"
										align="start"
									>
										<Command>
											<CommandInput
												v-model="songsSearchValue"
												placeholder="Search songs..."
											/>
											<CommandList>
												<CommandEmpty>No songs found.</CommandEmpty>
												<CommandGroup>
													<CommandItem
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
														<Avatar class="mr-2 size-5 rounded">
															<AvatarImage :src="option.image" />
															<AvatarFallback>{{ option.label.charAt(0) }}</AvatarFallback>
														</Avatar>
														<span>{{ option.label }}</span>
													</CommandItem>
												</CommandGroup>
											</CommandList>
										</Command>
									</PopoverContent>
								</Popover>
							</div>
						</TabsContent>

						<TabsContent
							value="translations"
							class="mt-4"
						>
							<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">notEnabled</Label>
									<Textarea
										v-model="formValue.translations.notEnabled"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">nowPlaying</Label>
									<Textarea
										v-model="formValue.translations.nowPlaying"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">noText</Label>
									<Textarea
										v-model="formValue.translations.noText"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">acceptOnlyWhenOnline</Label>
									<Textarea
										v-model="formValue.translations.acceptOnlyWhenOnline"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">song.notFound</Label>
									<Textarea
										v-model="formValue.translations.song.notFound"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">song.alreadyInQueue</Label>
									<Textarea
										v-model="formValue.translations.song.alreadyInQueue"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">song.ageRestrictions</Label>
									<Textarea
										v-model="formValue.translations.song.ageRestrictions"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">song.cannotGetInformation</Label>
									<Textarea
										v-model="formValue.translations.song.cannotGetInformation"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">song.live</Label>
									<Textarea
										v-model="formValue.translations.song.live"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">song.denied</Label>
									<Textarea
										v-model="formValue.translations.song.denied"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">song.requestedMessage</Label>
									<Textarea
										v-model="formValue.translations.song.requestedMessage"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">song.maximumOrdered</Label>
									<Textarea
										v-model="formValue.translations.song.maximumOrdered"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">song.minViews</Label>
									<Textarea
										v-model="formValue.translations.song.minViews"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">song.maxLength</Label>
									<Textarea
										v-model="formValue.translations.song.maxLength"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">song.minLength</Label>
									<Textarea
										v-model="formValue.translations.song.minLength"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">user.denied</Label>
									<Textarea
										v-model="formValue.translations.user.denied"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">user.maxRequests</Label>
									<Textarea
										v-model="formValue.translations.user.maxRequests"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">user.minMessages</Label>
									<Textarea
										v-model="formValue.translations.user.minMessages"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">user.minWatched</Label>
									<Textarea
										v-model="formValue.translations.user.minWatched"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">user.minFollow</Label>
									<Textarea
										v-model="formValue.translations.user.minFollow"
										class="min-h-20"
									/>
								</div>

								<div class="space-y-1">
									<Label class="text-muted-foreground text-xs">channel.denied</Label>
									<Textarea
										v-model="formValue.translations.channel.denied"
										class="min-h-20"
									/>
								</div>
							</div>
						</TabsContent>
					</Tabs>
				</div>
			</ScrollArea>

			<DialogFooter class="p-6 pt-0">
				<Button
					variant="outline"
					@click="isOpen = false"
				>
					{{ t('sharedButtons.close') }}
				</Button>
				<Button @click="save">
					{{ t('sharedButtons.save') }}
				</Button>
			</DialogFooter>
		</DialogContent>
	</Dialog>
</template>
