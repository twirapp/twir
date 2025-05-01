import { useOapi } from '~/composables/use-oapi'

export default defineEventHandler(async (event) => {
	const api = useOapi()
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
	if (!req.data?.content) {
		throw createError({
			statusCode: 500,
			statusMessage: 'Internal error',
		})
	}

	return req.data.content
})
