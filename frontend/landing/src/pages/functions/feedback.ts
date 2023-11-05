import { config } from '@twir/config';
import type { Profile } from '@twir/grpc/generated/api/api/auth';
import { APIRoute } from 'astro';
import Redis from 'ioredis';

const { REDIS_URL, DISCORD_FEEDBACK_URL } = config;

const redis = new Redis(REDIS_URL);

const internalError = new Response(JSON.stringify({ error: 'internal error, contact developers in discord' }), { status: 500 });

export type ReviewBody = {
	author: string
	message: string
	profile?: Profile
}

export const post: APIRoute = async ({ request }) => {
	if (!DISCORD_FEEDBACK_URL) {
		console.error('No env setted');
		return internalError;
	}

	const realIp = request.headers.get('x-real-ip');
	if (!realIp) {
		console.error('no real ip');
		return internalError;
	}

	const realIpRedisKey = `landing:feedback-limit:${realIp}`;

	if (await redis.exists(realIpRedisKey)) {
		return new Response(JSON.stringify({ error: 'You already sent an review, please wait 15 minutes.' }), { status: 429 });
	}

	const body: ReviewBody = await new Response(request.body).json();
	if (!body.author || !body.message || body.message.length > 200 || body.author.length > 25) {
		return new Response(JSON.stringify({ error: 'wrong body' }), { status: 400 });
	}

	const embed: Record<string, any> = {
		'type': 'rich',
		'title': `New feedback`,
		'description': body.message,
		'color': 0x00FFFF,
		'author': {
			'name': body.profile?.login ?? body.author,
		},
		thumbnail: {
			url: body.profile?.avatar,
		},
	};

	if (body.profile) {
		embed.author.icon_url = body.profile.avatar;
		embed.thumbnail.url = body.profile.avatar;
	}

	const discordReq = await fetch(DISCORD_FEEDBACK_URL, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify({
			content: 'New feedback.',
			embeds: [embed],
		}),
	});

	if (!discordReq.ok) {
		console.log(await discordReq.text());
		return internalError;
	}

	await redis.set(realIpRedisKey, realIpRedisKey, 'EX', 60 * 15);

	return new Response(JSON.stringify({}), { status: 201 });
};
