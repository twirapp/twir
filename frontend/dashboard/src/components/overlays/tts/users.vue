<script setup lang='ts'>
import { IconSpeakerphone } from '@tabler/icons-vue';
import {
	NAlert,
	NAvatar,
	NButton,
	NCard,
	NCheckbox,
	NGrid,
	NGridItem,
	NInput,
	NPopconfirm,
	NRow,
	NSkeleton,
	NSpace,
	NText,
	useThemeVars,
} from 'naive-ui';
import { UnwrapRef, computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { useTtsOverlayManager, useTwitchGetUsers } from '@/api/index.js';

const themeVars = useThemeVars();
const descriptionColor = computed(() => themeVars.value.textColor3);

const ttsManager = useTtsOverlayManager();
const ttsSettings = ttsManager.getSettings();
const ttsSay = ttsManager.useSay();
const { data: ttsUsersData, isLoading } = ttsManager.getUsersSettings();
const usersSettingsDeleter = ttsManager.deleteUsersSettings();

const usersIdsForRequest = computed(() => {
	if (!ttsUsersData?.value?.data) return [];
	return ttsUsersData.value.data.map((user) => user.userId);
});
const twitchUsers = useTwitchGetUsers({
	ids: usersIdsForRequest,
});

type ListUser = NonNullable<UnwrapRef<typeof ttsUsersData>>['data'][number] & {
	avatar: string,
	name: string,
	markedForDelete: boolean,
}

const mappedUsers = computed<ListUser[]>(() => {
	return ttsUsersData?.value?.data?.map((user) => {
		const twitchUser = twitchUsers?.data?.value?.users.find((twitchUser) => twitchUser.id === user.userId);
		if (!twitchUser) return null;

		return {
			...user,
			avatar: twitchUser.profileImageUrl ?? '',
			name: twitchUser.displayName ?? twitchUser?.login,
			markedForDelete: false,
		};
	}).filter(Boolean) as ListUser[] ?? [];
});
const users = ref<ListUser[]>([]);
watch(mappedUsers, (u) => {
	users.value = u;
});

const isSomeUserMarked = computed(() => {
	return users.value.some(u => u.markedForDelete);
});

function changeMarkedStateForAllUsers(state: boolean) {
	users.value = users.value.map(u => ({
		...u,
		markedForDelete: state,
	}));
}

const testText = ref('');

async function deleteUsers() {
	await usersSettingsDeleter.mutateAsync(users.value.filter(u => u.markedForDelete).map(u => u.userId));
}

async function testUserVoice(user: ListUser) {
	await ttsSay.mutateAsync({
		volume: ttsSettings.data?.value?.data?.volume || 50,
		voice: user.voice,
		pitch: user.pitch,
		rate: user.rate,
		text: testText.value || 'Hello world, привет мир',
	});
}

const { t } = useI18n();
</script>

<template>
	<div class="p-4">
		<n-grid v-if="isLoading || !ttsUsersData" cols="1 s:2 m:2 l:2" responsive="screen" :x-gap="10" :y-gap="10">
			<n-grid-item v-for="i in 16" :key="i" :span="1">
				<n-skeleton v-if="!isLoading" size="large" height="60px" :sharp="false" />
			</n-grid-item>
		</n-grid>

		<div v-else>
			<n-space justify="space-between" class="mb-2.5">
				<n-input v-model:value="testText" placeholder="Text for test user settings" />
				<n-space>
					<n-button
						secondary
						type="info"
						:disabled="!users.length"
						@click="changeMarkedStateForAllUsers(!isSomeUserMarked)"
					>
						{{ t(`overlays.tts.users.${isSomeUserMarked ? 'undoSelection' : 'selectAll'}`) }}
					</n-button>

					<n-popconfirm
						:positive-text="t('deleteConfirmation.confirm')"
						:negative-text="t('deleteConfirmation.cancel')"
						@positive-click="deleteUsers"
					>
						<template #trigger>
							<n-button
								secondary
								type="error"
								:disabled="!users.some(u => u.markedForDelete)"
							>
								{{ t('sharedButtons.delete') }} {{ users.filter(u => u.markedForDelete).length }}
							</n-button>
						</template>

						{{ t('deleteConfirmation.text') }}
					</n-popconfirm>
				</n-space>
			</n-space>

			<n-alert v-if="!users.length" :title="t('overlays.tts.users.empty')" type="info" />
			<n-grid v-else cols="1 s:1 m:2 l:2" responsive="screen" :x-gap="10" :y-gap="10">
				<n-grid-item v-for="(user, index) of users" :key="index" :span="1">
					<n-card
						class="user-card p-1"
						@click="user.markedForDelete = !user.markedForDelete"
					>
						<n-space align="center" justify="space-between">
							<n-row align-items="center" class="gap-2.5">
								<n-avatar :src="user.avatar" size="large" />
								<n-space vertical size="small" class="info">
									<n-text>{{ user.name }}</n-text>
									<n-text class="description">
										{{ t('overlays.tts.voice') }}: {{ user.voice }} | {{ t('overlays.tts.pitch') }}: {{ user.pitch }} | {{ t('overlays.tts.rate') }}: {{ user.rate }}
									</n-text>
								</n-space>
							</n-row>

							<n-row align-items="center">
								<n-space align="center">
									<IconSpeakerphone
										class="flex cursor-pointer"
										@click.stop="testUserVoice(user)"
									/>
									<n-checkbox
										:checked="user.markedForDelete"
									/>
								</n-space>
							</n-row>
						</n-space>
					</n-card>
				</n-grid-item>
			</n-grid>
		</div>
	</div>
</template>

<style scoped>
.user-card {
	width: 100%;
	border-radius: 10px
}

.user-card :deep(.description) {
	font-size: 12px;
	color: v-bind(descriptionColor)
}

.user-card :deep(.info) {
	gap: 0px !important;
	cursor: pointer
}
</style>
