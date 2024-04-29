<script setup lang="ts">
import { useThemeVars } from 'naive-ui'

const themeVars = useThemeVars()
</script>

<template>
	<div
		class="flex-wrap w-full border rounded-md" :style="{
			backgroundColor: themeVars.cardColor,
			color: themeVars.textColor2,
		}"
	>
		<slot />
	</div>
</template>
