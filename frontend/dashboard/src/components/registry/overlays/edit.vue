<script setup lang="ts">
import { Copy, Save } from 'lucide-vue-next'
import { NButton, NDivider, NFormItem, NInput, NInputNumber, NModal, useMessage } from 'naive-ui'
import { computed, ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import Moveable from 'vue3-moveable'

import HtmlLayer from './layers/html.vue'
import HtmlLayerForm from './layers/htmlForm.vue'

import type { OnDrag, OnResize } from 'vue3-moveable'

import {
	useChannelOverlayByIdQuery,
	useChannelOverlayCreate,
	useChannelOverlaysQuery,
	useChannelOverlayUpdate,
	useProfile,
} from '@/api/index.js'
import NewSelector from '@/components/registry/overlays/newSelector.vue'
import { copyToClipBoard } from '@/helpers'
import {
	type ChannelOverlay,
	type ChannelOverlayLayerInput,
	ChannelOverlayLayerType,
} from '@/gql/graphql'

const { t } = useI18n()

const route = useRoute()
const overlayId = computed(() => {
	const id = route.params.id
	if (typeof id !== 'string' || id === 'new') {
		return ''
	}

	return id
})

const { data: overlayData } = useChannelOverlayByIdQuery(overlayId)
const overlay = computed(() => overlayData.value?.channelOverlayById)

const { executeQuery: refetchOverlays } = useChannelOverlaysQuery()
const createOverlayMutation = useChannelOverlayCreate()
const updateOverlayMutation = useChannelOverlayUpdate()

type OverlayForm = Omit<
	ChannelOverlay,
	'updatedAt' | 'channelId' | 'createdAt' | '__typename' | 'layers'
> & {
	layers: Array<
		Omit<
			ChannelOverlay['layers'][number],
			'__typename' | 'id' | 'overlayId' | 'createdAt' | 'updatedAt' | 'settings'
		> & {
			settings: Omit<ChannelOverlay['layers'][number]['settings'], '__typename'>
		}
	>
}

const formValue = ref<OverlayForm>({
	id: '',
	name: '',
	layers: [],
	width: 1920,
	height: 1080,
})

watch(overlay, (v) => {
	if (!v) return

	formValue.value.id = v.id
	formValue.value.name = v.name
	formValue.value.layers = v.layers.map((layer) => ({
		type: layer.type,
		posX: layer.posX,
		posY: layer.posY,
		width: layer.width,
		height: layer.height,
		periodicallyRefetchData: layer.periodicallyRefetchData,
		settings: {
			htmlOverlayHtml: layer.settings.htmlOverlayHtml,
			htmlOverlayCss: layer.settings.htmlOverlayCss,
			htmlOverlayJs: layer.settings.htmlOverlayJs,
			htmlOverlayDataPollSecondsInterval: layer.settings.htmlOverlayDataPollSecondsInterval,
		},
	}))
	formValue.value.width = v.width
	formValue.value.height = v.height
})

const messages = useMessage()

async function save() {
	const data = toRaw(formValue.value)

	if (!data.name || data.name.length > 30) {
		messages.error(t('overlaysRegistry.validations.name'))
		return
	}

	if (!data.layers.length || data.layers.length > 15) {
		messages.error(t('overlaysRegistry.validations.layers'))
		return
	}

	const layersInput: ChannelOverlayLayerInput[] = data.layers.map((layer) => ({
		type: layer.type,
		posX: layer.posX,
		posY: layer.posY,
		width: layer.width,
		height: layer.height,
		periodicallyRefetchData: layer.periodicallyRefetchData,
		settings: {
			htmlOverlayHtml: layer.settings?.htmlOverlayHtml ?? '',
			htmlOverlayCss: layer.settings?.htmlOverlayCss ?? '',
			htmlOverlayJs: layer.settings?.htmlOverlayJs ?? '',
			htmlOverlayDataPollSecondsInterval: layer.settings?.htmlOverlayDataPollSecondsInterval ?? 5,
		},
	}))

	if (data.id) {
		const result = await updateOverlayMutation.executeMutation({
			id: data.id,
			input: {
				name: data.name,
				width: data.width,
				height: data.height,
				layers: layersInput,
			},
		})
		if (result.error) {
			messages.error(result.error.message)
			return
		}
		// Invalidate cache to refresh the overlays list
		refetchOverlays({ requestPolicy: 'network-only' })
	} else {
		const result = await createOverlayMutation.executeMutation({
			input: {
				name: data.name,
				width: data.width,
				height: data.height,
				layers: layersInput,
			},
		})
		if (result.error) {
			messages.error(result.error.message)
			return
		}

		if (result.data?.channelOverlayCreate) {
			const created = result.data.channelOverlayCreate
			formValue.value.id = created.id
			formValue.value.name = created.name
			formValue.value.width = created.width
			formValue.value.height = created.height
			formValue.value.layers = created.layers.map((layer) => ({
				type: layer.type,
				posX: layer.posX,
				posY: layer.posY,
				width: layer.width,
				height: layer.height,
				periodicallyRefetchData: layer.periodicallyRefetchData,
				settings: {
					htmlOverlayHtml: layer.settings.htmlOverlayHtml,
					htmlOverlayCss: layer.settings.htmlOverlayCss,
					htmlOverlayJs: layer.settings.htmlOverlayJs,
					htmlOverlayDataPollSecondsInterval: layer.settings.htmlOverlayDataPollSecondsInterval,
				},
			}))
		}
		// Invalidate cache to refresh the overlays list
		refetchOverlays({ requestPolicy: 'network-only' })
	}

	messages.success(t('sharedTexts.saved'))
}

const currentlyFocused = ref(0)
function focus(index: number) {
	currentlyFocused.value = index
}

interface EventWithLayerIndex {
	index: number
}

function onDrag({ target, transform, index }: OnDrag & EventWithLayerIndex) {
	focus(index)
	target.style.transform = transform
	const [x, y] = transform.match(/(\d+\.\d+|\d+)px/g)!

	formValue.value.layers[index].posX = Number.parseInt(x)
	formValue.value.layers[index].posY = Number.parseInt(y)
}

function onResize({ target, width, height, transform, index }: OnResize & EventWithLayerIndex) {
	focus(index)

	target.style.width = `${width}px`
	target.style.height = `${height}px`
	target.style.transform = transform

	formValue.value.layers[index].height = height
	formValue.value.layers[index].width = width
}

function removeLayer(index: number) {
	formValue.value.layers = formValue.value.layers.filter((_, i) => i !== index)
	focus(-1)
}

const isOverlayNewModalOpened = ref(false)

const { data: profile } = useProfile()
const selectedDashboardTwitchUser = computed(() => {
	return profile.value?.availableDashboards.find((d) => d.id === profile.value?.selectedDashboardId)
})

async function copyUrl(id: string) {
	await copyToClipBoard(
		`${window.location.origin}/overlays/${selectedDashboardTwitchUser.value?.apiKey}/registry/overlays/${id}`
	)
}

const innerWidth = computed(() => window.innerWidth)
</script>

<template>
	<div class="flex max-w-full">
		<div class="w-[85%]">
			<div
				class="container mx-auto"
				:style="{
					width: `${formValue.width}px`,
					height: `${formValue.height}px`,
					transform: `scale(${(innerWidth / formValue.width) * 0.7})`,
				}"
			>
				<div v-for="(layer, index) of formValue.layers" :key="index">
					<HtmlLayer
						v-if="layer.type === ChannelOverlayLayerType.Html"
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
						:target="`#layer-${index}`"
						:draggable="true"
						:resizable="true"
						:rotatable="false"
						:snappable="true"
						:bounds="{ left: 0, top: 0, right: 0, bottom: 0, position: 'css' }"
						:persistData="{
							height: layer.height,
							width: layer.width,
							left: layer.posX,
							top: layer.posY,
						}"
						:origin="false"
						:renderDirections="
							currentlyFocused === index ? ['nw', 'n', 'ne', 'w', 'e', 'sw', 's', 'se'] : []
						"
						@drag="(opts) => onDrag({ ...opts, index })"
						@resize="(opts) => onResize({ ...opts, index })"
						@click="focus(index)"
					>
					</Moveable>
				</div>
			</div>
		</div>
		<div class="flex flex-col gap-1">
			<NButton
				:disabled="!formValue.name || !formValue.layers.length"
				block
				secondary
				type="success"
				@click="save"
			>
				<Save class="size-4" />
				{{ t('sharedButtons.save') }}
			</NButton>
			<NButton block secondary type="info" :disabled="!formValue.id" @click="copyUrl(formValue.id)">
				<Copy class="size-4" />
				{{ t('overlays.copyOverlayLink') }}
			</NButton>

			<NFormItem :label="t('overlaysRegistry.name')">
				<NInput
					v-model:value="formValue.name"
					:placeholder="t('overlaysRegistry.name')"
					:maxlength="30"
				/>
			</NFormItem>

			<NFormItem :label="t('overlaysRegistry.customWidth')">
				<NInputNumber
					v-model:value="formValue.width"
					:min="50"
					:placeholder="t('overlaysRegistry.customWidth')"
				/>
			</NFormItem>

			<NFormItem :label="t('overlaysRegistry.customHeight')">
				<NInputNumber
					v-model:value="formValue.height"
					:min="50"
					:placeholder="t('overlaysRegistry.customHeight')"
				/>
			</NFormItem>

			<NDivider />

			<NButton secondary type="success" @click="isOverlayNewModalOpened = true">
				{{ t('overlaysRegistry.createNewLayer') }}
			</NButton>

			<div class="flex flex-col gap-3 w-full">
				<template v-for="(layer, index) of formValue.layers">
					<HtmlLayerForm
						v-if="layer.type === ChannelOverlayLayerType.Html"
						:key="index"
						v-model:html="formValue.layers[index].settings!.htmlOverlayHtml"
						v-model:css="formValue.layers[index].settings!.htmlOverlayCss"
						v-model:js="formValue.layers[index].settings!.htmlOverlayJs"
						v-model:pollInterval="
							formValue.layers[index].settings!.htmlOverlayDataPollSecondsInterval
						"
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

	<NModal
		v-model:show="isOverlayNewModalOpened"
		class="w-[50vw]"
		preset="card"
		:title="t('sharedButtons.create')"
	>
		<NewSelector
			@select="
				(v) => {
					formValue.layers.push(v)
					isOverlayNewModalOpened = false
				}
			"
		/>
	</NModal>
</template>

<style scoped>
.container {
	background-color: rgb(18, 18, 18);
	transform-origin: 0px 0px;

	background-image:
		linear-gradient(45deg, rgb(34, 34, 34) 25%, transparent 25%),
		linear-gradient(135deg, rgb(34, 34, 34) 25%, transparent 25%),
		linear-gradient(45deg, transparent 75%, rgb(34, 34, 34) 75%),
		linear-gradient(135deg, transparent 75%, rgb(34, 34, 34) 75%);
	background-size: 20px 20px;
	background-position:
		0px 0px,
		10px 0px,
		10px -10px,
		0px 10px;
}
</style>
