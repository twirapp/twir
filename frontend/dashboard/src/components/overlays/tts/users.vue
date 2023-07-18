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
	useThemeVars,
} from 'naive-ui';
import { computed } from 'vue';

import { useTtsOverlayManager, useTwitchGetUsers } from '@/api/index.js';

const themeVars = useThemeVars();
const descriptionColor = computed(() => themeVars.value.textColor3);

const ttsManager = useTtsOverlayManager();
const { data, isLoading } = ttsManager.getUsersSettings();

const usersIds = computed(() => {
	if (!data?.value?.data) return [];
	return data.value.data.map((user) => user.userId);
});
const twitchUsers = useTwitchGetUsers({
	ids: usersIds,
});

const users = computed(() => {
	return data?.value?.data?.map((user) => {
		const twitchUser = twitchUsers?.data?.value?.users.find((twitchUser) => twitchUser.id === user.userId);
		return {
			...user,
			avatar: twitchUser?.profileImageUrl ?? '',
			name: twitchUser?.displayName,
		};
	}) ?? [];
});
</script>

<template>
  <div style="padding: 15px">
    <n-grid v-if="isLoading || !data" :cols="24" :x-gap="10" :y-gap="10">
      <n-grid-item v-for="i in 16" :key="i" :span="12">
        <n-skeleton v-if="!isLoading" size="large" height="60px" :sharp="false" />
      </n-grid-item>
    </n-grid>
    <n-alert v-if="isLoading || !data?.data || !data.data.length" title="It's too quiet in here..." />
    <n-grid v-else :cols="24" :x-gap="10" :y-gap="10">
      <n-grid-item v-for="(user, index) of users" :key="index" :span="12">
        <n-card class="user-card" content-style="padding: 5px">
          <n-space align="center" justify="space-between">
            <n-row align-items="center" style="gap: 10px">
              <n-avatar :src="user.avatar" />
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
                <n-checkbox />
              </n-space>
            </n-row>
          </n-space>
        </n-card>
      </n-grid-item>
    </n-grid>
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
}
</style>
