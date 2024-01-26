import { RpcError } from '@protobuf-ts/runtime-rpc';
import type {
	GetDataResponse,
	UpdateDataRequest,
} from '@twir/api/messages/integrations_seventv/integrations_seventv';
import { defineStore } from 'pinia';
import { computed, MaybeRef, ref, unref, watch } from 'vue';

import { useSevenTvIntegration } from '@/api/integrations/seventv.js';

export const useSevenTv = defineStore('sevenTvIntegration', () => {
	const data = ref<GetDataResponse>();
	const error = ref<unknown | null>(null);

	const manager = useSevenTvIntegration();
	const updater = manager.useUpdate();
	const { data: sevenTvData, error: sevenTvError } = manager.useData();

	const sevenTvProfileLink = computed(() => {
		return `https://7tv.app/users/${data?.value?.userSeventvProfile?.id}`;
	});

	watch(sevenTvData, (value) => {
		data.value = value;
	});

	watch(sevenTvError, (value) => {
		error.value = value;
	});

	const isNotRegistered = computed(() => {
		if (error.value && error.value instanceof RpcError) {
			return error.value.message === 'profile_not_found';
		}

		return false;
	});

	async function save(data: MaybeRef<UpdateDataRequest>) {
		const raw = unref(data);
		await updater.mutateAsync({
			rewardIdForAddEmote: raw.rewardIdForAddEmote ?? undefined,
			rewardIdForRemoveEmote: raw.rewardIdForRemoveEmote ?? undefined,
			deleteEmotesOnlyAddedByApp: raw.deleteEmotesOnlyAddedByApp,
		});
	}

	return {
		data,
		error,
		sevenTvProfileLink,
		isNotRegistered,
		save,
	};
});
