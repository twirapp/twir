<script setup lang="ts">
import { useI18n } from 'vue-i18n';

import BadgesForm from './components/badges-form.vue';
import BadgesList from './components/badges-list.vue';

const { t } = useI18n();
</script>

<template>
	<div class="w-full flex flex-wrap gap-6">
		<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
			{{ t('adminPanel.manageBadges.formTitle') }}
		</h4>
		<badges-form />
		<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
			{{ t('adminPanel.manageBadges.title') }}
		</h4>
		<badges-list />
	</div>
</template>
