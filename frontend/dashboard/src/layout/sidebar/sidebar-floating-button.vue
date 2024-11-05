<script setup lang="ts">
import { Menu } from 'lucide-vue-next'

import { Button } from '@/components/ui/button'
import { useSidebar } from '@/components/ui/sidebar'

const { isMobile, setOpenMobile } = useSidebar()
</script>

<template>
	<Button
		v-if="isMobile"
		class="fixed right-2 bottom-2 z-50"
		variant="ghost"
		size="icon"
		@click="setOpenMobile(true)"
	>
		<Menu />
	</Button>
</template>
