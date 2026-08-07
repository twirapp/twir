interface AnchorLink {
	href: string
	label: string
	header?: boolean
}

export const anchorLinks: AnchorLink[] = [
	{
		href: '/#features',
		label: 'landing.nav.features',
	},
	{
		href: '/#integrations',
		label: 'landing.nav.integrations',
	},
	{
		href: '/compare',
		label: 'landing.nav.compare',
	},
	{
		href: '/url-shortener',
		label: 'landing.nav.urlShortener',
	},
	{
		href: '/h',
		label: 'landing.nav.hastebin',
	},
	{
		href: '/terms',
		label: 'landing.nav.terms',
		header: false,
	},
	{
		href: '/privacy',
		label: 'landing.nav.privacy',
		header: false,
	},
]

export const headerLinks = anchorLinks.filter((link) => link.header !== false)
