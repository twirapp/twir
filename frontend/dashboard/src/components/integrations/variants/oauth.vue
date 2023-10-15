<script setup lang='ts'>
import { IconLogin, IconLogout } from '@tabler/icons-vue';
import { NButton, NTooltip, NAvatar, NText, NTag } from 'naive-ui';
import type { FunctionalComponent } from 'vue';
import { useI18n } from 'vue-i18n';

import { useUserAccessFlagChecker } from '@/api/index.js';

const props = withDefaults(defineProps<{
	name: string,
	data: { userName?: string, avatar?: string } | undefined
	logout: () => any
	authLink?: string,
	icon: FunctionalComponent<any>
	iconWidth?: number
	iconColor?: string
	description?: string
}>(), {
	iconWidth: 30,
	authLink: '',
});

async function login() {
	if (!props.authLink) return;
	window.location.replace(props.authLink);
}

const userCanManageIntegrations = useUserAccessFlagChecker('MANAGE_INTEGRATIONS');

const { t } = useI18n();
</script>

<template>
	<tr>
		<td>
			<n-tooltip trigger="hover" placement="left">
				<template #trigger>
					<component
						:is="props.icon" :width="props.iconWidth" :style="{ fill: props.iconColor }"
						class="icon"
					/>
				</template>
				{{ name }}
			</n-tooltip>
		</td>
		<td>
			<div style="display: flex; flex-direction: column">
				<div>
					<div v-if="data?.userName" class="profile">
						<n-avatar :src="data.avatar" class="avatar" round />
						<n-text>
							{{ data.userName }}
						</n-text>
					</div>
					<n-tag v-else :bordered="false" type="info">
						{{ t('integrations.notLoggedIn') }}
					</n-tag>
				</div>
				<span v-if="description" style="font-size: 11px">{{ description }}</span>
			</div>
		</td>
		<td>
			<div class="actions">
				<n-button
					v-if="data?.userName" :disabled="!userCanManageIntegrations" strong secondary
					type="error" @click="logout"
				>
					<IconLogout />
					{{ t('sharedButtons.logout') }}
				</n-button>
				<n-button
					v-else :disabled="!userCanManageIntegrations || !authLink" trong secondary
					type="success" @click="login"
				>
					<IconLogin />
					{{ t('sharedButtons.login') }}
				</n-button>
			</div>
		</td>
	</tr>
</template>

<style scoped>
.icon {
	display: flex;
}

.actions {
	display: flex;
	justify-content: flex-end;
	align-items: center;
	gap: 5px;
	width: auto;
}

.profile {
	display: flex;
	align-items: center;
	gap: 5px;
}
</style>
