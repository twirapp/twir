<script setup lang="ts">
import { useWindowSize } from '@vueuse/core'

import { DialogContent } from '@/components/ui/dialog'
import { SheetContent } from '@/components/ui/sheet'

defineOptions({
	inheritAttrs: false,
})

const { width: windowWidth } = useWindowSize()

function onInteractOutside(event: any) {
	if ((event.target as HTMLElement)?.closest('[role="dialog"]')) return
	event.preventDefault()
}
</script>

<template>
	<DialogContent
		v-if="windowWidth > 800"
		class="max-w-3xl max-h-[90dvh] overflow-auto rounded-2xl outline-none"
		v-bind="$attrs"
		@interact-outside="onInteractOutside"
	>
		<slot />
	</DialogContent>
	<SheetContent v-else side="bottom" class="max-h-[90dvh] overflow-auto" v-bind="$attrs" @interact-outside="onInteractOutside">
		<slot />
	</SheetContent>
</template>
