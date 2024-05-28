<script setup lang="ts">
import { IconBrandDiscord, IconBrandGithub, IconUsers } from '@tabler/icons-vue'
import { DISCORD_INVITE_URL, GITHUB_REPOSITORY_URL } from '@twir/brand'
import { NMenu } from 'naive-ui'
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'

import type { MenuDividerOption, MenuOption } from 'naive-ui'

import { renderIcon } from '@/helpers/renderIcon'
import { usePublicPageHref } from '@/layout/use-public-page-href.js'

const { t } = useI18n()
const activeKey = 'empty'
const publicPageHref = usePublicPageHref()

function getOptionUrlParams(path: string, isExternal: boolean) {
	if (isExternal) {
		return {
			href: path,
			target: '__blank',
		}
	}

	return {
		to: {
			path,
		},
	}
}

const menuOptions = computed<(MenuOption | MenuDividerOption)[]>(() => {
	return [
		{
			label: 'Github',
			icon: renderIcon(IconBrandGithub),
			path: GITHUB_REPOSITORY_URL,
			isExternal: true,
		},
		{
			label: 'Discord',
			icon: renderIcon(IconBrandDiscord),
			path: DISCORD_INVITE_URL,
			isExternal: true,
		},
		{
			label: t('navbar.publicPage'),
			icon: renderIcon(IconUsers),
			path: publicPageHref.value,
			isExternal: true,
		},
	].map((item) => ({
		...item,
		key: item.path ?? item.label,
		label: !item.path
			? item.label ?? undefined
			: () => h(
				item.isExternal ? 'a' : RouterLink,
				// eslint-disable-next-line
				// @ts-ignore
				getOptionUrlParams(item.path, item.isExternal),
				{
					default: () => item.label,
				},
			),
	}))
})
</script>

<template>
	<NMenu
		v-model:value="activeKey"
		:collapsed-width="64"
		:collapsed-icon-size="22"
		:options="menuOptions"
		style="padding-bottom: 0"
	/>
</template>

<style scoped>
:deep(.n-menu-item-content-header) {
	@apply self-stretch flex items-center;
}
</style>
