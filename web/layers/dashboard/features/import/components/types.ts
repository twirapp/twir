import type { ImportFailureReason } from '~/gql/graphql.js'

export interface ImportFailureData {
	name: string
	reason: ImportFailureReason
}

export interface ImportReportData {
	importedCount: number
	failedCount: number
	failures: ImportFailureData[]
}
