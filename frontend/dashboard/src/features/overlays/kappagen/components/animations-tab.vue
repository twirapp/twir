<script setup lang="ts">
import { useFieldArray } from 'vee-validate'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Switch } from '@/components/ui/switch'
import type { KappagenOverlayAnimationsSettings } from '@/gql/graphql'
import { PlayIcon, SettingsIcon } from 'lucide-vue-next'

const { fields: animations } = useFieldArray<KappagenOverlayAnimationsSettings>('animations')
</script>

<template>
	<div class="space-y-6">
		<Card>
			<CardHeader>
				<CardTitle>Animations</CardTitle>
			</CardHeader>
			<CardContent class="space-y-4">
				<div class="grid grid-cols-1 md:grid-cols-2 gap-2">
					<div
						v-for="animation of animations"
						:key="animation.key"
						class="flex flex-row items-center justify-between bg-background/60 p-2 rounded-md"
					>
						<div class="flex gap-2 items-center">
							<button class="p-2 rounded-md bg-indigo-400/70">
								<PlayIcon class="size-4" />
							</button>
							{{ animation.value.style }}
						</div>

						<div class="flex gap-2 items-center">
							<button
								class="p-1 border-border border rounded-md bg-zinc-600/50 hover:bg-zinc-600/30 transition-colors"
								v-if="animation.value.prefs"
							>
								<SettingsIcon class="size-4" />
							</button>
							<Switch
								:value="animation.value.enabled"
								@update:value="
									(val) => animation.handleChange({ ...animation.value, enabled: val })
								"
							/>
						</div>
					</div>
				</div>
			</CardContent>
		</Card>
	</div>
</template>
