<script setup lang="ts">
import { computed } from 'vue';

import JpFlag from '@/assets/flags/jp.svg?use';
import RuFlag from '@/assets/flags/ru.svg?use';
import UaFlag from '@/assets/flags/ua.svg?use';
import UsFlag from '@/assets/flags/us.svg?use';

// TODO: fix locale name
const flags = {
	ru: RuFlag,
	uk: UaFlag,
	en: UsFlag,
	ja: JpFlag,
} as const;

export type Locale = keyof typeof flags

const props = defineProps<{
	locale: Locale;
}>();

const flag = computed(() => flags[props.locale]);
</script>

<template>
	<flag />
</template>
