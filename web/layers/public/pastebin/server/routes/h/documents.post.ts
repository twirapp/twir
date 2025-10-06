import { useOapi } from '~/composables/use-oapi'

export default defineEventHandler(async (event) => {
	if (event.node.req.method !== 'POST') {
		throw createError({
			statusCode: 405,
			statusMessage: 'Method not allowed',
		})
	}
	const url = getRequestURL(event)

	const api = useOapi({ headers: event.node.req.headers })
	const req = await api.v1.pastebinCreate({
		content: await readBody(event),
	})
	if (req.error) {
		throw createError({
			statusCode: req.status || 500,
			statusMessage: req.error,
		})
	}

	if (!req.data?.id) {
		throw createError({
			statusCode: 500,
			statusMessage: 'No ID returned from api, contact administrator',
		})
	}

	setResponseHeader(event, 'Content-Type', 'application/json; charset=utf-8')
	return {
		// key is for hastebin compatibility format
		key: req.data.id,
		id: req.data.id,
		url: `${url.origin}/h/${req.data.id}`,
		raw_url: `${url.origin}/h/${req.data.id}/raw`,
	}
})
