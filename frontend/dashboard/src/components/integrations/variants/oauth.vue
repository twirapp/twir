<script setup lang="ts">
import { IconLogin, IconLogout } from '@tabler/icons-vue';
import { NButton, NAvatar, useThemeVars } from 'naive-ui';
import { FunctionalComponent } from 'vue';
import { useI18n } from 'vue-i18n';

import { useUserAccessFlagChecker } from '@/api';
import Card from '@/components/card/card.vue';

const themeVars = useThemeVars();

const props = withDefaults(defineProps<{
	title: string,
	description?: string
	isLoading?: boolean
	data: { userName?: string, avatar?: string } | undefined
	logout: () => any
	authLink?: string,
	icon: FunctionalComponent<any>
	iconWidth?: string
	iconColor?: string
}>(), {
	authLink: '',
	description: '',
});

async function login() {
	if (!props.authLink) return;

	window.open(props.authLink, 'Twir connect integration', 'width=800,height=600');
}

const userCanManageIntegrations = useUserAccessFlagChecker('MANAGE_INTEGRATIONS');

const { t } = useI18n();
</script>

<template>
	<card
		:title="title"
		style="height: 100%;"
		:with-stroke="false"
		:icon="icon"
		:icon-width="iconWidth"
		:is-loading="isLoading"
	>
		<template #content>
			<span class="description" v-html="description" />
		</template>

		<template #footer>
			<div class="footer">
				<n-button
					:disabled="!userCanManageIntegrations || !authLink"
					secondary
					size="large"
					:type="data?.userName ? 'error' : 'success'"
					@click="data?.userName ? logout() : login()"
				>
					<div class="button-content">
						<span>
							{{ t(`sharedButtons.${data?.userName ? 'logout' : 'login'}`) }}
						</span>
						<IconLogout v-if="data?.userName" />
						<IconLogin v-else />
					</div>
				</n-button>
				<div v-if="data?.userName" class="profile">
					<n-avatar round :src="data?.avatar" />
					<span>{{ data.userName }}</span>
				</div>
			</div>
		</template>
	</card>
</template>

<style scoped>
.footer {
	display: flex;
	justify-content: space-between;
	align-items: center;
	gap: 5px;
	width: 100%;
	flex-wrap: wrap;
}

.button-content {
	display: flex;
	gap: 4px;
}

.profile {
	display: flex;
	align-items: center;
	padding-top: 4px;
	padding-bottom: 4px;
	padding-left: 18px;
	padding-right: 18px;
	background-color: v-bind('themeVars.buttonColor2');
	border-radius: 4px;
	gap: 8px;
}

.description :deep(a) {
	color: v-bind('themeVars.successColor');
	text-decoration: none;
}
</style>
