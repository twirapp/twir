<template>
  <router-view v-slot="{ Component }">
    <Suspense>
      <component :is="Component" />
    </Suspense>
  </router-view>
  <ul>
    <li v-for="(item, key) in appMenu" :key="key">
      <router-link :to="item.path">
        {{ AppMenu[key] }}
      </router-link>
    </li>
  </ul>
</template>

<script lang="ts" setup>
import { appMenu, AppMenu } from './router.js';
</script>
