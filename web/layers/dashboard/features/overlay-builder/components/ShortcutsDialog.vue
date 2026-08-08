<script setup lang="ts">
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
} from '@/components/ui/dialog'

const open = defineModel<boolean>('open', { default: false })

const { t } = useI18n()

interface ShortcutItem {
	combos: string[][]
	labelKey: string
}

const shortcuts: ShortcutItem[] = [
	{ combos: [['Alt', 'Click']], labelKey: 'overlayBuilder.shortcuts.deepSelect' },
	{ combos: [['Esc']], labelKey: 'overlayBuilder.shortcuts.clearSelection' },
	{ combos: [['Ctrl', 'S']], labelKey: 'overlayBuilder.toolbar.save' },
	{ combos: [['Ctrl', 'Z']], labelKey: 'overlayBuilder.toolbar.undo' },
	{ combos: [['Ctrl', 'Shift', 'Z'], ['Ctrl', 'Y']], labelKey: 'overlayBuilder.toolbar.redo' },
	{ combos: [['Ctrl', 'C']], labelKey: 'overlayBuilder.toolbar.copy' },
	{ combos: [['Ctrl', 'X']], labelKey: 'overlayBuilder.toolbar.cut' },
	{ combos: [['Ctrl', 'V']], labelKey: 'overlayBuilder.shortcuts.paste' },
	{ combos: [['Ctrl', 'D']], labelKey: 'overlayBuilder.toolbar.duplicate' },
	{ combos: [['Ctrl', 'A']], labelKey: 'overlayBuilder.shortcuts.selectAll' },
	{ combos: [['Delete'], ['Backspace']], labelKey: 'overlayBuilder.toolbar.delete' },
	{ combos: [['?']], labelKey: 'overlayBuilder.shortcuts.toggleDialog' },
]
</script>

<template>
	<Dialog v-model:open="open">
		<DialogContent class="sm:max-w-lg">
			<DialogHeader>
				<DialogTitle class="flex items-center gap-2">
					<Icon name="lucide:keyboard" class="h-5 w-5" />
					<span>{{ t('overlayBuilder.shortcuts.title') }}</span>
				</DialogTitle>
				<DialogDescription>
					{{ t('overlayBuilder.shortcuts.description') }}
				</DialogDescription>
			</DialogHeader>

			<div class="flex flex-col gap-2">
				<div
					v-for="shortcut in shortcuts"
					:key="shortcut.labelKey"
					class="flex items-center justify-between gap-4"
				>
					<span class="text-sm text-muted-foreground">{{ t(shortcut.labelKey) }}</span>
					<span class="flex shrink-0 items-center gap-1">
						<template v-for="(combo, comboIndex) in shortcut.combos" :key="comboIndex">
							<span v-if="comboIndex > 0" class="text-xs text-muted-foreground">/</span>
							<template v-for="(key, keyIndex) in combo" :key="keyIndex">
								<span v-if="keyIndex > 0" class="text-xs text-muted-foreground">+</span>
								<kbd class="rounded border bg-muted px-1.5 py-0.5 font-mono text-xs">{{ key }}</kbd>
							</template>
						</template>
					</span>
				</div>
			</div>
		</DialogContent>
	</Dialog>
</template>
