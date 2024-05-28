<script setup lang="ts">
import { IconSettings } from '@tabler/icons-vue'
import { NButton } from 'naive-ui'
import { useI18n } from 'vue-i18n'

import type { FunctionalComponent } from 'vue'

import { useUserAccessFlagChecker } from '@/api'
import Card from '@/components/card/card.vue'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

withDefaults(defineProps<{
	description: string
	title: string
	icon: FunctionalComponent
	iconStroke?: number
	showSettings?: boolean
	iconFill?: string
}>(), { showSettings: true })

defineEmits<{
	openSettings: []
}>()

const { t } = useI18n()

const userCanManageGames = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageGames)
</script>

<template>
	<Card :title="title" :icon="icon" :icon-stroke="iconStroke" :icon-fill="iconFill" class="h-full">
		<template #content>
			<p>{{ description }}</p>
		</template>
		<template #footer>
			<NButton
				v-if="showSettings"
				:disabled="!userCanManageGames"
				secondary
				size="large"
				@click="$emit('openSettings')"
			>
				<div class="flex gap-[6px]">
					<span>{{ t('sharedButtons.settings') }}</span>
					<IconSettings />
				</div>
			</NButton>
		</template>
	</Card>
</template>
