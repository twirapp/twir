<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import { useNotificationsFilters } from '../composables/use-notifications-filters.js'

import SearchBar from '@/components/search-bar.vue'
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue
} from '@/components/ui/select'
import { NotificationType } from '@/gql/graphql.js'

const { t } = useI18n()

const filters = useNotificationsFilters()
</script>

<template>
	<div class="flex gap-2 max-lg:w-full">
		<SearchBar v-model="filters.searchInput" />
		<Select v-model="filters.filterInput">
			<SelectTrigger class="h-9 w-[120px]">
				<SelectValue />
			</SelectTrigger>
			<SelectContent align="end">
				<SelectGroup>
					<SelectItem :value="NotificationType.Global">
						{{ t('adminPanel.notifications.globals') }}
					</SelectItem>
					<SelectItem :value="NotificationType.User">
						{{ t('adminPanel.notifications.users') }}
					</SelectItem>
				</SelectGroup>
			</SelectContent>
		</Select>
	</div>
</template>
