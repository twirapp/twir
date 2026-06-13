<script setup lang="ts">
import CommandsPrefix from './ui/commands-prefix.vue'

import PageLayout from '@/layout/page-layout.vue'
</script>

<template>
	<PageLayout>
		<template #title>
			Bot Settings
		</template>

		<template #content>
			<div class="flex flex-col gap-12">
				<CommandsPrefix />
			</div>
		</template>
	</PageLayout>
</template>
