<script lang="ts" setup>
import { computed, onMounted, ref, watch } from 'vue';
import { RouterView } from 'vue-router';
import { useDisplay } from 'vuetify';

import Sidebar from './components/sidebar/menu.vue';
import Profile from './components/sidebar/profile.vue';

const drawer = ref(false);
const display = useDisplay();

const isMobile = computed(() => {
  return display.xs.value || display.sm.value;
});

onMounted(() => {
  if (!isMobile.value) drawer.value = true;
});

watch(isMobile, (_, v) => {
  console.log(isMobile);
  drawer.value = !v;
});
</script>

<template>
  <v-layout>
    <v-navigation-drawer v-model="drawer" color="grey-darken-5" :expand-on-hover="!isMobile" :rail="!isMobile"> 
      <Profile />
      <v-divider></v-divider>
      <Sidebar />
    </v-navigation-drawer>
    <v-app-bar color="grey-darken-5">
      <v-app-bar-nav-icon v-if="isMobile" variant="text" @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
    </v-app-bar>
    <v-main>
      <div style="padding: 8px;">
        <RouterView />
      </div>
    </v-main>
  </v-layout>
</template>