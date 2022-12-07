<script lang="ts" setup>
import { computed, onMounted, ref, watch } from 'vue';
import { RouterView } from 'vue-router';
import { useDisplay } from 'vuetify';

import Sidebar from './components/sidebar/menu.vue';
import Profile from './components/sidebar/profile.vue';

const drawer = ref(false);

const { mobile } = useDisplay();

onMounted(() => {
  if (!mobile.value) drawer.value = true;
});

watch(mobile, (_, v) => {
  drawer.value = !v;
});
</script>

<template>
  <v-layout>
    <v-navigation-drawer v-model="drawer" color="grey-darken-5" :expand-on-hover="!mobile" :rail="!mobile"> 
      <Profile />
      <v-divider></v-divider>
      <Sidebar />
    </v-navigation-drawer>
    <v-app-bar color="grey-darken-5">
      <v-app-bar-nav-icon v-if="mobile" variant="text" @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
    </v-app-bar>
    <v-main>
      <div style="padding: 8px;">
        <RouterView />
      </div>
    </v-main>
  </v-layout>
</template>