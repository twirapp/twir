import { APIRoute } from 'astro';

const db = new Set<string>();

const internalError = new Response(JSON.stringify({ error: 'internal error, contact developers in discord' }), { status: 500 });

export const post: APIRoute = async ({ request }) => {
	const { DISCORD_FEEDBACK_URL } = import.meta.env;
	if (!DISCORD_FEEDBACK_URL) {
		console.error('No env setted');
		return internalError;
	}

	const realIp = request.headers.get('x-real-ip');
	if (!realIp) {
		console.error('no real ip');
		return internalError;
	}

	if (db.has(realIp)) {
		return new Response(JSON.stringify({ error: 'You already sent an review.' }), { status: 429 });
	}

	const body = await new Response(request.body).json();
	if (!body.author || !body.message || body.message.length > 200 || body.author.length > 25) {
		return new Response(JSON.stringify({ error: 'wrong body' }), { status: 400 });
	}

	const discordReq = await fetch(DISCORD_FEEDBACK_URL, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify({
			content: 'New feedback.',
			embeds: [
				{
					'type': 'rich',
					'title': `New feedback`,
					'description': body.message,
					'color': 0x00FFFF,
					'author': {
						'name': body.author,
					},
				},
			],
		}),
	});

	if (!discordReq.ok) {
		console.log(await discordReq.text());
		return new Response(JSON.stringify({ error: 'cannot send info, this is unexpected' }), { status: 500 });
	}

	db.add(realIp);
	setTimeout(() => {
		db.delete(realIp);
	}, 1 * 60 * 60 * 1000);

	return new Response(JSON.stringify({}), { status: 201 });
};
