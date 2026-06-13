import { useQuery, useSubscription } from '@urql/vue';
import { computed } from 'vue';

import { useMutation } from '@/composables/use-mutation.ts';
import { graphql } from '@/gql';

export const dashboardWidgetsLayoutCacheKey = 'dashboardWidgetsLayout';

const dashboardWidgetsLayoutQuery = graphql(`
	query DashboardWidgetsLayout {
		dashboardWidgetsLayout {
			id
			widgetId
			x
			y
			w
			h
			minW
			minH
			visible
			stackId
			stackOrder
			type
			customName
			customUrl
		}
	}
`);

export function useDashboardWidgetsLayout() {
	const query = useQuery({
		query: dashboardWidgetsLayoutQuery,
	});

	const layout = computed(() => query.data.value?.dashboardWidgetsLayout ?? []);

	return {
		layout,
		fetching: query.fetching,
		error: query.error,
	};
}

const dashboardWidgetsLayoutSubscription = graphql(`
	subscription DashboardWidgetsLayoutChanged {
		dashboardWidgetsLayoutChanged {
			layout {
				id
				widgetId
				x
				y
				w
				h
				minW
				minH
				visible
				stackId
				stackOrder
				type
				customName
				customUrl
			}
		}
	}
`);

export function useDashboardWidgetsLayoutSubscription() {
	const { data, isPaused, fetching } = useSubscription({
		query: dashboardWidgetsLayoutSubscription,
	});

	const layout = computed(() => {
		return data.value?.dashboardWidgetsLayoutChanged.layout ?? [];
	});

	return {
		layout,
		isPaused,
		fetching,
	};
}

export function useDashboardWidgetsLayoutUpdate() {
	return useMutation(
		graphql(`
			mutation DashboardWidgetsLayoutUpdate($input: [DashboardWidgetLayoutInput!]!) {
				dashboardWidgetsLayoutUpdate(input: $input) {
					id
					widgetId
					x
					y
					w
					h
					minW
					minH
					visible
					stackId
					stackOrder
					type
					customName
					customUrl
				}
			}
		`),
		[dashboardWidgetsLayoutCacheKey],
	);
}

export function useDashboardWidgetsCreateCustom() {
	return useMutation(
		graphql(`
			mutation DashboardWidgetsCreateCustom($input: DashboardWidgetCreateCustomInput!) {
				dashboardWidgetsCreateCustom(input: $input) {
					id
					widgetId
					x
					y
					w
					h
					minW
					minH
					visible
					stackId
					stackOrder
					type
					customName
					customUrl
				}
			}
		`),
		[dashboardWidgetsLayoutCacheKey],
	);
}

export function useDashboardWidgetsUpdateCustom() {
	return useMutation(
		graphql(`
			mutation DashboardWidgetsUpdateCustom($input: DashboardWidgetUpdateCustomInput!) {
				dashboardWidgetsUpdateCustom(input: $input) {
					id
					widgetId
					x
					y
					w
					h
					minW
					minH
					visible
					stackId
					stackOrder
					type
					customName
					customUrl
				}
			}
		`),
		[dashboardWidgetsLayoutCacheKey],
	);
}

export function useDashboardWidgetsDelete() {
	return useMutation(
		graphql(`
			mutation DashboardWidgetsDelete($widgetId: String!) {
				dashboardWidgetsDelete(widgetId: $widgetId)
			}
		`),
		[dashboardWidgetsLayoutCacheKey],
	);
}
