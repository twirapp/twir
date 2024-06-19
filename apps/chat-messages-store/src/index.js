import { config } from '@twir/config'
import { JSONCodec, connect as natsConnect } from 'nats'
import { createClient } from 'redis'
import { EntityId, Repository, Schema } from 'redis-om'

const redisClient = createClient({
	url: config.REDIS_URL,
})
await redisClient.connect()

const schema = new Schema('chat-messages-store:messages', {
	message_id: { type: 'string' },
	channel_id: { type: 'string' },
	user_id: { type: 'string' },
	user_login: { type: 'string' },
	text: { type: 'text' },
	can_be_deleted: { type: 'boolean' },
	created_at: { type: 'string' },
})

const repository = new Repository(schema, redisClient)
await repository.createIndex()

const sc = JSONCodec()
const nc = await natsConnect({
	servers: config.NATS_URL,
})

const chatMessagesSub = nc.subscribe(
	'chat.messages',
	{
		queue: 'chat-messages-store',
	},
)
const chatMessagesStorePub = nc.subscribe(
	'chat_messages_store.get_by_text_for_delete',
	{ queue: 'chat-messages-store' },
)

const ignoredBadges = ['broadcaster', 'moderator', 'vip', 'subscriber'];

(async () => {
	for await (const m of chatMessagesSub) {
		const data = sc.decode(m.data)

		const canBeDeleted = !data.badges.some(b => ignoredBadges.includes(b.set_id))

		const entity = await repository.save({
			message_id: data.message_id,
			channel_id: data.broadcaster_user_id,
			user_id: data.chatter_user_id,
			user_login: data.chatter_user_login,
			text: data.message.text,
			can_be_deleted: canBeDeleted,
			created_at: new Date().toISOString(),
		})
		await repository.expire(entity[EntityId], 60 * 60)
	}
})();

(async () => {
	for await (const m of chatMessagesStorePub) {
		const data = sc.decode(m.data)
		if (!data.channel_id || !data.text) {
			m.respond(sc.encode({ messages: [] }))
			continue
		}

		const messages = await repository
			.search()
			.where('channel_id').equal(data.channel_id)
			.and('can_be_deleted').equal(true)
			.and('text').match(`*${data.text}*`)
			.return.all()

		m.respond(sc.encode({
			messages,
		}))
	}

	// eslint-disable-next-line style/semi
})();

console.info('[chat-messages-store] is running')
