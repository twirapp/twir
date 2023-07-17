<script setup lang='ts'>
import { IconLogin, IconLogout } from '@tabler/icons-vue';
import type { FunctionalComponent } from 'vue';

const props = withDefaults(defineProps<{
	name: string,
	data: { userName: string, avatar: string } | undefined
	logout: () => Promise<void>
	getLoginLink: () => Promise<{ data: { link: string } }>
	icon: FunctionalComponent<any>
	iconWidth?: number
}>(), {
	iconWidth: 30,
});

async function login() {
	const req = await props.getLoginLink();
	if (!req.data?.link) return;
	window.location.replace(req.data?.link);
}
</script>

<template>
  <tr>
    <td>
      <n-tooltip trigger="hover" placement="left">
        <template #trigger>
          <component :is="props.icon" :width="props.iconWidth" class="icon" />
        </template>
        {{ name }}
      </n-tooltip>
    </td>
    <td>
      <div v-if="data?.userName" class="profile">
        <n-avatar :src="data.avatar" class="avatar" round />
        <n-text>
          {{ data.userName }}
        </n-text>
      </div>
      <n-tag v-else :bordered="false" type="info">
        Not Logged In
      </n-tag>
    </td>
    <td>
      <div class="actions">
        <n-button v-if="data?.userName" strong secondary type="error" @click="logout">
          <IconLogout />
          Logout
        </n-button>
        <n-button v-else trong secondary type="success" @click="login">
          <IconLogin />
          Login
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
