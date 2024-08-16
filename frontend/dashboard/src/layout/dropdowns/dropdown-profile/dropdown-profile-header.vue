<script setup lang="ts">
import { NAvatar, NText } from 'naive-ui'

import { useProfile } from '@/api/auth.js'
import { resolveUserName } from '@/helpers/resolveUserName.js'

const { data: profile } = useProfile()
</script>

<template>
	<div class="flex items-center px-3 py-2">
		<div v-if="profile" class="flex gap-3 items-center">
			<NAvatar :src="profile?.avatar" round />

			<div class="flex flex-col">
				<NText :depth="3" class="text-xs">
					Logged as
				</NText>

				<NText>
					{{ resolveUserName(profile.login, profile.displayName) }}
				</NText>
			</div>
		</div>
	</div>
</template>
