<script lang="ts" setup>
import { useLocalStorage } from '@vueuse/core'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { useLogout, useProfile } from '~~/layers/dashboard/api/auth'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
	DropdownMenu,
	DropdownMenuCheckboxItem,
	DropdownMenuContent,
	DropdownMenuGroup,
	DropdownMenuItem,
	DropdownMenuPortal,
	DropdownMenuSeparator,
	DropdownMenuSub,
	DropdownMenuSubContent,
	DropdownMenuSubTrigger,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

import { useTheme } from '~~/layers/dashboard/composables/use-theme.js'
import { AVAILABLE_LOCALES } from '~~/layers/dashboard/config/i18n-locales.js'
import { Button } from '@/components/ui/button'
import KickIcon from '~~/layers/dashboard/components/kick-icon.vue'
import TwitchIcon from '~~/layers/dashboard/components/twitch-icon.vue'
import type { DropdownMenuContentProps } from 'reka-ui'

const { t, locale } = useI18n()
const { data: profileData } = useProfile()

const theme = useTheme()
const logout = useLogout()
const currentLocale = useLocalStorage<string>('twirLocale', 'en')

const dropdownProps = computed((): DropdownMenuContentProps & { class?: string } => {
	return {
		class: 'w-(--radix-dropdown-menu-trigger-width)',
		side: 'bottom',
		align: 'end',
		sideOffset: 4,
	}
})

const linkedAccounts = computed(() => profileData.value?.linkedAccounts ?? [])
</script>

<template>
	<DropdownMenu v-if="profileData">
		<DropdownMenuTrigger as-child>
			<Button
				class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground flex gap-2"
				variant="secondary"
			>
				<Avatar class="size-6 rounded-full">
					<AvatarImage :src="profileData.avatar" :alt="profileData.login" />
					<AvatarFallback class="rounded-lg">
						{{ profileData.login.slice(0, 2).toUpperCase() }}
					</AvatarFallback>
				</Avatar>
				<span class="truncate font-semibold">{{ profileData.displayName }}</span>
				<div class="flex items-center gap-1">
					<template v-for="account in linkedAccounts" :key="account.platform">
						<KickIcon
							v-if="account.platform === 'kick'"
							class="size-4 text-[#53FC18]"
						/>
						<TwitchIcon
							v-else-if="account.platform === 'twitch'"
							class="size-4 text-[#9146FF]"
						/>
					</template>
				</div>
				<Icon name="lucide:chevrons-up-down" class="ml-auto size-4" />
			</Button>
		</DropdownMenuTrigger>

		<DropdownMenuContent class="min-w-56 rounded-lg" v-bind="dropdownProps">
			<DropdownMenuGroup>
				<DropdownMenuItem as-child>
					<RouterLink to="/dashboard/settings" class="flex items-center">
						<Icon name="lucide:settings" class="mr-2 size-4" />
						{{ t('sharedButtons.settings') }}
					</RouterLink>
				</DropdownMenuItem>

				<DropdownMenuSub>
					<DropdownMenuSubTrigger>
						<Icon name="lucide:languages" class="mr-2 h-4 w-4" />
						{{ t('sidebar.lang') }}
					</DropdownMenuSubTrigger>
					<DropdownMenuPortal>
						<DropdownMenuSubContent>
							<DropdownMenuItem
								v-for="lang of AVAILABLE_LOCALES"
								:key="lang.code"
								@select="
									() => {
										locale = lang.code
										currentLocale = lang.code
									}
								"
							>
								<DropdownMenuCheckboxItem :checked="currentLocale === lang.code" />
								{{ lang.name }}
							</DropdownMenuItem>
						</DropdownMenuSubContent>
					</DropdownMenuPortal>
				</DropdownMenuSub>

				<DropdownMenuItem as-child @select.prevent="theme.toggleTheme">
					<div>
						<template v-if="theme.isDark.value">
							<Icon name="lucide:sun" class="mr-2 size-4" />
							{{ t('sidebar.lightTheme') }}
						</template>
						<template v-else>
							<Icon name="lucide:moon" class="mr-2 size-4" />
							{{ t('sidebar.darkTheme') }}
						</template>
					</div>
				</DropdownMenuItem>

				<DropdownMenuItem v-if="profileData.isBotAdmin" as-child>
					<RouterLink to="/dashboard/admin" class="flex items-center">
						<Icon name="lucide:shield" class="mr-2 size-4" />
						{{ t('adminPanel.title') }}
					</RouterLink>
				</DropdownMenuItem>
			</DropdownMenuGroup>

			<DropdownMenuSeparator />

			<DropdownMenuItem @click="logout" class="text-red-500">
				<Icon name="lucide:log-out" class="mr-2 size-4" />
				{{ t('sidebar.logout') }}
			</DropdownMenuItem>
		</DropdownMenuContent>
	</DropdownMenu>
</template>
