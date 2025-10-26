import { useOapi } from '~/composables/use-oapi'

export default defineEventHandler(async (event) => {
	const api = useOapi({ headers: event.node.req.headers })
	const id = getRouterParam(event, 'id')
	if (!id) {
		throw createError({
			statusCode: 404,
			statusMessage: 'Not found',
		})
	}

	const req = await api.v1.pastebinGetById(id as string)
	if (req.error) {
		throw createError({
			statusCode: 404,
			statusMessage: req.error,
		})
	}
	if (!req.data?.data?.content) {
		throw createError({
			statusCode: 500,
			statusMessage: 'Internal error',
		})
	}

	setResponseHeader(event, 'Content-Type', 'text/plain; charset=utf-8')
	return req.data.data.content
})
