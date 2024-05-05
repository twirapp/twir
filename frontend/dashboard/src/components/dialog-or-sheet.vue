<script setup lang="ts">
import { useWindowSize } from '@vueuse/core'

import { DialogContent } from '@/components/ui/dialog'
import { SheetContent } from '@/components/ui/sheet'

const { width: windowWidth } = useWindowSize()
</script>

<template>
	<DialogContent
		v-if="windowWidth > 800"
		class="max-w-3xl overflow-hidden rounded-2xl outline-none"
	>
		<slot />
	</DialogContent>
	<SheetContent v-else side="bottom">
		<slot />
	</SheetContent>
</template>
