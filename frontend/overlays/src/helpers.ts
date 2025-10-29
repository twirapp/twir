import { addZero, randomNum } from '@zero-dependency/utils'

export function generateSocketUrlWithParams(
	path: string,
	params: Record<string, string | undefined>
): string {
	const protocol = location.protocol === 'https:' ? 'wss' : 'ws'
	const url = new URL(`${protocol}://${location.host}/socket${path}`)

	for (const [key, value] of Object.entries(params)) {
		if (!value) continue
		url.searchParams.append(key, value)
	}

	return url.toString()
}

export function base64DecodeUnicode(str: string): string {
	return decodeURIComponent(
		atob(str)
			.split('')
			.map(function (c) {
				return `%${`00${c.charCodeAt(0).toString(16)}`.slice(-2)}`
			})
			.join('')
	)
}

export function getTimeDiffInMilliseconds(minutes: number): number {
	const startDate = new Date()
	const endDate = new Date(startDate.getTime() + minutes * 60 * 1000)
	const diff = endDate.getTime() - startDate.getTime()

	return diff
}

export function millisecondsToTime(ms: number): string {
	const milliseconds = ms % 1000
	ms = (ms - milliseconds) / 1000
	const seconds = ms % 60
	ms = (ms - seconds) / 60
	const minutes = ms % 60
	const hours = (ms - minutes) / 60

	return `${hours ? `${addZero(hours)}:` : ''}${addZero(minutes)}:${addZero(seconds)}`
}

export async function requestWithOutCache<T>(url: string): Promise<T> {
	const res = await fetch(url, { cache: 'no-cache' })
	return await res.json()
}

export function randomRgbColor(): string {
	return `rgb(${randomNum(0, 255)}, ${randomNum(0, 255)}, ${randomNum(0, 255)})`
}

const CHAR_RANGE = {
	emoticons: [0x1F600, 0x1F64F],
	food: [0x1F32D, 0x1F37F],
	animals: [0x1F400, 0x1F4D3],
	expressions: [0x1F910, 0x1F92F],
}

type NamedCharRange = keyof typeof CHAR_RANGE

export function randomEmoji(range: NamedCharRange): string {
	const [max, min] = CHAR_RANGE[range]
	const codePoint = Math.floor(Math.random() * (max - min) + min)
	return String.fromCodePoint(codePoint)
}

export function normalizeDisplayName(displayName: string, userName: string): string {
	if (displayName.toLowerCase() === userName) {
		return displayName
	}

	return userName
}
