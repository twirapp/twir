<script setup lang="ts">
import { useThemeVars } from 'naive-ui';
import { TabsList, TabsRoot, TabsTrigger, TabsContent } from 'radix-vue';
import { ref } from 'vue';
import { onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import AccountSettings from './user-settings/account.vue';
import Notifications from './user-settings/notifications.vue';
import PublicSettings from './user-settings/public.vue';

import { useTheme } from '@/composables/use-theme';

const { t } = useI18n();
const themeVars = useThemeVars();
const { theme } = useTheme();

const router = useRouter();

const AVAILABLE_TABS = ['account', 'public', 'notifications'];
const activeTab = ref('account');

function onChangeTab(tab: string, replace = false): void {
	router.push({ query: { ...router.currentRoute.value.query, tab }, replace });
}

onMounted(() => {
	const tab = router.currentRoute.value.query.tab as string;
	if (AVAILABLE_TABS.includes(tab)) {
		activeTab.value = tab;
	} else {
		onChangeTab(activeTab.value, true);
	}
});
</script>

<template>
	<TabsRoot v-model="activeTab" @update:model-value="onChangeTab">
		<div
			class="after:inset-0 after:bottom-0 after:block after:h-px after:w-full after:content-['']"
			:class="[theme === 'dark' ? 'after:bg-white/[.15]' : 'after:bg-zinc-600/[.15]']"
			:style="{ 'background-color': themeVars.cardColor }"
		>
			<div class="container flex flex-col pt-9 gap-2">
				<h1 class="text-4xl">
					{{ t('userSettings.title') }}
				</h1>
				<div class="flex gap-2">
					<TabsList class="flex overflow-x-auto -mb-px">
						<TabsTrigger class="tabs-trigger" value="account">
							{{ t('userSettings.account.title') }}
						</TabsTrigger>
						<TabsTrigger class="tabs-trigger" value="public">
							{{ t('userSettings.public.title') }}
						</TabsTrigger>
						<TabsTrigger class="tabs-trigger" value="notifications">
							{{ t('userSettings.notifications.title') }}
						</TabsTrigger>
					</TabsList>
				</div>
			</div>
		</div>
		<div class="container py-8">
			<TabsContent value="account">
				<AccountSettings />
			</TabsContent>
			<TabsContent value="public">
				<PublicSettings />
			</TabsContent>
			<TabsContent value="notifications">
				<Notifications />
			</TabsContent>
		</div>
	</TabsRoot>
</template>

<style scoped>
.tabs-trigger {
	@apply relative z-[1] flex whitespace-nowrap px-3 py-4 text-sm  transition-colors before:absolute before:left-0 before:top-2 before:-z-[1] before:block before:h-9 before:w-full before:rounded-md before:transition-colors before:content-[''] hover:text-white hover:before:bg-zinc-800 data-[state=active]:after:absolute data-[state=active]:after:bottom-0 data-[state=active]:after:left-2 data-[state=active]:after:right-2 data-[state=active]:after:block data-[state=active]:after:h-0 data-[state=active]:after:border-b-2 data-[state=active]:after:border-white data-[state=active]:after:content-[''] data-[state=active]:after:rounded-t-sm font-medium
}
</style>
