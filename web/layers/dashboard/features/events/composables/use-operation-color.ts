import { EventOperationType } from '~/gql/graphql.ts'
import { flatOperations } from '~/features/events/constants/helpers.ts'

export function getOperationColor(operationType?: EventOperationType): string {
	if (!operationType) {
		return 'bg-gray-500'
	}

	const operation = flatOperations[operationType]

	switch (operation?.color) {
		case 'default':
		case undefined:
			return 'bg-gray-500'
		case 'info':
			return 'bg-blue-500'
		case 'success':
			return 'bg-green-500'
		case 'error':
			return 'bg-red-500'
		case 'warning':
			return 'bg-yellow-500'
		default:
			return 'bg-gray-500'
	}
}
