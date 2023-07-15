<script setup lang='ts'>
import { useRolesManager } from '@/api/index.js';

const rolesManager = useRolesManager();
const roles = rolesManager.getAll({});
</script>

<template>
  {{ roles.data }}
</template>
