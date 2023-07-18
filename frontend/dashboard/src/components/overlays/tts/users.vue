<script setup lang='ts'>
import { IconSpeakerphone } from '@tabler/icons-vue';
import {
	NGrid,
	NGridItem,
	NCard,
	NSkeleton,
	NText,
	NAlert,
	NRow,
	NAvatar,
	NSpace,
	NCheckbox,
	NButton,
	NInput,
	useThemeVars,
} from 'naive-ui';
import { computed, ref, UnwrapRef, watch } from 'vue';

import { useTtsOverlayManager, useTwitchGetUsers } from '@/api/index.js';

const themeVars = useThemeVars();
const descriptionColor = computed(() => themeVars.value.textColor3);

const ttsManager = useTtsOverlayManager();
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
		if (!twitchUser) return;

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
</script>

<template>
  <div style="padding: 15px">
    <n-grid v-if="isLoading || !ttsUsersData" :cols="24" :x-gap="10" :y-gap="10">
      <n-grid-item v-for="i in 16" :key="i" :span="12">
        <n-skeleton v-if="!isLoading" size="large" height="60px" :sharp="false" />
      </n-grid-item>
    </n-grid>

    <div v-else>
      <n-space justify="space-between" style="margin-bottom: 10px">
        <n-input v-model:value="testText" placeholder="Text for test user settings" />
        <n-space>
          <n-button
            secondary
            type="info"
            :disabled="!users.length"
            @click="changeMarkedStateForAllUsers(!isSomeUserMarked)"
          >
            {{ isSomeUserMarked ? 'Undo select' : 'Select all' }}
          </n-button>

          <n-button
            secondary
            type="error"
            :disabled="!users.some(u => u.markedForDelete)"
            @click="deleteUsers"
          >
            Delete {{ users.filter(u => u.markedForDelete).length }}
          </n-button>
        </n-space>
      </n-space>

      <n-alert v-if="!users.length" title="No one's created their own customizations yet." type="info" />
      <n-grid v-else cols="1 s:1 m:24 l:24 xl:24" responsive="screen" :x-gap="10" :y-gap="10">
        <n-grid-item v-for="(user, index) of users" :key="index" :span="12">
          <n-card
            class="user-card"
            content-style="padding: 5px"
            @click="user.markedForDelete = !user.markedForDelete"
          >
            <n-space align="center" justify="space-between">
              <n-row align-items="center" style="gap: 10px">
                <n-avatar :src="user.avatar" size="large" />
                <n-space vertical size="small" class="info">
                  <n-text>{{ user.name }}</n-text>
                  <n-text class="description">
                    Voice: {{ user.voice }} | Pitch: {{ user.pitch }} | Rate: {{ user.rate }}
                  </n-text>
                </n-space>
              </n-row>

              <n-row align-items="center">
                <n-space align="center">
                  <IconSpeakerphone style="display: flex" />
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
