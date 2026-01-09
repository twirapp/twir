import { defineStore } from 'pinia'

import { graphql } from '~/gql/gql.js'

export const useDashboardSelection = defineStore('dashboard-selection', () => {
	const { executeMutation } = useMutation(
		graphql(`
			mutation SetDashboard($dashboardId: String!) {
				authenticatedUserSelectDashboard(dashboardId: $dashboardId)
			}
		`)
	)

	async function selectDashboard(dashboardId: string) {
		await executeMutation({ dashboardId })
		await navigateTo('/dashboard', { replace: true })

		// Reload to refresh all data with new dashboard context
		if (typeof window !== 'undefined') {
			window.location.reload()
		}
	}

	return { selectDashboard }
})
