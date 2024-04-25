<script setup lang='ts'>
import { useI18n } from 'vue-i18n';

import CommunityUsers from '@/features/community-users/community-users.vue';
import PageLayout from '@/layout/page-layout.vue';

const { t } = useI18n();
</script>

<template>
	<page-layout is-not-tabs>
		<template #title>
			{{ t('community.users.title') }}
		</template>
		<template #content>
			<community-users />
		</template>
	</page-layout>
</template>
