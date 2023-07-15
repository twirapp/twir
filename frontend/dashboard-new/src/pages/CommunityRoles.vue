<script setup lang='ts'>
import {
	NCard,
	NSpace,

} from 'naive-ui';

import { useRolesManager } from '@/api/index.js';

const rolesManager = useRolesManager();
const roles = rolesManager.getAll({});
</script>

<template>
  <n-space>
    <n-card v-for="role in roles.data.value" :key="role.id" title="Small Card" size="small">
      {{ role.name }}
    </n-card>
  </n-space>
</template>
