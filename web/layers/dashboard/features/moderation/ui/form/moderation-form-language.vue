<script setup lang="ts">
import { Search, Trash } from 'lucide-vue-next'
import { useField } from 'vee-validate'
import { computed, ref } from 'vue'


import { useModerationAvailableLanguages } from '#layers/dashboard/api/moderation'





const languagesApi = useModerationAvailableLanguages()
const { data: availableLanguages } = languagesApi.query()
const { t } = useI18n()

const sourceSearch = ref('')
const targetSearch = ref('')

const allLanguages = computed(() =>
	availableLanguages?.value?.moderationLanguagesAvailableLanguages.languages.map(l => ({
		label: l.name,
		value: l.iso_639_1,
	})) ?? [],
)

const { value: deniedLanguages, setValue: setDenyList } = useField<string[]>('deniedChatLanguages')
const { value: excludedWords } = useField<string[]>('languageExcludedWords')

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
	setDenyList([...deniedLanguages.value, value])
}

function moveToAllowed(value: string) {
	setDenyList(deniedLanguages.value.filter(v => v !== value))
}

const maxExcludedWords = 1000
</script>

<template>
	<div class="flex gap-2 w-full flex-col">
		<div class="grid grid-cols-2 gap-4">
			<div class="flex flex-col gap-2 p-4">
				<h3 class="font-medium">
					{{ t('moderation.types.language.allowedLanguages') }} ({{ allowedLanguages.length }})
				</h3>
				<div class="relative">
					<Search class="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
					<UiInput
						v-model="sourceSearch"
						class="pl-8"
						placeholder="Search language..."
					/>
				</div>
				<UiScrollArea class="h-[300px] w-full rounded-md border" type="auto">
					<div class="p-2 space-y-2">
						<UiButton
							v-for="lang in filteredAllowedLanguages"
							:key="lang.value"
							type="button"
							variant="ghost"
							class="w-full justify-start"
							@click="moveToDisallowed(lang.value)"
						>
							{{ lang.label }}
						</UiButton>
					</div>
				</UiScrollArea>
			</div>

			<!-- Disallowed Languages -->
			<div class="flex flex-col gap-2 p-4">
				<h3 class="font-medium">
					{{ t('moderation.types.language.disallowedLanguages') }} ({{ deniedLanguages.length }})
				</h3>
				<div class="relative">
					<Search class="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
					<UiInput
						v-model="targetSearch"
						class="pl-8"
						placeholder="Search language..."
					/>
				</div>
				<UiScrollArea class="h-[300px] w-full rounded-md border" type="auto">
					<div class="p-2 space-y-2">
						<UiButton
							v-for="lang in filteredDeniedLanguages"
							:key="lang.value"
							variant="ghost"
							class="w-full justify-start"
							@click="moveToAllowed(lang.value)"
						>
							{{ lang.label }}
						</UiButton>
					</div>
				</UiScrollArea>
			</div>
		</div>

		<div class="flex flex-col gap-2 w-full">
			<h3 class="font-medium">
				Ignored words
			</h3>

			<UiAlert v-if="!excludedWords.length">
				<UiAlertDescription>
					Create ignored words for excluded them from detect.
				</UiAlertDescription>
			</UiAlert>

			<div
				v-for="(_, index) of excludedWords"
				:key="index"
				class="flex gap-2 w-full"
			>
				<UiInput
					v-model="excludedWords[index]"
					placeholder="Yes"
					class="flex-1"
				/>

				<UiButton
					type="button"
					variant="destructive"
					size="icon"
					@click="() => {
						excludedWords = excludedWords.filter((_, i) => i !== index)
					}"
				>
					<Trash class="h-4 w-4" />
				</UiButton>
			</div>

			<UiButton
				type="button"
				variant="default"
				class="w-full"
				:disabled="excludedWords.length >= maxExcludedWords"
				@click="() => excludedWords.push('')"
			>
				{{ t('sharedButtons.create') }}
			</UiButton>
		</div>
	</div>
</template>
