<script lang="ts" setup>
import type { ToasterProps } from 'vue-sonner'
import { Toaster as Sonner } from 'vue-sonner'
import 'vue-sonner/style.css'

const props = defineProps<ToasterProps>()
</script>

<template>
	<ClientOnly>
		<Sonner
			class="toaster group"
			v-bind="props"
			:toast-options="{
				unstyled: true,
				classes: {
					toast:
						'group toast group-[.toaster]:bg-background group-[.toaster]:text-foreground group-[.toaster]:border-border group-[.toaster]:border group-[.toaster]:shadow-lg group-[.toaster]:p-4 group-[.toaster]:rounded-lg group-[.toaster]:flex group-[.toaster]:w-[356px] group-[.toaster]:items-center group-[.toaster]:text-sm group-[.toaster]:gap-3 group-[.toaster]:backdrop-blur-md ',
					title: 'group-[.toast]:font-bold',
					description: 'group-[.toast]:text-white/65 group-[.toast]:text-sm',
					actionButton: 'group-[.toast]:bg-primary group-[.toast]:text-primary-foreground',
					cancelButton: 'group-[.toast]:bg-muted group-[.toast]:text-muted-foreground',
				},
			}"
		/>
	</ClientOnly>
</template>
