<script setup lang="ts">
import { IconTrash, IconSettings } from '@tabler/icons-vue';
import { IconCopy } from '@tabler/icons-vue';
import { IconDeviceFloppy } from '@tabler/icons-vue';
import { type Overlay } from '@twir/grpc/generated/api/api/overlays';
import { NInput, NFormItem, NButton, NCard, NDivider, useThemeVars, NInputNumber } from 'naive-ui';
import { computed, ref, toRaw, watch } from 'vue';
import { useRoute } from 'vue-router';
import Moveable from 'vue3-moveable';
import type { OnDrag, OnResize } from 'vue3-moveable';

import { convertOverlayLayerTypeToText } from './helpers.js';

import { useOverlaysRegistry } from '@/api/index.js';

const route = useRoute();
const overlayId = computed(() => {
	const id = route.params.id;
	if (typeof id !== 'string' || id === 'new') {
		return '';
	}

	return id;
});

const overlaysManager = useOverlaysRegistry();
const creator = overlaysManager.create;
const updater = overlaysManager.update!;
const { data: overlay, refetch } = overlaysManager.getOne!({
	id: overlayId.value,
	isQueryDisabled: !!overlayId.value,
});

watch(overlayId, (v) => {
	if (!v) return;
	refetch();
}, { immediate: true });

type OverlayForm = Omit<Overlay, 'updatedAt' | 'channelId' | 'createdAt'>

const formValue = ref<OverlayForm>({
	id: '',
	name: '',
	layers: [],
	width: 1920,
	height: 1080,
});

watch(overlay, (v) => {
	if (!v) return;

	const raw = toRaw(v);
	console.log(raw);

	formValue.value.id = raw.id;
	formValue.value.name = raw.name;
	formValue.value.layers = raw.layers;
	formValue.value.width = raw.width;
	formValue.value.height = raw.height;
});

async function save() {
	const data = toRaw(formValue.value);

	if (data.id) {
		await updater.mutateAsync({
			...data,
			id: data.id,
		});
	} else {
		await creator.mutateAsync(data);
	}
}

const theme = useThemeVars();
const activeLayourCardColor = computed(() => theme.value.infoColor);

const currentlyFocused = ref(0);
const focus = (index: number) => {
	currentlyFocused.value = index;
};

type EventWithLayerIndex = { index: number }
function onDrag({ target, transform, index }: OnDrag & EventWithLayerIndex) {
	focus(index);
	target.style.transform = transform;
	const [x, y] = transform.match(/(\d+\.\d+|\d+)px/g)!;

	formValue.value.layers[index].posX = parseInt(x);
	formValue.value.layers[index].posY = parseInt(y);
}
function onResize({ target, width, height, transform, index }: OnResize & EventWithLayerIndex) {
	focus(index);

	target.style.width = `${width}px`;
	target.style.height = `${height}px`;
	target.style.transform = transform;

	formValue.value.layers[index].height = height;
	formValue.value.layers[index].width = width;
}
</script>

<template>
	<div style="display: flex; max-width: 100%;">
		<div style="width: 85%;">
			<div class="container" :style="{ width: `${formValue.width}px`, height: `${formValue.height}px` }">
				<div v-for="(layer, index) of formValue.layers" :key="index">
					<div
						:id="'layer-' + index"
						style="position: absolute;"
						:style="{
							transform: `translate(${layer.posX}px, ${layer.posY}px)`,
							width: `${layer.width}px`,
							height: `${layer.height}px`
						}"
					>
						qwe
					</div>
					<Moveable
						className="moveable"
						:target="'#layer-' + index"
						:draggable="true"
						:resizable="true"
						:rotatable="false"
						:snappable="true"
						:bounds="{ left: 0, top: 0, right: 0, bottom: 0, position: 'css' }"
						:persistData="({
							height: layer.height,
							width: layer.width,
							left: layer.posX,
							top: layer.posY
						})"
						:origin="false"
						:renderDirections="currentlyFocused === index ? ['nw', 'n', 'ne', 'w', 'e', 'sw', 's', 'se'] : []"
						@drag="(opts) => onDrag({ ...opts, index })"
						@resize="(opts) => onResize({ ...opts, index })"
						@click="focus(index)"
					>
						<!-- @render="(e) => {
							e.target.style.cssText += e.cssText
							e.moveable.getRect();
						}" -->
					</moveable>
				</div>
			</div>
		</div>
		<div style="display: flex; gap: 4px; flex-direction: column;">
			<n-button block secondary type="success" @click="save">
				<IconDeviceFloppy /> Save
			</n-button>
			<n-button block secondary type="info">
				<IconCopy /> Copy link
			</n-button>
			<n-button block secondary type="error">
				<IconTrash /> Delete
			</n-button>

			<n-form-item label="Name">
				<n-input v-model:value="formValue.name" placeholder="Overlay name" />
			</n-form-item>

			<n-form-item label="Custom width">
				<n-input-number v-model:value="formValue.width" :min="50" placeholder="Custom width" />
			</n-form-item>

			<n-form-item label="Custom height">
				<n-input-number v-model:value="formValue.height" :min="50" placeholder="Custom height" />
			</n-form-item>

			<n-divider />

			<div style="display: flex; flex-direction: column; gap: 12px; width: 100%">
				<n-card
					v-for="(layer, index) of formValue.layers" :key="index"
					:title="convertOverlayLayerTypeToText(layer.type)"
					style="cursor: pointer;"
					:style="{
						border: currentlyFocused === index ? `1px solid ${activeLayourCardColor}` : undefined
					}"
					@click="focus(index)"
				>
					<div style="display: flex; gap: 12px; width: 100%">
						<n-button
							style="flex: 1" secondary
							@click="(e) => {
								e.stopPropagation();
							}"
						>
							<IconSettings />
						</n-button>
						<n-button
							style="flex: 1" secondary type="error"
							@click="(e) => {
								e.stopPropagation();
								formValue.layers = formValue.layers.filter((_, i) => i != index)
							}"
						>
							<IconTrash />
						</n-button>
					</div>
				</n-card>
			</div>
		</div>
	</div>
</template>

<style scoped>
.container {
	background-color: rgb(18, 18, 18);
	transform-origin: 0px 0px;
	transform: scale(0.7);
	background-image: linear-gradient(45deg, rgb(34, 34, 34) 25%, transparent 25%), linear-gradient(135deg, rgb(34, 34, 34) 25%, transparent 25%), linear-gradient(45deg, transparent 75%, rgb(34, 34, 34) 75%), linear-gradient(135deg, transparent 75%, rgb(34, 34, 34) 75%);
	background-size: 20px 20px;
	background-position: 0px 0px, 10px 0px, 10px -10px, 0px 10px;
}
</style>
