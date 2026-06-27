import { getSecretsForChannel } from './libs/db'
import { executeCode } from './libs/executor'
import { twirBus } from './libs/twirbus'

console.info('Executron service starting...')

twirBus.Executron.Execute.subscribeGroup('executron', async (data: any) => {
	console.info(`Executing script for channel ${data.channelId}`)

	if (data.language !== 'javascript') {
		return {
			result: '',
			error: `Unsupported language: ${data.language}`,
		}
	}

	const secrets = await getSecretsForChannel(data.channelId)
	const result = await executeCode(data.code, data.channelId, secrets)

	return result
})

console.info('Executron service started')

process.on('uncaughtException', console.error)
process.on('unhandledRejection', console.error)
