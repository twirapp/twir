<script setup lang="ts">
import { Settings } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import type { FunctionalComponent } from 'vue'

import { useUserAccessFlagChecker } from '@/api'
import Card from '@/components/card/card.vue'
import { Button } from '@/components/ui/button'
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
	<Card :title="title" :icon="icon" :icon-stroke="iconStroke" :icon-fill="iconFill" class="h-full" icon-width="28px" icon-height="28px">
		<template #content>
			<p>{{ description }}</p>
		</template>
		<template #footer>
			<Button
				v-if="showSettings"
				:disabled="!userCanManageGames"
				variant="secondary"
				@click="$emit('openSettings')"
			>
				{{ t('sharedButtons.settings') }}
				<Settings class="ml-2 h-4 w-4" />
			</Button>
		</template>
	</Card>
</template>
