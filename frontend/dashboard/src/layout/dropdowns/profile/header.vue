<script setup lang="ts">
import { NAvatar, NText } from 'naive-ui';

import { useProfile } from '@/api';
import { createUserName } from '@/helpers';

const { data: profileData } = useProfile();
</script>

<template>
	<div class="header">
		<div v-if="profileData" style="display: flex; gap: 12px; align-items: center;">
			<n-avatar :src="profileData?.avatar" round />
			<div style="display: flex; flex-direction: column;">
				<n-text :depth="3" style="font-size: 11px">
					Logged as
				</n-text>
				<n-text>
					{{ createUserName(profileData.login, profileData.displayName) }}
				</n-text>
			</div>
		</div>
	</div>
</template>

<style scoped>
.header {
	display: flex;
	align-items: center;
	padding: 8px 12px;
}
</style>
