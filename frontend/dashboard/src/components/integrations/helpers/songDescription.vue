<script setup lang="ts">
import { useThemeVars } from 'naive-ui';

const themeVars = useThemeVars();
</script>

<template>
	<i18n-t
		keypath="integrations.songServicesInfo"
	>
		<b :style="{ color: themeVars.successColor }">!song</b>
	</i18n-t>
</template>
