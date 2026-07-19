export type CompareSupport = 'yes' | 'partial' | 'no'

export type CompareBotId = 'twir' | 'nightbot' | 'streamelements' | 'moobot' | 'fossabot'

export interface CompareBot {
	id: CompareBotId
	name: string
	siteUrl: string
	logoType: 'icon' | 'img'
	logo: string
	logoClass: string
}

export interface CompareCell {
	support: CompareSupport
	noteKey?: string
}

export interface CompareFeatureRow {
	labelKey: string
	cells: Record<CompareBotId, CompareCell>
}

export interface CompareTextRow {
	labelKey: string
	cells: Record<CompareBotId, string>
}

const yes = (noteKey?: string): CompareCell => ({ support: 'yes', noteKey })
const partial = (noteKey?: string): CompareCell => ({ support: 'partial', noteKey })
const no = (noteKey?: string): CompareCell => ({ support: 'no', noteKey })

export const compareBots: CompareBot[] = [
	{
		id: 'twir',
		name: 'Twir',
		siteUrl: 'https://twir.app',
		logoType: 'icon',
		logo: 'twir-logo',
		logoClass: 'h-9 w-9',
	},
	{
		id: 'nightbot',
		name: 'Nightbot',
		siteUrl: 'https://nightbot.tv',
		logoType: 'icon',
		logo: 'twir-compare:nightbot',
		logoClass: 'h-8 w-8 text-white',
	},
	{
		id: 'streamelements',
		name: 'StreamElements',
		siteUrl: 'https://streamelements.com',
		logoType: 'icon',
		logo: 'twir-compare:streamelements',
		logoClass: 'h-8 w-auto max-w-[140px]',
	},
	{
		id: 'moobot',
		name: 'Moobot',
		siteUrl: 'https://moobot.tv',
		logoType: 'img',
		logo: '/compare/moobot.png',
		logoClass: 'h-9 w-auto',
	},
	{
		id: 'fossabot',
		name: 'Fossabot',
		siteUrl: 'https://fossabot.com',
		logoType: 'icon',
		logo: 'twir-compare:fossabot',
		logoClass: 'h-8 w-auto max-w-[130px]',
	},
]

export const compareFeatureRows: CompareFeatureRow[] = [
	{
		labelKey: 'compare.features.commands',
		cells: { twir: yes(), nightbot: yes(), streamelements: yes(), moobot: yes(), fossabot: yes() },
	},
	{
		labelKey: 'compare.features.timers',
		cells: { twir: yes(), nightbot: yes(), streamelements: yes(), moobot: yes(), fossabot: yes() },
	},
	{
		labelKey: 'compare.features.moderation',
		cells: { twir: yes(), nightbot: yes(), streamelements: yes(), moobot: yes(), fossabot: yes() },
	},
	{
		labelKey: 'compare.features.songRequests',
		cells: {
			twir: yes(),
			nightbot: yes(),
			streamelements: yes(),
			moobot: yes(),
			fossabot: yes(),
		},
	},
	{
		labelKey: 'compare.features.giveaways',
		cells: { twir: yes(), nightbot: yes(), streamelements: yes(), moobot: yes(), fossabot: yes() },
	},
	{
		labelKey: 'compare.features.statsTracking',
		cells: {
			twir: yes('compare.notes.statsTwir'),
			nightbot: no(),
			streamelements: yes(),
			moobot: partial('compare.notes.statsMoobot'),
			fossabot: no(),
		},
	},
	{
		labelKey: 'compare.features.greetings',
		cells: {
			twir: yes(),
			nightbot: no(),
			streamelements: no(),
			moobot: partial('compare.notes.greetingsMoobot'),
			fossabot: no(),
		},
	},
	{
		labelKey: 'compare.features.keywords',
		cells: { twir: yes(), nightbot: no(), streamelements: yes(), moobot: no(), fossabot: yes() },
	},
	{
		labelKey: 'compare.features.customScripting',
		cells: {
			twir: yes('compare.notes.scriptingTwir'),
			nightbot: partial(),
			streamelements: partial(),
			moobot: no(),
			fossabot: partial('compare.notes.scriptingFossabot'),
		},
	},
	{
		labelKey: 'compare.features.roles',
		cells: {
			twir: yes(),
			nightbot: no(),
			streamelements: partial('compare.notes.rolesSe'),
			moobot: partial('compare.notes.rolesMoobot'),
			fossabot: yes(),
		},
	},
	{
		labelKey: 'compare.features.vipManagement',
		cells: {
			twir: yes('compare.notes.vipTwir'),
			nightbot: no(),
			streamelements: no(),
			moobot: partial('compare.notes.vipMoobot'),
			fossabot: no(),
		},
	},
	{
		labelKey: 'compare.features.emoteStats',
		cells: { twir: yes(), nightbot: no(), streamelements: no(), moobot: no(), fossabot: no() },
	},
	{
		labelKey: 'compare.features.overlays',
		cells: {
			twir: yes('compare.notes.overlaysTwir'),
			nightbot: no(),
			streamelements: yes(),
			moobot: no(),
			fossabot: no(),
		},
	},
	{
		labelKey: 'compare.features.alerts',
		cells: {
			twir: yes(),
			nightbot: no(),
			streamelements: yes(),
			moobot: partial('compare.notes.alertsMoobot'),
			fossabot: no(),
		},
	},
	{
		labelKey: 'compare.features.games',
		cells: {
			twir: yes('compare.notes.gamesTwir'),
			nightbot: no(),
			streamelements: partial('compare.notes.gamesSe'),
			moobot: no(),
			fossabot: partial('compare.notes.gamesFossabot'),
		},
	},
	{
		labelKey: 'compare.features.musicIntegrations',
		cells: {
			twir: yes('compare.notes.musicTwir'),
			nightbot: partial('compare.notes.musicNightbot'),
			streamelements: partial('compare.notes.musicSe'),
			moobot: no(),
			fossabot: no(),
		},
	},
	{
		labelKey: 'compare.features.musicRecognition',
		cells: { twir: yes(), nightbot: no(), streamelements: no(), moobot: no(), fossabot: no() },
	},
	{
		labelKey: 'compare.features.obsIntegration',
		cells: {
			twir: yes('compare.notes.obsTwir'),
			nightbot: no(),
			streamelements: no(),
			moobot: no(),
			fossabot: no(),
		},
	},
	{
		labelKey: 'compare.features.openSource',
		cells: {
			twir: yes('compare.notes.openSourceTwir'),
			nightbot: no(),
			streamelements: no(),
			moobot: no(),
			fossabot: no(),
		},
	},
]

export const compareTextRows: CompareTextRow[] = [
	{
		labelKey: 'compare.features.platforms',
		cells: {
			twir: 'compare.values.platforms.twir',
			nightbot: 'compare.values.platforms.nightbot',
			streamelements: 'compare.values.platforms.streamelements',
			moobot: 'compare.values.platforms.moobot',
			fossabot: 'compare.values.platforms.fossabot',
		},
	},
	{
		labelKey: 'compare.features.price',
		cells: {
			twir: 'compare.values.price.free',
			nightbot: 'compare.values.price.free',
			streamelements: 'compare.values.price.free',
			moobot: 'compare.values.price.moobotPlus',
			fossabot: 'compare.values.price.free',
		},
	},
]

export const compareFaqKeys = ['q1', 'q2', 'q3', 'q4'] as const
