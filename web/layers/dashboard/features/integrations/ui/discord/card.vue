<script setup lang="ts">
import { ExternalLink, Settings } from 'lucide-vue-next'
import { computed } from 'vue'


import { useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import DiscordIcon from '~/assets/integrations/discord.svg?use'




import { useDiscordIntegration } from '~/features/integrations/composables/discord/use-discord-integration.js'
import { ChannelRolePermissionEnum } from '~/gql/graphql'

const { t } = useI18n()

const { guilds, isLoading } = useDiscordIntegration()

const userCanManageIntegrations = useUserAccessFlagChecker(
	ChannelRolePermissionEnum.ManageIntegrations
)

const connectedGuildsCount = computed(() => guilds.value?.length ?? 0)

function getGuildIconUrl(guildId: string, icon: string | null | undefined) {
	if (!icon) return undefined
	return `https://cdn.discordapp.com/icons/${guildId}/${icon}.png`
}
</script>

<template>
	<UiCard class="flex flex-col h-full">
		<UiCardHeader>
			<UiCardTitle class="flex items-center gap-2">
				<DiscordIcon class="w-8 h-8" style="fill: #5865f2" />
				Discord
			</UiCardTitle>
		</UiCardHeader>

		<UiCardContent class="grow">
			<p class="text-sm text-muted-foreground">
				{{ t('integrations.discord.description') }}
			</p>
		</UiCardContent>

		<UiCardFooter class="mt-auto">
			<div class="flex justify-between flex-wrap items-center gap-4 w-full">
				<RouterLink custom v-slot="{ href, navigate }" :to="{ name: 'DiscordIntegration' }">
					<UiButton
						:disabled="!userCanManageIntegrations"
						variant="secondary"
						size="sm"
						@click="navigate"
						as="a"
						:href="href"
					>
						<Settings class="mr-2 h-4 w-4" />
						{{ t('sharedButtons.settings') }}
						<ExternalLink class="ml-2 h-3 w-3" />
					</UiButton>
				</RouterLink>

				<div class="flex items-center gap-2">
					<template v-if="isLoading">
						<span class="text-sm text-muted-foreground">{{ t('sharedTexts.loading') }}</span>
					</template>
					<template v-else-if="connectedGuildsCount > 0">
						<div class="flex -space-x-2">
							<UiAvatar
								v-for="guild in guilds.slice(0, 3)"
								:key="guild.id"
								class="h-6 w-6 border-2 border-background"
							>
								<UiAvatarImage v-if="guild.icon" :src="getGuildIconUrl(guild.id, guild.icon)!" />
								<UiAvatarFallback class="text-xs">{{ guild.name.charAt(0) }}</UiAvatarFallback>
							</UiAvatar>
						</div>
						<UiBadge variant="secondary" class="ml-2">
							{{
								t('integrations.discord.connectedGuilds', {
									guilds: t('integrations.discord.guildPluralization', connectedGuildsCount),
								})
							}}
						</UiBadge>
					</template>
					<template v-else>
						<UiBadge variant="outline" class="text-muted-foreground">
							{{ t('integrations.discord.noGuilds') }}
						</UiBadge>
					</template>
				</div>
			</div>
		</UiCardFooter>
	</UiCard>
</template>
