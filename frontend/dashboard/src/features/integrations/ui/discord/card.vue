<script setup lang="ts">
import { ExternalLink, Settings } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { useUserAccessFlagChecker } from '@/api'
import DiscordIcon from '@/assets/integrations/discord.svg?use'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { useDiscordIntegration } from '@/features/integrations/composables/discord/use-discord-integration.js'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

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
	<Card class="flex flex-col h-full">
		<CardHeader>
			<CardTitle class="flex items-center gap-2">
				<DiscordIcon class="w-8 h-8" style="fill: #5865f2" />
				Discord
			</CardTitle>
		</CardHeader>

		<CardContent class="grow">
			<p class="text-sm text-muted-foreground">
				{{ t('integrations.discord.description') }}
			</p>
		</CardContent>

		<CardFooter class="mt-auto">
			<div class="flex justify-between flex-wrap items-center gap-4 w-full">
				<RouterLink custom v-slot="{ href, navigate }" :to="{ name: 'DiscordIntegration' }">
					<Button
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
					</Button>
				</RouterLink>

				<div class="flex items-center gap-2">
					<template v-if="isLoading">
						<span class="text-sm text-muted-foreground">{{ t('sharedTexts.loading') }}</span>
					</template>
					<template v-else-if="connectedGuildsCount > 0">
						<div class="flex -space-x-2">
							<Avatar
								v-for="guild in guilds.slice(0, 3)"
								:key="guild.id"
								class="h-6 w-6 border-2 border-background"
							>
								<AvatarImage v-if="guild.icon" :src="getGuildIconUrl(guild.id, guild.icon)!" />
								<AvatarFallback class="text-xs">{{ guild.name.charAt(0) }}</AvatarFallback>
							</Avatar>
						</div>
						<Badge variant="secondary" class="ml-2">
							{{
								t('integrations.discord.connectedGuilds', {
									guilds: t('integrations.discord.guildPluralization', connectedGuildsCount),
								})
							}}
						</Badge>
					</template>
					<template v-else>
						<Badge variant="outline" class="text-muted-foreground">
							{{ t('integrations.discord.noGuilds') }}
						</Badge>
					</template>
				</div>
			</div>
		</CardFooter>
	</Card>
</template>
