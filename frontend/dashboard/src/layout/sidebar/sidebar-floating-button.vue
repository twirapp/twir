<script setup lang="ts">
import { Menu } from 'lucide-vue-next'

import { Button } from '@/components/ui/button'
import { useSidebar } from '@/components/ui/sidebar'

const { isMobile, setOpenMobile } = useSidebar()
</script>

<template>
	<Button
		v-if="isMobile"
		class="fixed right-[2rem] bottom-[2rem] z-50 h-14 w-14"
		variant="secondary"
		size="icon"
		@click="setOpenMobile(true)"
	>
		<Menu class="size-8" />
	</Button>
</template>
