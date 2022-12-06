<script lang="ts" setup>
import { useRouter } from 'vue-router';

import { publicRoutes } from './menu';

const router = useRouter();
</script>

<template>
  <v-list lines="one">
    <v-list-item
      v-for="item in publicRoutes"
      :key="item.name"
      :title="item.name.charAt(0).toUpperCase() + item.name.substring(1)"
      :prepend-icon="item.icon"
      @click="router.push(item.path)"
    >
    </v-list-item>
  </v-list>
</template>