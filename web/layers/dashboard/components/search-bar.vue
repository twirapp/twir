<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import { Input } from '@/components/ui/input'

defineProps<{ placeholder?: string }>()

const model = defineModel<string>({ default: '' })

const { t } = useI18n()
</script>

<template>
	<div class="relative w-full items-center">
		<Input v-model="model" type="text" :placeholder="placeholder ?? t('sharedTexts.searchPlaceholder')" class="h-9 pl-10" />
		<span class="absolute start-2 inset-y-0 flex items-center justify-center px-2">
			<Icon name="lucide:search" class="size-4 text-muted-foreground" />
		</span>
	</div>
</template>
