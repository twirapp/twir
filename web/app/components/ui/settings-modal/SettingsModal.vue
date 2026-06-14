<script setup lang="ts">
import { IconSettings } from '@tabler/icons-vue'
import { clsx } from 'clsx'
import { useI18n } from 'vue-i18n'

import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'

defineProps<{
	open: boolean
	title: string
	showSave?: boolean
	saveDisabled?: boolean
	buttonDisabled?: boolean
	contentClass?: string
}>()

defineEmits<{
	'update:open': [value: boolean]
	'save': []
}>()

const { t } = useI18n()
</script>

<template>
	<Dialog :open="open" @update:open="$emit('update:open', $event)">
		<DialogTrigger asChild>
			<Button
				variant="secondary"
				:disabled="buttonDisabled"
			>
				<div class="flex items-center gap-1">
					<span>{{ t('sharedButtons.settings') }}</span>
					<IconSettings class="h-4 w-4" />
				</div>
			</Button>
		</DialogTrigger>
		<DialogContent :class="clsx('sm:max-w-[80vw] sm:max-h-[80vh] overflow-y-auto', contentClass)">
			<DialogHeader>
				<DialogTitle>{{ title }}</DialogTitle>
			</DialogHeader>

			<div class="py-4">
				<slot />
			</div>

			<DialogFooter>
				<Button variant="outline" @click="$emit('update:open', false)">
					{{ t('sharedButtons.close') }}
				</Button>
				<Button
					v-if="showSave"
					:disabled="saveDisabled"
					@click="$emit('save')"
				>
					{{ t('sharedButtons.save') }}
				</Button>
			</DialogFooter>
		</DialogContent>
	</Dialog>
</template>
