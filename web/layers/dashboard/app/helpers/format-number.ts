const intl = new Intl.NumberFormat(navigator.language);

export function formatNumber(num: number): string {
	return intl.format(num);
}
