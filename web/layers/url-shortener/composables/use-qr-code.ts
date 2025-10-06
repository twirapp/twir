import QRCode from 'qrcode'

interface QRCodeOptions {
	url: string
	color?: string
	backgroundColor?: string
	logoComponent?: HTMLElement | null
	logoSize?: number
	size?: number
}

export function useQRCode() {
	async function generateQRCode(options: QRCodeOptions): Promise<string> {
		const {
			url,
			color = '#000000',
			backgroundColor = '#ffffff',
			logoComponent,
			logoSize = 60,
			size = 300,
		} = options

		const canvas = document.createElement('canvas')
		canvas.width = size
		canvas.height = size
		const ctx = canvas.getContext('2d')

		if (!ctx) throw new Error('Canvas context not available')

		await QRCode.toCanvas(canvas, url, {
			width: size,
			margin: 0,
			color: {
				dark: color,
				light: backgroundColor,
			},
			errorCorrectionLevel: 'H',
		})

		if (logoComponent) {
			const logoData = await componentToImage(logoComponent, logoSize)

			const logoImg = new Image()
			await new Promise((resolve, reject) => {
				logoImg.onload = resolve
				logoImg.onerror = reject
				logoImg.src = logoData
			})

			const logoX = (size - logoSize) / 2
			const logoY = (size - logoSize) / 2

			ctx.fillStyle = backgroundColor
			ctx.beginPath()
			ctx.arc(size / 2, size / 2, logoSize / 2 + 8, 0, Math.PI * 2)
			ctx.fill()

			ctx.save()
			ctx.beginPath()
			ctx.arc(size / 2, size / 2, logoSize / 2, 0, Math.PI * 2)
			ctx.closePath()
			ctx.clip()
			ctx.drawImage(logoImg, logoX, logoY, logoSize, logoSize)
			ctx.restore()
		}

		return canvas.toDataURL('image/png')
	}

	async function componentToImage(element: HTMLElement, size: number): Promise<string> {
		const canvas = document.createElement('canvas')
		canvas.width = size
		canvas.height = size
		const ctx = canvas.getContext('2d')

		if (!ctx) throw new Error('Canvas context not available')

		const svgElement = element.querySelector('svg')
		if (svgElement) {
			const svgData = new XMLSerializer().serializeToString(svgElement)
			const svgBlob = new Blob([svgData], { type: 'image/svg+xml;charset=utf-8' })
			const url = URL.createObjectURL(svgBlob)

			const img = new Image()
			await new Promise((resolve, reject) => {
				img.onload = resolve
				img.onerror = reject
				img.src = url
			})

			ctx.drawImage(img, 0, 0, size, size)
			URL.revokeObjectURL(url)
		}

		return canvas.toDataURL('image/png')
	}

	return {
		generateQRCode,
	}
}
