<script setup lang="ts">
import { IconSettings, IconUser } from '@tabler/icons-vue'
import { NButton } from 'naive-ui'
import { computed } from 'vue'
import { RouterLink } from 'vue-router'

import { useProfile } from '@/api'
import ButtonLogout from '@/layout/buttons/button-logout.vue'

const profile = useProfile()
const isAdmin = computed(() => profile.data.value?.isBotAdmin)
</script>

<template>
	<div class="flex flex-col gap-1 p-2">
		<RouterLink :to="{ name: 'Settings' }" #="{ navigate, href }" custom>
			<NButton
				:href="href"
				secondary
				block
				type="info"
				tag="a"
				@click="navigate"
			>
				<template #icon>
					<IconSettings />
				</template>

				Settings
			</NButton>
		</RouterLink>

		<RouterLink v-if="isAdmin" :to="{ name: 'AdminPanel' }" #="{ navigate, href }" custom>
			<NButton
				:href="href"
				secondary
				block
				type="info"
				tag="a"
				@click="navigate"
			>
				<template #icon>
					<IconUser />
				</template>

				Admin Panel
			</NButton>
		</RouterLink>

		<ButtonLogout />
	</div>
</template>
