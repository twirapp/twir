type DonatePayEvent = {
	data: {
		notification: {
			type: 'donation',
			vars: {
				name: string,
				comment: string,
				sum: number,
				currency: 'string'
			}
		}
	}
}
