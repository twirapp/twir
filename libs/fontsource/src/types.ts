export type FontStyle = 'normal' | 'italic'
export type FontSubset = 'latin' | 'cyrillic' | string
export type FontType = 'woff2' | 'woff'

export interface FontItem {
	category: string
	defSubset: string
	family: string
	id: string
	lastModified: string
	styles: FontStyle[]
	subsets: FontSubset[]
	type: string
	variable: boolean
	version: string
	weights: number[]
}

export type FontVariant = {
	[weight: string]: {
		[style: string]: {
			[key in FontSubset]: {
				url: {
					[key in FontType]: string
				}
			}
		}
	}
}

export interface Font {
	id: string
	family: string
	subsets: FontSubset
	weights: number[]
	styles: string[]
	unicodeRange: {
		[key in FontSubset]: string
	}
	defSubset: string
	variable: boolean
	category: string
	version: string
	type: string
	variants: FontVariant
}
