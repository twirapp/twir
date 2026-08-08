/**
 * Humanized remaining time until the given date, e.g. "in 29 days" / "через 29 дней".
 * Returns null when the date is invalid or already in the past (caller falls back
 * to the absolute date).
 */
export function formatTimeUntil(value: string, locale: string): string | null {
	const date = new Date(value)
	if (Number.isNaN(date.getTime())) return null

	const diffMs = date.getTime() - Date.now()
	if (diffMs <= 0) return null

	const minutes = diffMs / 60_000
	const hours = diffMs / 3_600_000
	const days = diffMs / 86_400_000

	const rtf = new Intl.RelativeTimeFormat(locale, { numeric: 'auto' })
	if (days >= 1) return rtf.format(Math.round(days), 'day')
	if (hours >= 1) return rtf.format(Math.round(hours), 'hour')
	return rtf.format(Math.max(Math.round(minutes), 1), 'minute')
}
