<script setup lang="ts">
import { ref, watch } from 'vue'

import { useFaceitIntegration } from '@/api/index.js'
import IconFaceit from '@/assets/integrations/faceit.svg?use'
import OauthComponent from '@/components/integrations/variants/oauth.vue'
import { Label } from '@/components/ui/label'
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'

const manager = useFaceitIntegration()
const { data } = manager.useData()
const logout = manager.useLogout()
const { data: authLink } = manager.useAuthLink()
const updater = manager.update!()

const game = ref('cs2')

const gameOptions = [
	{ label: 'Counter-Strike', value: 'cs2' },
	{ label: 'Counter-Strike: Global Offensive', value: 'csgo' },
	{ label: 'Dota 2', value: 'dota2' },
]

watch(data, (newData) => {
	if (!newData?.game) return
	game.value = newData.game
})

async function save() {
	await updater.mutateAsync({ game: game.value })
}
</script>

<template>
	<OauthComponent
		title="Faceit"
		:data="data"
		:logout="() => logout.mutateAsync({})"
		:authLink="authLink?.link"
		:icon="IconFaceit"
		:withSettings="true"
		:save="save"
	>
		<template #description>
			<i18n-t keypath="integrations.faceit.info">
				<b class="text-foreground">$(faceit.elo)</b>
				<b class="text-foreground">$(faceit.lvl)</b>
			</i18n-t>
		</template>
		<template #settings>
			<div class="space-y-4 py-2">
				<div class="space-y-2">
					<Label>Game</Label>
					<Select v-model:modelValue="game">
						<SelectTrigger>
							<SelectValue :placeholder="gameOptions[0].label" />
						</SelectTrigger>
						<SelectContent>
							<SelectItem
								v-for="option in gameOptions"
								:key="option.value"
								:value="option.value"
							>
								{{ option.label }}
							</SelectItem>
						</SelectContent>
					</Select>
				</div>
			</div>
		</template>
	</OauthComponent>
</template>

<style scoped>
/* Removed naive-ui specific styling */
</style>
