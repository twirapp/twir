export const convertBytesToSize = (bytes: number) => {
	const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
	if (bytes === 0) return '0 Byte'
	const i = Math.floor(Math.log(bytes) / Math.log(1024))
	// oxlint-disable-next-line prefer-exponentiation-operator
	return `${(bytes / Math.pow(1024, i)).toFixed(2)} ${sizes[i]}`
}
