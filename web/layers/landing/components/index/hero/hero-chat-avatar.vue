<script setup lang="ts">
import { computed } from 'vue'

import TwirLogo from '../../../../../../libs/brand/src/logo.svg?url'

export type ColorVariant =
  | 'blue'
  | 'purple'
  | 'red'
  | 'turquoise'
  | 'orange'
  | 'lime'
  | 'gray'
  | 'pink'

interface Props {
	variant?: ColorVariant
	isBot?: boolean
}

const props = withDefaults(defineProps<Props>(), {
	variant: 'purple',
})

const DARK = '#060002'
const WHITE = '#FFF'

const variantStyles: Record<ColorVariant, { avatar: string, bg: string }> = {
	blue: { bg: '#1E69FF', avatar: WHITE },
	red: { bg: '#FD6675', avatar: DARK },
	purple: { bg: '#8205B4', avatar: WHITE },
	turquoise: { bg: '#00C8AF', avatar: DARK },
	orange: { bg: '#ffca5f', avatar: DARK },
	lime: { bg: '#beff00', avatar: DARK },
	gray: { bg: '#d2d2e6', avatar: DARK },
	pink: { bg: '#fc0fc0', avatar: WHITE },
}

const variant = computed(() => variantStyles[props.variant])
</script>

<template>
	<img v-if="isBot" :src="TwirLogo" class="w-9 h-9 shrink-0" />
	<svg v-else fill="none" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 36 36" class="h-9 w-9 shrink-0">
		<path d="M36 18a18 18 0 1 1-36 0 18 18 0 0 1 36 0Z" :fill="variant.bg"></path>
		<path
			fill-rule="evenodd"
			clip-rule="evenodd"
			d="M11.25 14A6.8 6.8 0 0 1 18 7.25 6.8 6.8 0 0 1 24.75 14a6.81 6.81 0 0 1-5.21 6.57c.18.3.46.58.83.82.8.52 2.04.86 3.63.86 1.2 0 2.14.6 2.75 1.3.6.68 1 1.59 1 2.45v3.25h-2.5V26c0-.14-.1-.48-.38-.8-.26-.3-.56-.45-.87-.45-1.9 0-3.67-.4-5-1.26-.31-.21-.61-.45-.88-.71-.22.25-.48.48-.79.69-1.31.9-3.24 1.28-5.8 1.28a1.46 1.46 0 0 0-.4.38c-.13.16-.23.34-.3.51-.06.18-.08.3-.08.36v3.25h-2.5V26c0-.85.4-1.73.88-2.38.45-.6 1.27-1.37 2.37-1.37 2.43 0 3.74-.38 4.42-.84.3-.21.52-.47.66-.81a6.81 6.81 0 0 1-5.33-6.6ZM18 9.75A4.3 4.3 0 0 0 13.75 14 4.3 4.3 0 0 0 18 18.25 4.3 4.3 0 0 0 22.25 14 4.3 4.3 0 0 0 18 9.75Z"
			:fill="variant.avatar"
		></path>
	</svg>
</template>
