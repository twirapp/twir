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
	<div v-if="renderValue" class="flex flex-col">
		<template v-for="block in renderValue.blocks" :key="block.id">
			<p v-if="block.type === 'paragraph'" v-html="block.data.text" />
			<ul v-if="block.type === 'list'" class="ul">
				<li v-for="item of block.data.items" :key="item" v-html="item" />
			</ul>

			<component :is="`h${block.data.level}`" v-if="block.type === 'header'" :class="`h${block.data.level}`" v-html="block.data.text" />

			<blockquote v-if="block.type === 'quote'" class="bq" v-html="block.data.text" />

			<div v-if="block.type === 'delimiter'" class="border-2 border-b-border my-2" />

			<img v-if="block.type === 'image'" :src="block.data.url" />
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
