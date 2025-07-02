<script setup lang="ts">
import { NButton, NResult } from 'naive-ui';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

const { t } = useI18n();
const router = useRouter();

function goToDashboard() {
	router.push('/dashboard');
}
</script>

<template>
	<div class="flex items-center justify-center h-full">
		<n-result
			status="404"
			:title="t('pageNotFound.title')"
			:description="t('pageNotFound.description')"
		>
			<template #footer>
				<n-button @click="goToDashboard">
					{{ t('sharedButtons.goToDashboard') }}
				</n-button>
			</template>
		</n-result>
	</div>
</template>
