<script lang="ts" setup>
import { useLocalStorage } from '@vueuse/core'
import {
	ChevronsUpDown,
	Languages,
	LogOut,
	MoonIcon,
	Settings,
	Shield,
	SunIcon,
} from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import type { DropdownMenuContentProps } from 'radix-vue'

import { useLogout, useProfile } from '@/api'
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

import { useTheme } from '@/composables/use-theme.ts'
import { AVAILABLE_LOCALES } from '@/plugins/i18n.ts'
import { Button } from '@/components/ui/button'

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
				<ChevronsUpDown class="ml-auto size-4" />
			</Button>
		</DropdownMenuTrigger>

		<DropdownMenuContent class="min-w-56 rounded-lg" v-bind="dropdownProps">
			<DropdownMenuGroup>
				<DropdownMenuItem as-child>
					<RouterLink to="/dashboard/settings" class="flex items-center">
						<Settings class="mr-2 size-4" />
						{{ t('sharedButtons.settings') }}
					</RouterLink>
				</DropdownMenuItem>

				<DropdownMenuSub>
					<DropdownMenuSubTrigger>
						<Languages class="mr-2 h-4 w-4" />
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
							<SunIcon class="mr-2 size-4" />
							{{ t('sidebar.lightTheme') }}
						</template>
						<template v-else>
							<MoonIcon class="mr-2 size-4" />
							{{ t('sidebar.darkTheme') }}
						</template>
					</div>
				</DropdownMenuItem>

				<DropdownMenuItem v-if="profileData.isBotAdmin" as-child>
					<RouterLink to="/dashboard/admin" class="flex items-center">
						<Shield class="mr-2 size-4" />
						{{ t('adminPanel.title') }}
					</RouterLink>
				</DropdownMenuItem>
			</DropdownMenuGroup>

			<DropdownMenuSeparator />

			<DropdownMenuItem @click="logout" class="text-red-500">
				<LogOut class="mr-2 size-4" />
				{{ t('sidebar.logout') }}
			</DropdownMenuItem>
		</DropdownMenuContent>
	</DropdownMenu>
</template>
