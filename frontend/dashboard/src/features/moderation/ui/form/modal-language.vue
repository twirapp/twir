<script setup lang="ts">
import { Search } from 'lucide-vue-next'
import { useField } from 'vee-validate'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { useEditableItem } from './helpers.ts'
import Button from '../../../../components/ui/button/Button.vue'
import Input from '../../../../components/ui/input/Input.vue'

import { useModerationAvailableLanguages } from '@/api'
import { ScrollArea } from '@/components/ui/scroll-area'

const { editableItem } = useEditableItem()
const { data: availableLanguages } = useModerationAvailableLanguages()
const { t } = useI18n()

const sourceSearch = ref('')
const targetSearch = ref('')

const allLanguages = computed(() =>
	availableLanguages?.value?.langs.map(l => ({
		label: l.name,
		value: l.code.toString(),
	})) ?? [],
)

const deniedLanguages = computed({
	get: () => editableItem.value?.data?.deniedChatLanguages ?? [],
	set: (value) => {
		if (editableItem.value?.data) {
			editableItem.value.data.deniedChatLanguages = value
		}
	},
})

const allowedLanguages = computed(() =>
	allLanguages.value.filter(lang => !deniedLanguages.value.includes(lang.value)),
)

const filteredAllowedLanguages = computed(() =>
	allowedLanguages.value.filter(lang =>
		lang.label.toLowerCase().includes(sourceSearch.value.toLowerCase()),
	),
)

const filteredDeniedLanguages = computed(() =>
	allLanguages.value
		.filter(lang => deniedLanguages.value.includes(lang.value))
		.filter(lang =>
			lang.label.toLowerCase().includes(targetSearch.value.toLowerCase()),
		),
)

function moveToDisallowed(value: string) {
	deniedLanguages.value = [...deniedLanguages.value, value]
}

function moveToAllowed(value: string) {
	deniedLanguages.value = deniedLanguages.value.filter(v => v !== value)
}

const field = useField('deniedChatLanguages')

watch(deniedLanguages, () => {
	field.setValue(deniedLanguages.value)
})
</script>

<template>
	<div class="grid grid-cols-2 gap-4">
		<div class="flex flex-col gap-2 p-4">
			<h3 class="font-medium">
				{{ t('moderation.types.language.allowedLanguages') }} ({{ allowedLanguages.length }})
			</h3>
			<div class="relative">
				<Search class="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
				<Input
					v-model="sourceSearch"
					class="pl-8"
					placeholder="Search language..."
				/>
			</div>
			<ScrollArea class="h-[300px] w-full rounded-md border" type="auto">
				<div class="p-2 space-y-2">
					<Button
						v-for="lang in filteredAllowedLanguages"
						:key="lang.value"
						variant="ghost"
						class="w-full justify-start"
						@click="moveToDisallowed(lang.value)"
					>
						{{ lang.label }}
					</Button>
				</div>
			</ScrollArea>
		</div>

		<!-- Disallowed Languages -->
		<div class="flex flex-col gap-2 p-4">
			<h3 class="font-medium">
				{{ t('moderation.types.language.disallowedLanguages') }} ({{ deniedLanguages.length }})
			</h3>
			<div class="relative">
				<Search class="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
				<Input
					v-model="targetSearch"
					class="pl-8"
					placeholder="Search language..."
				/>
			</div>
			<ScrollArea class="h-[300px] w-full rounded-md border" type="auto">
				<div class="p-2 space-y-2">
					<Button
						v-for="lang in filteredDeniedLanguages"
						:key="lang.value"
						variant="ghost"
						class="w-full justify-start"
						@click="moveToAllowed(lang.value)"
					>
						{{ lang.label }}
					</Button>
				</div>
			</ScrollArea>
		</div>
	</div>
</template>
