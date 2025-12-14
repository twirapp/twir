<script setup lang="ts">
import { Plus, Trash2 } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import DiscordSvg from '@/assets/integrations/discord.svg?use'

import DiscordGuildSettingsForm from '../ui/discord/guild-settings-form.vue'

import { useProfile } from '@/api'
import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
	AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { toast } from 'vue-sonner'
import { useDiscordIntegration } from '@/features/integrations/composables/discord/use-discord-integration.ts'
import PageLayout from '@/layout/page-layout.vue'

const { t } = useI18n()
const route = useRoute()

const { guilds, authLink, disconnectGuild } = useDiscordIntegration()

const activeGuildId = ref<string | undefined>()

// Set initial active guild from route query or first guild
watch(
	[guilds, () => route.query.guild],
	([g, queryGuild]) => {
		if (!g || g.length === 0) {
			activeGuildId.value = undefined
			return
		}

		const guildIdFromQuery = typeof queryGuild === 'string' ? queryGuild : undefined
		const guildExists = guildIdFromQuery && g.some((guild) => guild.id === guildIdFromQuery)

		if (guildExists) {
			activeGuildId.value = guildIdFromQuery
		} else if (!activeGuildId.value || !g.some((guild) => guild.id === activeGuildId.value)) {
			activeGuildId.value = g[0]?.id
		}
	},
	{ immediate: true }
)

const isConnectDisabled = computed(() => (guilds.value?.length ?? 0) >= 2)

function getGuildIconUrl(guildId: string, icon: string | null | undefined) {
	if (!icon) return undefined
	return `https://cdn.discordapp.com/icons/${guildId}/${icon}.png`
}

function handleConnectGuild() {
	if (!authLink.value) return
	window.open(authLink.value, 'Twir connect discord', 'width=800,height=600')
}

async function handleDisconnectGuild(guildId: string) {
	const result = await disconnectGuild(guildId)
	if (result.error) {
		toast.error(t('sharedTexts.error'), {
			description: result.error.message,
			duration: 5000,
		})
	} else {
		toast.success(t('sharedTexts.saved'), {
			duration: 2500,
		})
	}
}

// Handle callback from OAuth popup
const { refetchData } = useDiscordIntegration()
if (typeof window !== 'undefined') {
	window.addEventListener('message', async (event) => {
		if (event.data === 'discord-connected') {
			await refetchData()
		}
	})
}

const { data: currentUser } = useProfile()
</script>

<template>
	<PageLayout show-back back-redirect-to="/dashboard/integrations">
		<template #title>
			<div class="flex items-center gap-3">
				<DiscordSvg class="h-8 w-8 text-white" alt="Discord" />
				Discord
			</div>
		</template>

		<template #title-footer>
			<span v-if="guilds.length > 0" class="text-sm text-muted-foreground">
				{{
					t('integrations.discord.connectedGuilds', {
						guilds: t('integrations.discord.guildPluralization', guilds.length),
					})
				}}
			</span>
		</template>

		<template #action>
			<div class="flex items-center gap-2">
				<Button
					:disabled="isConnectDisabled"
					variant="secondary"
					size="sm"
					@click="handleConnectGuild"
				>
					<Plus class="h-4 w-4 mr-2" />
					{{ t('integrations.discord.connectGuild') }}
				</Button>
			</div>
		</template>

		<template #content>
			<div v-if="guilds.length === 0" class="flex flex-col items-center justify-center py-12 gap-4">
				<p class="text-muted-foreground">{{ t('integrations.discord.noGuilds') }}</p>
				<Button variant="default" @click="handleConnectGuild">
					<Plus class="h-4 w-4 mr-2" />
					{{ t('integrations.discord.connectGuild') }}
				</Button>
			</div>

			<div v-else class="flex flex-col gap-6">
				<Tabs v-model="activeGuildId" class="w-full">
					<TabsList class="mb-4">
						<TabsTrigger v-for="guild in guilds" :key="guild.id" :value="guild.id">
							<div class="flex items-center justify-center gap-2">
								<Avatar class="h-5 w-5">
									<AvatarImage v-if="guild.icon" :src="getGuildIconUrl(guild.id, guild.icon)!" />
									<AvatarFallback class="text-xs">{{ guild.name.charAt(0) }}</AvatarFallback>
								</Avatar>
								<span>{{ guild.name }}</span>
							</div>
						</TabsTrigger>
					</TabsList>

					<TabsContent v-for="guild in guilds" :key="guild.id" :value="guild.id">
						<div class="flex flex-col gap-6">
							<DiscordGuildSettingsForm
								:guild-id="guild.id"
								:initial-values="{
									liveNotificationEnabled: guild.liveNotificationEnabled,
									liveNotificationChannelsIds: guild.liveNotificationChannelsIds,
									liveNotificationShowTitle: guild.liveNotificationShowTitle,
									liveNotificationShowCategory: guild.liveNotificationShowCategory,
									liveNotificationShowViewers: guild.liveNotificationShowViewers,
									liveNotificationMessage: guild.liveNotificationMessage,
									liveNotificationShowPreview: guild.liveNotificationShowPreview,
									liveNotificationShowProfileImage: guild.liveNotificationShowProfileImage,
									offlineNotificationMessage: guild.offlineNotificationMessage,
									shouldDeleteMessageOnOffline: guild.shouldDeleteMessageOnOffline,
									additionalUsersIdsForLiveCheck: guild.additionalUsersIdsForLiveCheck,
								}"
								:current-user="currentUser"
							/>

							<!-- Danger Zone -->
							<div class="border border-destructive/50 rounded-lg p-4">
								<h3 class="text-lg font-medium text-destructive mb-4">
									{{ t('sharedTexts.dangerZone') }}
								</h3>

								<AlertDialog>
									<AlertDialogTrigger as-child>
										<Button variant="destructive">
											<Trash2 class="h-4 w-4 mr-2" />
											{{ t('integrations.discord.disconnectGuild') }}
										</Button>
									</AlertDialogTrigger>
									<AlertDialogContent>
										<AlertDialogHeader>
											<AlertDialogTitle>{{ t('deleteConfirmation.text') }}</AlertDialogTitle>
											<AlertDialogDescription>
												{{ t('integrations.discord.disconnectGuild') }} - {{ guild.name }}
											</AlertDialogDescription>
										</AlertDialogHeader>
										<AlertDialogFooter>
											<AlertDialogCancel>{{ t('deleteConfirmation.cancel') }}</AlertDialogCancel>
											<AlertDialogAction @click="handleDisconnectGuild(guild.id)">
												{{ t('deleteConfirmation.confirm') }}
											</AlertDialogAction>
										</AlertDialogFooter>
									</AlertDialogContent>
								</AlertDialog>
							</div>
						</div>
					</TabsContent>
				</Tabs>
			</div>
		</template>
	</PageLayout>
</template>
