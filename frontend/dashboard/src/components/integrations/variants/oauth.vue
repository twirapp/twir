<script setup lang="ts">
import { LogIn, LogOut, Settings } from 'lucide-vue-next'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import type { FunctionalComponent } from 'vue'

import { useUserAccessFlagChecker } from '@/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import {
	Dialog,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from '@/components/ui/dialog'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const props = withDefaults(
	defineProps<{
		title: string
		isLoading?: boolean
		data: { userName?: string | null; avatar?: string | null } | undefined | null
		logout: () => any
		authLink?: string | null
		icon: FunctionalComponent<any>
		iconWidth?: string
		iconColor?: string
		withSettings?: boolean
		save?: () => any | Promise<any>
	}>(),
	{
		authLink: '',
		description: '',
	}
)

defineSlots<{
	settings?: FunctionalComponent
	description?: FunctionalComponent | string
}>()

const showSettings = ref(false)

async function login() {
	if (!props.authLink) return

	window.open(props.authLink, 'Twir connect integration', 'width=800,height=600')
}

async function saveSettings() {
	await props.save?.()
	showSettings.value = false
}

const userCanManageIntegrations = useUserAccessFlagChecker(
	ChannelRolePermissionEnum.ManageIntegrations
)

const { t } = useI18n()
</script>

<template>
	<Card class="flex flex-col h-full">
		<CardHeader>
			<CardTitle class="flex items-center gap-2">
				<component :is="icon" :style="{ width: iconWidth }" class="w-8 h-8" />
				{{ title }}
			</CardTitle>
		</CardHeader>

		<CardContent class="grow">
			<slot name="description" />
		</CardContent>

		<CardFooter class="mt-auto">
			<div class="flex justify-between flex-wrap items-center gap-4 w-full">
				<div class="flex gap-2 flex-wrap">
					<Button
						v-if="withSettings"
						:disabled="!userCanManageIntegrations"
						variant="secondary"
						size="sm"
						@click="showSettings = true"
					>
						<Settings class="mr-2 h-4 w-4" />
						{{ t('sharedButtons.settings') }}
					</Button>

					<Button
						:disabled="!userCanManageIntegrations || !authLink"
						:variant="data?.userName ? 'destructive' : 'default'"
						size="sm"
						@click="data?.userName ? logout() : login()"
					>
						<LogOut v-if="data?.userName" class="mr-2 h-4 w-4" />
						<LogIn v-else class="mr-2 h-4 w-4" />
						{{ t(`sharedButtons.${data?.userName ? 'logout' : 'login'}`) }}
					</Button>
				</div>

				<div
					v-if="data?.userName"
					class="flex items-center gap-2 rounded-md px-3 h-9 text-sm border-2 border-gray-700"
				>
					<img
						v-if="data?.avatar"
						:src="data.avatar"
						:alt="data.userName"
						class="h-5 w-5 rounded-full object-cover"
					/>
					<span class="font-medium">{{ data.userName }}</span>
				</div>
			</div>
		</CardFooter>

		<Dialog v-if="withSettings" v-model:open="showSettings">
			<DialogContent class="sm:max-w-[425px]">
				<DialogHeader>
					<DialogTitle>{{ title }}</DialogTitle>
				</DialogHeader>

				<slot name="settings" />

				<DialogFooter>
					<Button variant="outline" @click="showSettings = false">
						{{ t('sharedButtons.close') }}
					</Button>
					<Button v-if="save" variant="default" @click="saveSettings">
						{{ t('sharedButtons.save') }}
					</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	</Card>
</template>
