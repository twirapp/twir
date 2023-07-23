import { config } from '@twir/config';
import Knex from 'knex';

import { Integration } from '../types.js';

export const db = Knex({
	client: 'pg',
	connection: config.DATABASE_URL,
});

export const DONATIONALERTS = 'DONATIONALERTS';
export const STREAMLABS = 'STREAMLABS';
export const DONATEPAY = 'DONATEPAY';

export async function getIntegrations(integrationId: string): Promise<Integration>
export async function getIntegrations(): Promise<Integration[]>
export async function getIntegrations(integrationId?: string): Promise<Integration | Integration[]> {
	let query = db
		.from('channels_integrations')
		.select<Integration[]>([
			'channels_integrations.*',
			db.raw('(json_agg(integration.*) ->> 0)::json as integration'),
			db.raw('(json_agg(channel.*) ->> 0)::json as channel'),
		])
		.where({
			enabled: true,
		})
		.andWhere('integration.service', 'in', [
			DONATIONALERTS,
			STREAMLABS,
			DONATEPAY,
		])
		.andWhere('channel.isEnabled', true)
		.andWhere('channel.isBanned', false)
		.leftJoin('integrations as integration', 'integration.id', '=', 'channels_integrations.integrationId')
		.leftJoin('channels as channel', 'channel.id', '=', 'channels_integrations.channelId')
		.groupBy(['channels_integrations.id', 'integration.id', 'channel.id']);

	if (integrationId) {
		query = query.andWhere('channels_integrations.id', integrationId).first();
	}

	return query;
}
