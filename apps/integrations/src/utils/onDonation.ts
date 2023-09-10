import { randomUUID } from 'node:crypto';

import { botsGrpcClient } from '../libs/botsGrpc.js';
import { db } from '../libs/db.js';
import { eventsGrpcClient } from '../libs/eventsGrpc.js';
import { parserGrpcClient } from '../libs/parserGrpc.js';
import { ytsrGrpcClient } from '../libs/ytsrGrpc.js';

export type Donate = {
	twitchUserId: string;
	amount: number | string;
	currency: string;
	message?: string | null;
	userName?: string | null;
}

export const onDonation = async (opts: Donate) => {
	const userName = opts.userName ?? 'Anonymous';

	await db.insert({
		id: randomUUID(),
		channel_id: opts.twitchUserId,
		type: 'DONATION',
		data: {
			donationAmount: opts.amount.toString(),
			donationCurrency: opts.currency,
			donationMessage: opts.message,
			donationUsername: userName,
		},
	}).into('channels_events_list');

	const msg = opts.message || '';

	await eventsGrpcClient.donate({
		amount: opts.amount.toString(),
		message: msg,
		currency: opts.currency,
		baseInfo: { channelId: opts.twitchUserId },
		userName,
	});

	const songs = await ytsrGrpcClient.search({
		search: msg,
		onlyLinks: true,
	});

	if (!songs.songs.length) {
		return;
	}

	const srCommand = await getSrCommand(opts.twitchUserId);
	if (!srCommand?.enabled) {
		return;
	}

	for (const song of songs.songs) {
		try {
			const parseResult = await parserGrpcClient.processCommand({
				channel: {
					id: opts.twitchUserId,
				},
				message: {
					id: '',
					emotes: [],
					text: `!${srCommand.name} https://youtu.be/${song.id}`,
				},
				sender: {
					id: opts.twitchUserId,
					badges: ['BROADCASTER'],
					name: userName,
					displayName: userName,
				},
			});

			for (const response of parseResult.responses) {
				await botsGrpcClient.sendMessage({
					channelId: opts.twitchUserId,
					message: response,
					skipRateLimits: true,
				});
			}
		} catch (e) {
			console.error(e);
		}
	}
};

async function getSrCommand(channelId: string) {
	const result = await db
		.from('channels_commands')
		.select('*')
		.where({
			channelId,
			module: 'SONGS',
			defaultName: 'sr',
		})
		.first();

	if (!result) {
		return null;
	}

	return result as {
		name: string,
		enabled: boolean,
	};
}

// ttsCommand := &model.ChannelsCommands{}
// err = c.db.
// Where(`"channelId" = ?`, msg.Channel.ID).
// Where(`"module" = ?`, "TTS").
// Where(`"defaultName" = ?`, "tts").
// Find(&ttsCommand).
// Error
// if err != nil {
// 	c.logger.Error(
// 		"cannot find tts command",
// 		slog.Any("err", err),
// 		slog.String("channelId", msg.Channel.ID),
// 	)
// 	return
// }
