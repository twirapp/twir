import { newBus } from '@twir/bus-core'
import { connect } from 'nats'

const nc = await connect({
	servers: Bun.env.NODE_ENV === 'production'
		? 'nats://nats:4222'
		: 'nats://127.0.0.1:4222',
})

export const twirBus = newBus(nc)
