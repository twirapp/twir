<script lang="ts" setup>
const currentRoute = useRoute()

const routes = computed(() => {
	return [
		{
			name: 'Commands',
			icon: 'lucide:package',
			to: `/p/${currentRoute.params.channelName}`,
		},
		{
			name: 'Song requests',
			icon: 'lucide:audio-lines',
			to: `/p/${currentRoute.params.channelName}/songs-requests`,
		},
		{
			name: 'TTS Profiles',
			icon: 'lucide:square-activity',
			to: `/p/${currentRoute.params.channelName}/tts`,
		},
		{
			name: 'Users',
			icon: 'lucide:users',
			to: `/p/${currentRoute.params.channelName}/users`,
		},
	]
})
</script>

<template>
	<UiSidebarGroup>
		<UiSidebarGroupContent>
			<UiSidebarMenu>
				<UiSidebarMenuItem>
					<NuxtLink
						v-for="route in routes"
						:key="route.name"
						:to="route.to"
						custom #="{ isActive, href, navigate }"
						prefetch-on="interaction"
					>
						<UiSidebarMenuButton as="a" :href="href" :variant="isActive ? 'active' : 'default'" @click="navigate">
							<Icon :name="route.icon" />
							{{ route.name }}
						</UiSidebarMenuButton>
					</NuxtLink>
				</UiSidebarMenuItem>
			</UiSidebarMenu>
		</UiSidebarGroupContent>
	</UiSidebarGroup>
</template>
