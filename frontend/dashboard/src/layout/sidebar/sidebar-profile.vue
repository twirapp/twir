<script lang="ts" setup>
import { useLocalStorage } from '@vueuse/core'
import { ChevronsUpDown, Languages, LogOut, MoonIcon, Settings, Shield, SunIcon } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import type { DropdownMenuContentProps } from 'radix-vue'

import { useLogout, useProfile } from '@/api'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuPortal,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import DropdownMenuCheckboxItem from '@/components/ui/dropdown-menu/DropdownMenuCheckboxItem.vue'
import DropdownMenuGroup from '@/components/ui/dropdown-menu/DropdownMenuGroup.vue'
import DropdownMenuSub from '@/components/ui/dropdown-menu/DropdownMenuSub.vue'
import DropdownMenuSubContent from '@/components/ui/dropdown-menu/DropdownMenuSubContent.vue'
import DropdownMenuSubTrigger from '@/components/ui/dropdown-menu/DropdownMenuSubTrigger.vue'
import { SidebarMenuButton, useSidebar } from '@/components/ui/sidebar'
import { useTheme } from '@/composables/use-theme'
import { AVAILABLE_LOCALES } from '@/plugins/i18n'

const { t, locale } = useI18n()
const { data: profileData } = useProfile()

const { open: sidebarOpen, isMobile, setOpenMobile } = useSidebar()
const theme = useTheme()
const logout = useLogout()
const currentLocale = useLocalStorage<string>('twirLocale', 'en')

function toggleMobileSidebar() {
	if (isMobile.value) {
		setOpenMobile(false)
	}
}

const dropdownProps = computed((): DropdownMenuContentProps & { class?: string } => {
	if (sidebarOpen.value) return {
		class: 'w-[--radix-dropdown-menu-trigger-width]',
		side: 'bottom',
		align: 'end',
		sideOffset: 4,
	}

	return {
		class: 'w-[300px]',
		alignOffset: -4,
		align: 'start',
		sideOffset: 12,
		side: 'right',
	}
})
</script>

<template>
	<DropdownMenu v-if="profileData">
		<DropdownMenuTrigger as-child>
			<SidebarMenuButton
				size="lg"
				class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
			>
				<Avatar class="size-8 rounded-lg">
					<AvatarImage :src="profileData.avatar" :alt="profileData.login" />
					<AvatarFallback class="rounded-lg">
						{{ profileData.login.slice(0, 2).toUpperCase() }}
					</AvatarFallback>
				</Avatar>
				<div class="grid flex-1 text-left text-sm leading-tight">
					<span class="truncate font-semibold">{{ profileData.displayName }}</span>
					<span class="truncate text-xs">{{ t('sidebar.loggedAs') }}</span>
				</div>
				<ChevronsUpDown class="ml-auto size-4" />
			</SidebarMenuButton>
		</DropdownMenuTrigger>

		<DropdownMenuContent
			class="min-w-56 rounded-lg"
			v-bind="dropdownProps"
		>
			<DropdownMenuGroup>
				<DropdownMenuItem as-child @click="toggleMobileSidebar">
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
								@select="() => {
									locale = lang.code
									currentLocale = lang.code
									toggleMobileSidebar()
								}"
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

				<DropdownMenuItem
					v-if="profileData.isBotAdmin"
					as-child
					@click="toggleMobileSidebar"
				>
					<RouterLink to="/dashboard/admin" class="flex items-center">
						<Shield class="mr-2 size-4" />
						{{ t('adminPanel.title') }}
					</RouterLink>
				</DropdownMenuItem>
			</DropdownMenuGroup>

			<DropdownMenuSeparator />

			<DropdownMenuItem @click="logout">
				<LogOut class="mr-2 size-4" />
				{{ t('sidebar.logout') }}
			</DropdownMenuItem>
		</DropdownMenuContent>
	</DropdownMenu>
</template>
