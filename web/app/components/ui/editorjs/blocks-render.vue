<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
	data: string | object
}>()

const renderValue = ref<Record<string,any>>()

watch(() => props.data, () => {
	if (!props.data) return
	if (typeof props.data === 'string') {
		renderValue.value = JSON.parse(props.data)
	} else {
		renderValue.value = props.data
	}
}, { immediate: true })
</script>

<template>
	<div v-if="renderValue" class="flex flex-row flex-wrap items-start gap-3">
		<template v-for="block in renderValue.blocks" :key="block.id">
			<p v-if="block.type === 'paragraph'" class="w-full" v-html="block.data.text" />
			<ul v-if="block.type === 'list'" class="ul">
				<li v-for="item of block.data.items" :key="item" v-html="item" />
			</ul>

			<component :is="`h${block.data.level}`" v-if="block.type === 'header'" :class="`h${block.data.level}`" v-html="block.data.text" />

			<blockquote v-if="block.type === 'quote'" class="bq" v-html="block.data.text" />

			<div v-if="block.type === 'delimiter'" class="my-2 w-full border-2 border-b-border" />

			<figure
				v-if="block.type === 'image'"
				class="w-fit max-w-full flex-none overflow-hidden rounded-lg border bg-muted/30 shadow-sm sm:max-w-2xl"
			>
				<a
					:href="block.data.url"
					:aria-label="block.data.caption || 'Open image in a new tab'"
					class="group block"
					target="_blank"
					rel="noopener noreferrer"
				>
					<img
						:src="block.data.url"
						:alt="block.data.caption || ''"
						class="block h-auto max-h-96 w-auto max-w-full object-contain transition-opacity group-hover:opacity-90"
						loading="lazy"
						decoding="async"
					>
				</a>
				<figcaption
					v-if="block.data.caption"
					class="truncate border-t px-3 py-2 text-xs text-muted-foreground"
					:title="block.data.caption"
				>
					{{ block.data.caption }}
				</figcaption>
			</figure>
		</template>
	</div>
</template>

<style scoped>
:deep(a) {
	cursor: pointer;
	color: var(--notification-color);
}

.ul {
	list-style: unset !important;
	margin-left: 1.5em;
	width: 100%;
}

.ul li::marker {
	color: var(--notification-color);
	font-size: 1rem;
}

.h1 {
	font-size: 1.5rem;
	width: 100%;
}

.h2 {
	font-size: 1rem;
	width: 100%;
}

.h3 {
	font-size: 0.875rem;
	width: 100%;
}

.h4 {
	font-size: 0.75rem;
	width: 100%;
}

.h5 {
	font-size: 0.625rem;
	width: 100%;
}

.h6 {
	font-size: 0.5rem;
	width: 100%;
}

.bq {
	border-left: 4px solid var(--notification-color);
	margin: 0;
	padding: 10px 15px;
	background-color: hsl(var(--muted));
	border-radius: 4px;
	width: 100%;
}
</style>
