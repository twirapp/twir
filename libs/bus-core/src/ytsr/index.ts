export const YTSRSearchSubject = 'ytsr.search'

export interface YTSRSearchRequest {
	search: string
	onlyLinks: boolean
}

export interface YTSRSearchResponse {
	songs: YTSRSong[]
}

export interface YTSRSong {
	title: string
	id: string
	views: number
	duration: number
	thumbnailUrl: string | null
	isLive: boolean
	author: YTSRSongAuthor
	link: string | null
	authorName: string
	authorId: string
	authorImage: string
}

export interface YTSRSongAuthor {
	name: string
	channelId: string
	avatarUrl: string | null
}
