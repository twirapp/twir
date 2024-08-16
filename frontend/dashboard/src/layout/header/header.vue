<script setup lang="ts">
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core'

import HeaderDesktop from './header-desktop.vue'
import HeaderMobile from './header-mobile.vue'

import TwirLogo from '@/components/twir-logo.vue'

const breakPoints = useBreakpoints(breakpointsTailwind)
const isDesktopWindow = breakPoints.greaterOrEqual('md')
</script>

<template>
	<div class="header">
		<div class="content">
			<a href="/" class="logo">
				<TwirLogo style="width: 36px; height: 36px; display: flex" />
			</a>

			<template v-if="isDesktopWindow">
				<HeaderDesktop />
			</template>

			<template v-if="!isDesktopWindow">
				<HeaderMobile />
			</template>
		</div>
	</div>
</template>

<style scoped>
.header {
	height: 100%;
	padding-top: 10px;
	padding-left: 10px;
	padding-right: 10px;
}

.content .logo {
	display: flex;
	align-items: center;
	text-decoration: none;
	padding-right: 14px;
}

.content > div {
	display: flex;
	justify-content: flex-start;
	gap: 5px;
	height: 45px;
}

.content {
	display: flex;
	justify-content: space-between;
	align-items: center;
	height: 45px;
	column-gap: 12px;
}
</style>
