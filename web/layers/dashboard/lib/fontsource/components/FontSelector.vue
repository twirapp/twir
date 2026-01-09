<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from 'vue'

import { Check, ChevronsUpDown, LoaderCircle } from 'lucide-vue-next'
import { type Font, generateFontKey, useFontSource } from '@twir/fontsource'
import { useVirtualizer } from '@tanstack/vue-virtual'




import { cn } from '@/lib/utils'

const props = defineProps<{
	fontFamily: string
	fontWeight: number
	fontStyle: string
}>()

const { t } = useI18n()
const fontSource = useFontSource()
const font = defineModel<Font | null>('font')

const selectedFont = ref<string>('')
const open = ref(false)
const searchQuery = ref('')
const scrollContainerRef = ref<HTMLElement | null>(null)

// Watch for prop changes to sync selectedFont
watch(
	() => props.fontFamily,
	async (newFamily) => {
		if (newFamily && newFamily !== selectedFont.value) {
			selectedFont.value = newFamily

			// Load font immediately when prop changes
			if (!font.value || font.value.id !== newFamily) {
				const loadedFont = await fontSource.loadFont(newFamily, props.fontWeight, props.fontStyle)
				if (loadedFont) {
					font.value = loadedFont
				}
			}
		}
	},
	{ immediate: true }
)

watch(
	() => selectedFont.value,
	(selectedFontId) => {
		if (!selectedFontId) {
			font.value = null
			return
		}

		const foundFont = fontSource.getFont(selectedFontId)
		if (foundFont) {
			font.value = foundFont
		} else {
			// Font not in cache, load it
			fontSource.loadFont(selectedFontId, props.fontWeight, props.fontStyle)
				.then((loadedFont) => {
					if (loadedFont) {
						font.value = loadedFont
					}
				})
		}
	}
)

interface FontOption {
	label: string
	value: string
	fontWeight: number
	fontStyle: string
}

const fontOptions = computed((): FontOption[] => {
	const query = searchQuery.value.toLowerCase().trim()

	return fontSource.fontList.value
		.filter((font) => {
			// Filter by search query
			if (query && !font.family.toLowerCase().includes(query)) {
				return false
			}

			return true
		})
		.map((font) => ({
			label: font.family,
			value: font.id,
			fontWeight: font.weights.includes(400) ? 400 : font.weights[0],
			fontStyle: font.styles.includes('normal') ? 'normal' : font.styles[0],
		}))
})

const selectedFontLabel = computed(() => {
	const selected = fontOptions.value.find((f) => f.value === selectedFont.value)
	return selected?.label || t('overlays.chat.selectFont')
})



// Virtual scrolling setup
const virtualizer = useVirtualizer({
	get count() {
		return fontOptions.value.length
	},
	getScrollElement: () => scrollContainerRef.value,
	estimateSize: () => 36, // Approximate height of each item
	overscan: 10, // Render 10 extra items outside viewport
})

function loadFontPreview(option: FontOption) {
	if (!fontSource.loading.value) {
		fontSource.loadFont(option.value, option.fontWeight, option.fontStyle)
	}
}

function getFontFamily(option: FontOption): string {
	return generateFontKey(option.value, option.fontWeight, option.fontStyle)
}

function selectFont(optionValue: string) {
	selectedFont.value = optionValue
	open.value = false
}

// Scroll to selected item when opening popover
watch(open, async (isOpen) => {
	if (isOpen && selectedFont.value) {
		await nextTick()
		const selectedIndex = fontOptions.value.findIndex(f => f.value === selectedFont.value)
		if (selectedIndex !== -1) {
			virtualizer.value.scrollToIndex(selectedIndex, { align: 'center' })
		}
	}
})

onMounted(async () => {
	// Set initial selected font
	if (props.fontFamily) {
		selectedFont.value = props.fontFamily

		// Load the font if not already loaded
		const loadedFont = await fontSource.loadFont(props.fontFamily, props.fontWeight, props.fontStyle)
		if (loadedFont) {
			font.value = loadedFont
		}
	}
})
</script>

<template>
	<div class="space-y-2">
		<UiPopover v-model:open="open">
			<UiPopoverTrigger as-child>
				<UiButton
					variant="outline"
					role="combobox"
					:aria-expanded="open"
					:disabled="fontSource.loading.value"
					class="w-full justify-between"
				>
					<span v-if="fontSource.loading.value" class="flex items-center gap-2">
						<LoaderCircle class="h-4 w-4 animate-spin" />
						{{ t('sharedTexts.loading') }}
					</span>
					<span v-else class="truncate">{{ selectedFontLabel }}</span>
					<ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
				</UiButton>
			</UiPopoverTrigger>
			<UiPopoverContent class="w-100 p-0">
				<div class="flex flex-col">
					<!-- Search Input -->
					<div class="p-2 border-b">
						<UiInput
							v-model="searchQuery"
							:placeholder="t('overlays.chat.searchFont')"
							class="h-9"
						/>
					</div>

							<!-- Empty state -->
					<div v-if="fontOptions.length === 0" class="p-4 text-center text-sm text-muted-foreground">
						{{ t('overlays.chat.noFontsFound') }}
					</div>

					<!-- Virtual scrolling container -->
					<div
						v-else
						ref="scrollContainerRef"
						class="max-h-80 overflow-auto"
					>
						<div
							:style="{
								height: `${virtualizer.getTotalSize()}px`,
								width: '100%',
								position: 'relative',
							}"
						>
							<div
								v-for="virtualRow in virtualizer.getVirtualItems()"
								:key="virtualRow.index"
								:data-index="virtualRow.index"
								:style="{
									position: 'absolute',
									top: 0,
									left: 0,
									width: '100%',
									height: `${virtualRow.size}px`,
									transform: `translateY(${virtualRow.start}px)`,
								}"
							>
								<div
									class="flex items-center gap-2 px-2 py-2 cursor-pointer hover:bg-accent hover:text-accent-foreground rounded-sm transition-colors"
									:class="{
										'bg-accent text-accent-foreground': selectedFont === fontOptions[virtualRow.index].value
									}"
									@mouseenter="loadFontPreview(fontOptions[virtualRow.index])"
									@click="selectFont(fontOptions[virtualRow.index].value)"
								>
									<Check
										:class="
											cn(
												'h-4 w-4 shrink-0',
												selectedFont === fontOptions[virtualRow.index].value ? 'opacity-100' : 'opacity-0'
											)
										"
									/>
									<span
										class="truncate text-sm"
										:style="{
											fontFamily: `'${getFontFamily(fontOptions[virtualRow.index])}'`,
										}"
									>
										{{ fontOptions[virtualRow.index].label }}
									</span>
								</div>
							</div>
						</div>
					</div>
				</div>
			</UiPopoverContent>
		</UiPopover>
	</div>
</template>
