import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { computed } from 'vue'

import { useMutation } from '@/composables/use-mutation.js'
import { graphql } from '@/gql/gql.js'

const invalidationKey = 'StorageInvalidateKey'

export type StorageEntry = {
	key: string
	value: unknown
	createdAt: string
	updatedAt: string
}

export const useStorageApi = createGlobalState(() => {
	const storageQuery = useQuery({
		variables: {},
		context: { additionalTypenames: [invalidationKey] },
		query: graphql(`
			query GetStorage {
				storage {
					entries {
						key
						value
						createdAt
						updatedAt
					}
					totalSize
				}
			}
		`),
	})

	const entries = computed<StorageEntry[]>(() => {
		return (storageQuery.data.value?.storage?.entries ?? []) as StorageEntry[]
	})

	const isLoading = computed(() => {
		return storageQuery.fetching.value
	})

	const totalSize = computed(() => {
		return storageQuery.data.value?.storage?.totalSize ?? 0
	})

	const useMutationStorageSet = () =>
		useMutation(
			graphql(`
				mutation StorageSet($key: String!, $value: JSON!) {
					storageSet(key: $key, value: $value) {
						key
						value
						createdAt
						updatedAt
					}
				}
			`),
			[invalidationKey],
		)

	const useMutationStorageDelete = () =>
		useMutation(
			graphql(`
				mutation StorageDelete($key: String!) {
					storageDelete(key: $key)
				}
			`),
			[invalidationKey],
		)

	const useMutationStorageDeleteAll = () =>
		useMutation(
			graphql(`
				mutation StorageDeleteAll {
					storageDeleteAll
				}
			`),
			[invalidationKey],
		)

	return {
		storageQuery,
		entries,
		isLoading,
		totalSize,
		useMutationStorageSet,
		useMutationStorageDelete,
		useMutationStorageDeleteAll,
	}
})
