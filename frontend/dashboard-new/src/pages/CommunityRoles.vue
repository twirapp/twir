<script setup lang='ts'>
import { IconPlus } from '@tabler/icons-vue';
import {
	NCard,
	NSpace,
	NText,
} from 'naive-ui';
import { ref } from 'vue';

import { useRolesManager } from '@/api/index.js';

const rolesManager = useRolesManager();
const { data: roles } = rolesManager.getAll({});

const isModalShowed = ref(false);
</script>

<template>
  <n-space align="center" justify="center" vertical>
    <n-card class="card" size="small" bordered hoverable>
      <n-space align="center" justify="center" vertical>
        <n-text class="text">
          <IconPlus />
        </n-text>
      </n-space>
    </n-card>
    <n-card
      v-for="role in roles?.roles"
      :key="role.id"
      size="small"
      class="card"
      hoverable
      @click="() => console.log('clicked')"
    >
      <n-space align="center" justify="center" vertical>
        <n-text class="text">
          {{ role.name }}
        </n-text>
      </n-space>
    </n-card>
  </n-space>
</template>

<style scoped>
.card {
	min-width: 400px;
	cursor: pointer
}

.card .text {
	font-size: 30px;
}
</style>
