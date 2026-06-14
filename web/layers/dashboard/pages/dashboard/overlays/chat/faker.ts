const firstNames = [
	'Alysha',
	'Brian',
	'Bob',
	'Coty',
	'Jon',
	'Sasha',
	'Denis',
]

export const firstName = () => firstNames[Math.floor(Math.random() * firstNames.length)]

const loremText = 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.'
const splittedLorem = loremText.split(' ')

export const lorem = () => splittedLorem.slice(Math.floor(Math.random() * splittedLorem.length)).join(' ')
export const loremWithLen = (len: number) => splittedLorem.slice(Math.floor(Math.random() * splittedLorem.length)).slice(0, len).join(' ')

export const boolean = () => Math.random() < 0.5

const randomBetween = (min: number, max: number) => min + Math.floor(Math.random() * (max - min + 1))

export function rgb() {
	const r = randomBetween(0, 255)
	const g = randomBetween(0, 255)
	const b = randomBetween(0, 255)
	const rgb = `rgb(${r},${g},${b})`
	return rgb
}

export function randomObjectKey(obj: Record<string, unknown>) {
	const keys = Object.keys(obj)

	return keys[Math.floor(Math.random() * keys.length)]
}

export function randomArrayItem<T>(arr: T[]) {
	return arr[Math.floor(Math.random() * arr.length)]
}
