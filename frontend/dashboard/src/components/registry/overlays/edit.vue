<script setup lang="ts">
import { IconCopy } from '@tabler/icons-vue';
import { IconDeviceFloppy } from '@tabler/icons-vue';
import { OverlayLayerType, type Overlay } from '@twir/api/messages/overlays/overlays';
import { NInput, NFormItem, NButton, NDivider, NInputNumber, NModal, useMessage } from 'naive-ui';
import { computed, ref, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute } from 'vue-router';
import Moveable from 'vue3-moveable';
import type { OnDrag, OnResize } from 'vue3-moveable';

import HtmlLayer from './layers/html.vue';
import HtmlLayerForm from './layers/htmlForm.vue';

import {
	useOverlaysRegistry, useProfile,
} from '@/api/index.js';
import NewSelector from '@/components/registry/overlays/newSelector.vue';
import { copyToClipBoard } from '@/helpers';

const { t } = useI18n();

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
	isQueryDisabled: true,
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

	formValue.value.id = raw.id;
	formValue.value.name = raw.name;
	formValue.value.layers = raw.layers;
	formValue.value.width = raw.width;
	formValue.value.height = raw.height;
});

const messages = useMessage();

async function save() {
	const data = toRaw(formValue.value);

	if (!data.name || data.name.length > 30) {
		messages.error(t('overlaysRegistry.validations.name'));
		return;
	}

	if (!data.layers.length || data.layers.length > 15) {
		messages.error(t('overlaysRegistry.validations.layers'));
		return;
	}

	if (data.id) {
		await updater.mutateAsync({
			...data,
			id: data.id,
		});
	} else {
		const newOverlayData = await creator.mutateAsync(data);

		const raw = toRaw(newOverlayData);

		formValue.value.id = raw.id;
		formValue.value.name = raw.name;
		formValue.value.layers = raw.layers;
		formValue.value.width = raw.width;
		formValue.value.height = raw.height;
	}

	messages.success(t('sharedTexts.saved'));
}

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

const removeLayer = (index: number) => {
	formValue.value.layers = formValue.value.layers.filter((_, i) => i != index);
	focus(-1);
};

const isOverlayNewModalOpened = ref(false);

const userProfile = useProfile();
const copyUrl = async (id: string) => {
	await copyToClipBoard(`${window.location.origin}/overlays/${userProfile.data.value?.apiKey}/registry/overlays/${id}`);
};

const innerWidth = computed(() => window.innerWidth);
</script>

<template>
	<div style="display: flex; max-width: 100%;">
		<div style="width: 85%;">
			<div
				class="container"
				:style="{
					width: `${formValue.width}px`,
					height: `${formValue.height}px`,
					transform: `scale(${(innerWidth / formValue.width) * 0.7})`
				}"
			>
				<div v-for="(layer, index) of formValue.layers" :key="index">
					<HtmlLayer
						v-if="layer.type === OverlayLayerType.HTML"
						:posX="layer.posX"
						:posY="layer.posY"
						:width="layer.width"
						:height="layer.height"
						:index="index"
						:text="layer.settings?.htmlOverlayHtml ?? ''"
						:css="layer.settings?.htmlOverlayCss ?? ''"
						:js="layer.settings?.htmlOverlayJs ?? ''"
						:periodicallyRefetchData="layer.periodicallyRefetchData"
					/>

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
					</moveable>
				</div>
			</div>
		</div>
		<div style="display: flex; gap: 4px; flex-direction: column;">
			<n-button
				:disabled="!formValue.name || !formValue.layers.length" block secondary
				type="success" @click="save"
			>
				<IconDeviceFloppy />
				{{ t('sharedButtons.save') }}
			</n-button>
			<n-button
				block secondary type="info" :disabled="!formValue.id"
				@click="copyUrl(formValue.id)"
			>
				<IconCopy />
				{{ t('overlays.copyOverlayLink') }}
			</n-button>

			<n-form-item :label="t('overlaysRegistry.name')">
				<n-input
					v-model:value="formValue.name" :placeholder="t('overlaysRegistry.name')"
					:maxlength="30"
				/>
			</n-form-item>

			<n-form-item :label="t('overlaysRegistry.customWidth')">
				<n-input-number
					v-model:value="formValue.width" :min="50"
					:placeholder="t('overlaysRegistry.customWidth')"
				/>
			</n-form-item>

			<n-form-item :label="t('overlaysRegistry.customHeight')">
				<n-input-number
					v-model:value="formValue.height" :min="50"
					:placeholder="t('overlaysRegistry.customHeight')"
				/>
			</n-form-item>

			<n-divider />

			<n-button
				secondary
				type="success"
				@click="isOverlayNewModalOpened = true"
			>
				{{ t('overlaysRegistry.createNewLayer') }}
			</n-button>

			<div style="display: flex; flex-direction: column; gap: 12px; width: 100%">
				<template v-for="(layer, index) of formValue.layers">
					<html-layer-form
						v-if="layer.type === OverlayLayerType.HTML"
						:key="index"
						v-model:html="formValue.layers[index].settings!.htmlOverlayHtml"
						v-model:css="formValue.layers[index].settings!.htmlOverlayCss"
						v-model:js="formValue.layers[index].settings!.htmlOverlayJs"
						v-model:pollInterval="formValue.layers[index].settings!.htmlOverlayHtmlDataPollSecondsInterval"
						v-model:periodicallyRefetchData="formValue.layers[index].periodicallyRefetchData"
						:isFocused="currentlyFocused === index"
						:layerIndex="index"
						:type="layer.type"
						@remove="removeLayer"
						@focus="focus"
					/>
				</template>
			</div>
		</div>
	</div>

	<n-modal
		v-model:show="isOverlayNewModalOpened" style="width: 50vw" preset="card"
		:title="t('sharedButtons.create')"
	>
		<new-selector
			@select="v => {
				formValue.layers.push(v)
				isOverlayNewModalOpened = false
			}"
		/>
	</n-modal>
</template>

<style scoped>
.container {
	background-color: rgb(18, 18, 18);
	transform-origin: 0px 0px;

	background-image: linear-gradient(45deg, rgb(34, 34, 34) 25%, transparent 25%), linear-gradient(135deg, rgb(34, 34, 34) 25%, transparent 25%), linear-gradient(45deg, transparent 75%, rgb(34, 34, 34) 75%), linear-gradient(135deg, transparent 75%, rgb(34, 34, 34) 75%);
	background-size: 20px 20px;
	background-position: 0px 0px, 10px 0px, 10px -10px, 0px 10px;
}
</style>
