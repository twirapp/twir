<script setup lang='ts'>
import { IconLogin, IconLogout } from '@tabler/icons-vue';
import { NButton, NTooltip } from 'naive-ui';
import type { FunctionalComponent } from 'vue';
import { defineSlots } from 'vue';

const props = defineProps<{
	name: string,
	data: { userName: string, avatar: string } | undefined
	logout: () => Promise<void>
	getLoginLink: () => Promise<{ link: string }>
}>();

defineSlots<{
	icon: FunctionalComponent<any>
}>();

async function login() {
	const { link } = await props.getLoginLink();
	if (link) window.location.replace(link);
}
</script>

<template>
  <tr>
    <td>
      <n-tooltip trigger="hover" placement="left">
        <template #trigger>
          <slot name="icon" />
        </template>
        {{ name }}
      </n-tooltip>
    </td>
    <td>
      <n-text v-if="data">
        {{ data.userName }}
      </n-text>
      <span v-else class="badge">Not Logged In</span>
    </td>
    <td>
      <div class="actions">
        <n-button strong secondary type="error" @click="logout">
          <IconLogout />
          Logout
        </n-button>
        <n-button trong secondary type="success" @click="login">
          <IconLogin />
          Login
        </n-button>
      </div>
    </td>
  </tr>
</template>

<style scoped>
.badge {
	background-color: #59c4c3;
	color: #0a0a0a;
	padding: 5px 10px;
	border-radius: 5px;
	font-size: 12px;
	font-weight: 500;
}

.actions {
	display: flex;
	justify-content: flex-end;
	align-items: center;
	gap: 5px;
	width: auto;
}
</style>
