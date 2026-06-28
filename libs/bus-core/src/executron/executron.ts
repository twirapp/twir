export const ExecuteSubject = 'executron.execute'

export interface ExecuteRequest {
	channelId: string
	language: 'javascript'
	code: string
	userId?: string
}

export interface ExecuteResponse {
	result: string
	error: string
}
