<script setup lang="ts">
import { LinkIcon, UnlinkIcon } from 'lucide-vue-next'
import { computed } from 'vue'

import { useAuthLink, useProfile, useUnlinkPlatformAccount } from '@/api/auth'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'

const userSettingsPath = '/dashboard/settings'
const { data: profile, executeQuery } = useProfile()
const unlinkAccount = useUnlinkPlatformAccount()
const { data: twitchAuthLinkData, fetching: twitchAuthLinkFetching } = useAuthLink(userSettingsPath)

const accounts = computed(() => profile.value?.linkedAccounts || [])
const currentPlatform = computed(() => profile.value?.currentPlatform || '')
const twitchAuthLink = computed(() => twitchAuthLinkData.value?.authLink ?? null)

const isKickLinked = computed(() => accounts.value.some((a) => a.platform === 'kick'))
const isTwitchLinked = computed(() => accounts.value.some((a) => a.platform === 'twitch'))

async function handleUnlink(platform: string) {
	if (platform === currentPlatform.value) return
	await unlinkAccount.executeMutation({ platform })
	await executeQuery({ requestPolicy: 'network-only' })
}

function handleConnectKick() {
	window.location.href = `/api/auth/kick/authorize?redirect_to=${userSettingsPath}`
}

function handleConnectTwitch() {
	if (!twitchAuthLink.value) return
	window.location.href = twitchAuthLink.value
}
</script>

<template>
	<div class="flex flex-col gap-6">
		<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">Linked Accounts</h4>

		<div class="flex flex-col gap-4">
			<Card
				v-for="account in accounts"
				:key="account.platform"
			>
				<CardContent class="flex items-center justify-between p-4">
					<div class="flex items-center gap-4">
						<Avatar>
							<AvatarImage
								:src="account.platformAvatar || ''"
								:alt="account.platformLogin"
							/>
							<AvatarFallback>{{ account.platformLogin.slice(0, 2).toUpperCase() }}</AvatarFallback>
						</Avatar>
						<div class="flex flex-col">
							<div class="flex items-center gap-2">
								<span class="font-semibold">{{ account.platformLogin }}</span>
								<Badge
									variant="secondary"
									class="uppercase"
									>{{ account.platform }}</Badge
								>
								<Badge
									v-if="account.platform === currentPlatform"
									variant="default"
									>Primary</Badge
								>
							</div>
						</div>
					</div>
					<Button
						v-if="account.platform !== currentPlatform"
						variant="destructive"
						size="sm"
						@click="handleUnlink(account.platform)"
					>
						<UnlinkIcon class="mr-2 h-4 w-4" />
						Disconnect
					</Button>
				</CardContent>
			</Card>

			<Card v-if="!isTwitchLinked">
				<CardContent class="flex items-center justify-between p-4">
					<div class="flex items-center gap-4">
						<Avatar>
							<AvatarFallback>T</AvatarFallback>
						</Avatar>
						<div class="flex flex-col">
							<span class="font-semibold">Twitch</span>
							<span class="text-muted-foreground text-sm">Not connected</span>
						</div>
					</div>
					<Button
						variant="default"
						size="sm"
						:disabled="twitchAuthLinkFetching || !twitchAuthLink"
						@click="handleConnectTwitch"
					>
						<LinkIcon class="mr-2 h-4 w-4" />
						Connect Twitch
					</Button>
				</CardContent>
			</Card>

			<Card v-if="!isKickLinked">
				<CardContent class="flex items-center justify-between p-4">
					<div class="flex items-center gap-4">
						<Avatar>
							<AvatarFallback>K</AvatarFallback>
						</Avatar>
						<div class="flex flex-col">
							<span class="font-semibold">Kick</span>
							<span class="text-muted-foreground text-sm">Not connected</span>
						</div>
					</div>
					<Button
						variant="default"
						size="sm"
						@click="handleConnectKick"
					>
						<LinkIcon class="mr-2 h-4 w-4" />
						Connect Kick
					</Button>
				</CardContent>
			</Card>
		</div>
	</div>
</template>
