<script setup lang="ts">
import { useRouteQuery } from '@vueuse/router';
import { useThemeVars } from 'naive-ui';
import { TabsList, TabsRoot, TabsTrigger, TabsContent } from 'radix-vue';
import { onBeforeMount, ref, type Component } from 'vue';
import { watch } from 'vue';
import { useRouter } from 'vue-router';

import { useTheme } from '@/composables/use-theme';

const router = useRouter();
const themeVars = useThemeVars();
const { theme } = useTheme();

export interface PageLayoutProps {
	activeTab: string
	tabs: PageLayoutTab[]
}

export interface PageLayoutTab {
	name: string
	title: string
	component: Component
}

const props = defineProps<PageLayoutProps>();

const activeTab = ref(props.activeTab);

const queryActiveTab = useRouteQuery<string>('tab');
watch(queryActiveTab, setTab);

onBeforeMount(() => {
	setTab();
	onChangeTab(activeTab.value, true);
});

function setTab(): void {
	const tabValue = (queryActiveTab.value ?? props.activeTab).toLowerCase();
	if (props.tabs.some((tab) => tab.name === tabValue)) {
		activeTab.value = tabValue;
	}
}

function onChangeTab(tab: string, replace = false): void {
	router.push({ query: { tab }, replace });
}
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
					<slot name="title" />
				</h1>
				<div class="flex gap-2">
					<TabsList class="flex overflow-x-auto -mb-px">
						<TabsTrigger v-for="tab of props.tabs" :key="tab.name" class="tabs-trigger" :value="tab.name">
							{{ tab.title }}
						</TabsTrigger>
					</TabsList>
				</div>
			</div>
		</div>
		<div class="container py-8">
			<TabsContent v-for="tab of props.tabs" :key="tab.name" :value="tab.name">
				<component :is="tab.component" />
			</TabsContent>
		</div>
	</TabsRoot>
</template>

<style scoped>
.tabs-trigger {
	@apply relative z-[1] flex whitespace-nowrap px-3 py-4 text-sm  transition-colors before:absolute before:left-0 before:top-2 before:-z-[1] before:block before:h-9 before:w-full before:rounded-md before:transition-colors before:content-[''] hover:text-white hover:before:bg-zinc-800 data-[state=active]:after:absolute data-[state=active]:after:bottom-0 data-[state=active]:after:left-2 data-[state=active]:after:right-2 data-[state=active]:after:block data-[state=active]:after:h-0 data-[state=active]:after:border-b-2 data-[state=active]:after:border-white data-[state=active]:after:content-[''] data-[state=active]:after:rounded-t-sm font-medium
}
</style>
