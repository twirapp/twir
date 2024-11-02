<script lang="ts" setup>
import { ChevronsUpDown, LogOut, Settings, Shield } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import { useLogout, useProfile } from '@/api'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuGroup,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { SidebarMenuButton } from '@/components/ui/sidebar'

const { t } = useI18n()
const { data: profileData } = useProfile()

const logout = useLogout()
</script>

<template>
	<DropdownMenu v-if="profileData">
		<DropdownMenuTrigger as-child>
			<SidebarMenuButton
				size="lg"
				class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
			>
				<Avatar class="h-8 w-8 rounded-lg">
					<AvatarImage :src="profileData.avatar" :alt="profileData.login" />
					<AvatarFallback class="rounded-lg">
						{{ profileData.login.slice(0, 2).toUpperCase() }}
					</AvatarFallback>
				</Avatar>
				<div class="grid flex-1 text-left text-sm leading-tight">
					<span class="truncate font-semibold">{{ profileData.displayName }}</span>
					<span class="truncate text-xs">Logged as</span>
				</div>
				<ChevronsUpDown class="ml-auto size-4" />
			</SidebarMenuButton>
		</DropdownMenuTrigger>
		<DropdownMenuContent class="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg" side="bottom" align="end" :side-offset="4">
			<DropdownMenuLabel class="p-0 font-normal">
				<div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
					<Avatar class="h-8 w-8 rounded-lg">
						<AvatarImage :src="profileData.avatar" :alt="profileData.login" />
						<AvatarFallback class="rounded-lg">
							{{ profileData.login.slice(0, 2).toUpperCase() }}
						</AvatarFallback>
					</Avatar>
					<div class="grid flex-1 text-left text-sm leading-tight">
						<span class="truncate font-semibold">{{ profileData.displayName }}</span>
						<span class="truncate text-xs">Logged as</span>
					</div>
				</div>
			</DropdownMenuLabel>
			<DropdownMenuSeparator />
			<DropdownMenuGroup>
				<DropdownMenuItem as-child>
					<RouterLink to="/dashboard/settings" class="flex items-center">
						<Settings class="mr-2 size-4" />
						Settings
					</RouterLink>
				</DropdownMenuItem>
			</DropdownMenuGroup>
			<DropdownMenuSeparator />
			<template v-if="profileData.isBotAdmin">
				<DropdownMenuGroup>
					<DropdownMenuItem>
						<RouterLink to="/dashboard/admin" class="flex items-center">
							<Shield class="mr-2 size-4" />
							Admin panel
						</RouterLink>
					</DropdownMenuItem>
				</DropdownMenuGroup>
				<DropdownMenuSeparator />
			</template>
			<DropdownMenuItem @click="() => logout">
				<LogOut class="mr-2 size-4" />
				{{ t('navbar.logout') }}
			</DropdownMenuItem>
		</DropdownMenuContent>
	</DropdownMenu>
</template>
