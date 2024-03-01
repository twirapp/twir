import { randomUUID } from 'node:crypto';

import { db } from '../libs/db.js';
import { eventsGrpcClient } from '../libs/eventsGrpc.js';

/**
 * @param {Donate} donate
*/
export const onDonation = async (donate) => {
	const userName = donate.userName ?? 'Anonymous';

	await db.insert({
		id: randomUUID(),
		channel_id: donate.twitchUserId,
		type: 'DONATION',
		data: {
			donationAmount: donate.amount.toString(),
			donationCurrency: donate.currency,
			donationMessage: donate.message,
			donationUsername: userName,
		},
	}).into('channels_events_list');

	const msg = donate.message || '';

	await eventsGrpcClient.donate({
		amount: donate.amount.toString(),
		message: msg,
		currency: donate.currency,
		baseInfo: { channelId: donate.twitchUserId },
		userName,
	});

	// const songs = await ytsrGrpcClient.search({
	// 	search: msg,
	// 	onlyLinks: true,
	// });
	//
	// if (!songs.songs.length) {
	// 	return;
	// }

	// const srCommand = await getSrCommand(donate.twitchUserId);
	// if (!srCommand?.enabled) {
	// 	return;
	// }
	//
	// for (const song of songs.songs) {
	// 	try {
	// 		const parseResult = await parserGrpcClient.processCommand({
	// 			channel: {
	// 				id: donate.twitchUserId,
	// 			},
	// 			message: {
	// 				id: '',
	// 				emotes: [],
	// 				text: `!${srCommand.name} https://youtu.be/${song.id}`,
	// 			},
	// 			sender: {
	// 				id: donate.twitchUserId,
	// 				badges: ['BROADCASTER'],
	// 				name: userName,
	// 				displayName: userName,
	// 			},
	// 		});
	//
	// 		for (const response of parseResult.responses) {
	// 			await botsGrpcClient.sendMessage({
	// 				channelId: donate.twitchUserId,
	// 				message: response,
	// 				skipRateLimits: true,
	// 			});
	// 		}
	// 	} catch (e) {
	// 		console.error(e);
	// 	}
};

/**
	* @param {string} channelId
	* @returns {Promise<{ name: string, enabled: boolean } | null>}
*/
async function getSrCommand(channelId) {
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

	return result;
}
