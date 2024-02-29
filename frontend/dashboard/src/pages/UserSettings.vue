<script setup lang="ts">
import { useThemeVars } from 'naive-ui';
import { TabsList, TabsRoot, TabsTrigger, TabsContent } from 'radix-vue';

import AccountSettings from './user-settings/account.vue';
import PublicSettings from './user-settings/public.vue';

import { useTheme } from '@/composables/use-theme';

const themeVars = useThemeVars();
const { theme } = useTheme();
</script>

<template>
	<TabsRoot default-value="account">
		<div
			class="after:inset-0 after:bottom-0 after:block after:h-px after:w-full after:content-['']"
			:class="{
				'after:bg-white/[.15]': theme === 'dark',
				'after:bg-zinc-600/[.15]': theme === 'light'
			}"
			:style="{'background-color': themeVars.cardColor}"
		>
			<div class="container flex flex-col pt-9 gap-2">
				<h1 class="text-4xl">
					User settings
				</h1>
				<div class="flex gap-2">
					<TabsList class="flex overflow-x-auto -mb-px">
						<TabsTrigger class="tabs-trigger" value="account">
							Account
						</TabsTrigger>
						<TabsTrigger class="tabs-trigger" value="public">
							Public page
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
		</div>
	</TabsRoot>
</template>

<style scoped>
.tabs-trigger {
	@apply relative z-[1] flex whitespace-nowrap px-3 py-4 text-sm  transition-colors before:absolute before:left-0 before:top-2 before:-z-[1] before:block before:h-9 before:w-full before:rounded-md before:transition-colors before:content-[''] hover:text-white hover:before:bg-zinc-800 data-[state=active]:after:absolute data-[state=active]:after:bottom-0 data-[state=active]:after:left-2 data-[state=active]:after:right-2 data-[state=active]:after:block data-[state=active]:after:h-0 data-[state=active]:after:border-b-2 data-[state=active]:after:border-white data-[state=active]:after:content-[''] data-[state=active]:after:rounded-t-sm font-medium
}
</style>
