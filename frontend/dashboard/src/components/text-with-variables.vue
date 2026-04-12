<script lang="ts" setup>
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

	while (true) {
		const match = variableRegex.exec(text);
		if (match === null) break;
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
					class:
						"mx-0.5 font-mono text-xs inline-flex rounded-md items-center gap-1 !whitespace-normal break-words max-w-full",
				},
				{
					default: () => [
						h(VariableIcon, { class: "size-3 shrink-0" }),
						h("span", { class: "break-words" }, variableData?.description || fullMatch),
					],
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
	<div :class="props.class" class="break-words">
		<component :is="() => renderedContent" />
	</div>
</template>
