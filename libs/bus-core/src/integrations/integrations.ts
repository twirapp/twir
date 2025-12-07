export const AddIntegrationTopic = 'integrations.add'
export const RemoveIntegrationTopic = 'integrations.remove'

export enum IntegrationService {
	DONATIONALERTS = 'DONATIONALERTS',
	STREAMLABS = 'STREAMLABS',
	DONATEPAY = 'DONATEPAY',
	FACEIT = 'FACEIT',
}

export interface Request {
	id: string
	service: IntegrationService
}
