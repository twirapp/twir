<script setup lang="ts">
import { VariableIcon } from "lucide-vue-next";
import { computed, h } from "vue";

import { useVariablesApi } from "@/api/variables";
import { Badge } from "@/components/ui/badge";

import type { HTMLAttributes, VNode } from "vue";

const props = defineProps<{
	text: string;
	class?: HTMLAttributes["class"];
}>();

const { allVariables } = useVariablesApi();

// Создаем Map для быстрого поиска переменных
const variablesMap = computed(() => {
	const map = new Map<string, { description: string | null; example: string }>();

	for (const variable of allVariables.value) {
		map.set(variable.example, {
			description: variable.description ?? null,
			example: variable.example,
		});
	}

	return map;
});

// Функция для парсинга текста и создания элементов
function parseTextWithVariables(text: string): VNode[] {
	const result: VNode[] = [];
	// Регулярное выражение для поиска переменных вида $(...)
	const variableRegex = /\$\(([^)]+)\)/g;
	let lastIndex = 0;
	let match: RegExpExecArray | null;

	while ((match = variableRegex.exec(text)) !== null) {
		// Добавляем текст до переменной
		if (match.index > lastIndex) {
			result.push(h("span", {}, text.slice(lastIndex, match.index)));
		}

		const variableContent = match[1]; // содержимое внутри $()
		const fullMatch = match[0]; // полное совпадение $()
		const variableData = variablesMap.value.get(variableContent);

		// Создаем бейдж
		result.push(
			h(
				Badge,
				{
					variant: "default",
					class: "mx-0.5 p-0.5 font-mono text-xs rounded inline-flex align-center",
				},
				{
					default: () => [h(VariableIcon), variableData?.description || fullMatch],
				},
			),
		);

		lastIndex = variableRegex.lastIndex;
	}

	// Добавляем оставшийся текст после последней переменной
	if (lastIndex < text.length) {
		result.push(h("span", {}, text.slice(lastIndex)));
	}

	return result;
}

const renderedContent = computed(() => {
	return parseTextWithVariables(props.text);
});
</script>

<template>
	<div :class="props.class">
		<component :is="() => renderedContent" />
	</div>
</template>
