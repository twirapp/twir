<script setup lang="ts">
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
	<UiDialogContent
		v-if="windowWidth > 800"
		class="max-w-3xl max-h-[90dvh] overflow-auto rounded-2xl outline-hidden"
		v-bind="$attrs"
		@interact-outside="onInteractOutside"
	>
		<slot />
	</UiDialogContent>
	<UiSheetContent v-else side="bottom" class="max-h-[90dvh] overflow-auto" v-bind="$attrs" @interact-outside="onInteractOutside">
		<slot />
	</UiSheetContent>
</template>
