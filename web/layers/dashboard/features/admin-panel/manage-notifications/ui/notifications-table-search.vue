<script setup lang="ts">


import { useNotificationsFilters } from '../composables/use-notifications-filters.js'

import SearchBar from '#layers/dashboard/components/search-bar.vue'

import { NotificationType } from '~/gql/graphql.js'

const { t } = useI18n()

const filters = useNotificationsFilters()
</script>

<template>
	<div class="flex gap-2 max-sm:w-full">
		<SearchBar v-model="filters.searchInput.value" />
		<UiSelect v-model="filters.filterInput.value">
			<UiSelectTrigger class="h-9 w-[120px]">
				<UiSelectValue />
			</UiSelectTrigger>
			<UiSelectContent align="end">
				<UiSelectGroup>
					<UiSelectItem :value="NotificationType.Global">
						{{ t('adminPanel.notifications.globals') }}
					</UiSelectItem>
					<UiSelectItem :value="NotificationType.User">
						{{ t('adminPanel.notifications.users') }}
					</UiSelectItem>
				</UiSelectGroup>
			</UiSelectContent>
		</UiSelect>
	</div>
</template>
