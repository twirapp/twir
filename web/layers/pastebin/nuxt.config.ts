export default defineNuxtConfig({
	serverHandlers: ['/h', '/h/documents'].map((route) => ({
		route,
		handler: '~~/layers/pastebin/server/routes/h/documents.post.ts',
		method: 'post',
	})),
})
