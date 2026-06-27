import { createColumnHelper, getCoreRowModel, useVueTable } from "@tanstack/vue-table";
import { computed } from "vue";

import type { Secret } from "./use-secrets-api";

import { useSecretsApi } from "./use-secrets-api";

import SecretsActions from "../ui/secret-actions.vue";

const columnHelper = createColumnHelper<Secret>();

export function useSecretsTable() {
	const secretsApi = useSecretsApi();

	const { data, fetching } = secretsApi.secretsQuery;

	const secrets = computed<Secret[]>(() => {
		if (!data.value?.secrets) return [];
		return data.value.secrets;
	});

	const columns = [
		columnHelper.accessor("name", {
			header: () => "Name",
		}),
		columnHelper.accessor("description", {
			header: () => "Description",
			cell: (info) => info.getValue() ?? "-",
		}),
		columnHelper.display({
			id: "actions",
			cell: (info) => h(SecretsActions, { secret: info.row.original }),
		}),
	];

	const table = useVueTable({
		get data() {
			return secrets.value;
		},
		columns,
		getCoreRowModel: getCoreRowModel(),
	});

	return {
		table,
		isLoading: fetching,
	};
}
