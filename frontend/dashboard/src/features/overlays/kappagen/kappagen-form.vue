<script setup lang="ts">
import EventsTab from '@/features/overlays/kappagen/components/events-tab.vue'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import SettingsTab from '@/features/overlays/kappagen/components/settings-tab.vue'
import AnimationsTab from '@/features/overlays/kappagen/components/animations-tab.vue'
</script>

<template>
	<div class="w-[45%]">
		<Tabs default-value="settings">
			<TabsList class="grid w-full grid-cols-3">
				<TabsTrigger value="settings">Settings</TabsTrigger>
				<TabsTrigger value="animations">Animations</TabsTrigger>
				<TabsTrigger value="events">Events</TabsTrigger>
			</TabsList>

			<TabsContent force-mount value="settings" class="mt-6">
				<SettingsTab />
			</TabsContent>

			<TabsContent force-mount value="animations" class="mt-6">
				<AnimationsTab />
			</TabsContent>

			<TabsContent force-mount value="events" class="mt-6">
				<EventsTab />
			</TabsContent>
		</Tabs>
	</div>
</template>
