<script setup lang="ts">
import { NAvatar, NText } from 'naive-ui';

import { useProfile } from '@/api';
import { createUserName } from '@/helpers';

const { data: profileData } = useProfile();
</script>

<template>
	<div class="flex items-center px-3 py-2">
		<div v-if="profileData" class="flex gap-3 items-center">
			<n-avatar :src="profileData?.avatar" round />
			<div class="flex flex-col">
				<n-text :depth="3" class="text-xs">
					Logged as
				</n-text>
				<n-text>
					{{ createUserName(profileData.login, profileData.displayName) }}
				</n-text>
			</div>
		</div>
	</div>
</template>
