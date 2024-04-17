<script setup lang="ts">
import { IconExternalLink } from '@tabler/icons-vue';
import { NButton, NTooltip } from 'naive-ui';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { useProfile, useTwitchGetUsers } from '@/api';
import { storeToRefs } from 'pinia';

const { data: profileData } = storeToRefs(useProfile());
const { t } = useI18n();

const selectedUserId = computed(() => {
	return (profileData.value?.selectedDashboardId ?? profileData?.value?.id) || '';
});
const selectedDashboardTwitchUser = useTwitchGetUsers({
	ids: selectedUserId,
});

const publicPageHref = computed<string>(() => {
	if (!profileData.value || !selectedDashboardTwitchUser.data.value?.users.length) return '';

	const login = selectedDashboardTwitchUser.data.value.users.at(0)!.login;

	return `${window.location.origin}/p/${login}`;
});
</script>

<template>
	<n-tooltip>
		<template #trigger>
			<n-button tag="a" circle quaternary target="_blank" :href="publicPageHref">
				<IconExternalLink />
			</n-button>
		</template>
		{{ t('navbar.publicPage') }}
	</n-tooltip>
</template>
