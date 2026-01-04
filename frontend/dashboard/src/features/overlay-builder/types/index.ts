import type { ChannelOverlayLayerType } from '@/gql/graphql'

export interface Layer {
	id: string
	type: ChannelOverlayLayerType
	name: string
	posX: number
	posY: number
	width: number
	height: number
	rotation: number
	opacity: number
	visible: boolean
	locked: boolean
	zIndex: number
	periodicallyRefetchData: boolean
	settings: LayerSettings
}

export interface LayerSettings {
	htmlOverlayHtml?: string
	htmlOverlayCss?: string
	htmlOverlayJs?: string
	htmlOverlayDataPollSecondsInterval?: number
	imageUrl?: string
}

export interface OverlayProject {
	id: string
	name: string
	width: number
	height: number
	layers: Layer[]
}

export interface CanvasState {
	zoom: number
	panX: number
	panY: number
	selectedLayerIds: string[]
	clipboardLayers: Layer[]
	showGrid: boolean
	snapToGrid: boolean
	gridSize: number
	showRulers: boolean
	showGuides: boolean
}

export interface HistoryState {
	past: OverlayProject[]
	present: OverlayProject
	future: OverlayProject[]
}

export interface AlignmentGuide {
	type: 'horizontal' | 'vertical'
	position: number
	matchedLayers: string[]
}

export interface Transform {
	x: number
	y: number
	width: number
	height: number
	rotation: number
}

export enum ToolType {
	SELECT = 'select',
	PAN = 'pan',
}

export interface BuilderAction {
	type: 'add' | 'remove' | 'update' | 'reorder' | 'duplicate'
	layerId?: string
	layerIds?: string[]
	data?: Partial<Layer> | Layer | Layer[]
	previousData?: Partial<Layer> | Layer | Layer[]
}
