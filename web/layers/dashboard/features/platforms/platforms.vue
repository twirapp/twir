<script setup lang="ts">
import CommandsPrefix from './ui/commands-prefix.vue'

import PlatformBindings from '~~/layers/dashboard/features/channel-platforms/ui/platform-bindings.vue'
import PageLayout from '~~/layers/dashboard/layout/page-layout.vue'
</script>

<template>
	<PageLayout>
		<template #title>
			Platforms
		</template>

		<template #content>
			<div class="flex flex-col gap-12">
				<PlatformBindings />
				<CommandsPrefix />
			</div>
		</template>
	</PageLayout>
</template>
