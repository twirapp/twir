export function padTo2Digits(num: number) {
	return num.toString().padStart(2, '0')
}

export function convertMillisToTime(millis: number) {
	let seconds = Math.floor(millis / 1000)
	let minutes = Math.floor(seconds / 60)
	const hours = Math.floor(minutes / 60)

	seconds = seconds % 60
	minutes = minutes % 60

	return `${hours ? `${padTo2Digits(hours)}:` : ''}${padTo2Digits(minutes)}:${padTo2Digits(
		seconds,
	)}`
}
