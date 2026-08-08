export function formatBytes(bytes: number, decimals = 1): string {
	if (!Number.isFinite(bytes) || bytes <= 0) return '0 B'

	const step = 1024
	const units = ['B', 'KB', 'MB', 'GB', 'TB']
	const unitIndex = Math.min(Math.floor(Math.log(bytes) / Math.log(step)), units.length - 1)
	const value = bytes / step ** unitIndex

	return `${value.toFixed(unitIndex === 0 ? 0 : decimals)} ${units[unitIndex]}`
}
