<script setup lang="ts">
const { locale, locales, setLocale } = useI18n()

const currentLocaleName = computed(() => {
	return locales.value.find((l) => l.code === locale.value)?.name ?? locale.value
})

const availableLocales = computed(() => {
	return locales.value.filter((l) => l.code !== locale.value)
})
</script>

<template>
	<UiDropdownMenu>
		<UiDropdownMenuTrigger
			class="inline-flex items-center gap-2 rounded-md px-3 py-2 font-medium text-[#ADB0B8] hover:text-[#D5D8DF] transition-colors"
			as="button"
		>
			<Icon name="lucide:globe" class="w-4 h-4 shrink-0" />
			<span class="max-lg:hidden">{{ currentLocaleName }}</span>
		</UiDropdownMenuTrigger>

		<UiDropdownMenuContent align="end" side="bottom" :side-offset="4" class="min-w-40">
			<UiDropdownMenuItem
				v-for="availableLocale in availableLocales"
				:key="availableLocale.code"
				as="button"
				class="flex w-full items-center"
				@click="setLocale(availableLocale.code)"
			>
				{{ availableLocale.name ?? availableLocale.code }}
			</UiDropdownMenuItem>
		</UiDropdownMenuContent>
	</UiDropdownMenu>
</template>
