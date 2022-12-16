<script lang="ts" setup>
import { computed, onMounted, ref, watch } from 'vue';
import { RouterView } from 'vue-router';
import { useDisplay } from 'vuetify';

import Profile from './components/layout/profile.vue';
import Channel from './components/sidebar/channel.vue';
import Sidebar from './components/sidebar/menu.vue';

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
    <v-navigation-drawer v-model="drawer" color="#202020" :expand-on-hover="!mobile" :rail="!mobile"> 
      <Channel />
      <v-divider></v-divider>
      <Sidebar />
    </v-navigation-drawer>
    <v-app-bar color="#202020">
      <template #prepend>
        <v-app-bar-nav-icon v-if="mobile" variant="text" @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      </template>

      <template #append>
        <Profile />
      </template>
    </v-app-bar>
    <v-main>
      <div class="main">
        <RouterView />
      </div>
    </v-main>
  </v-layout>
</template>

<style>
.main {
  width: 80%;
  margin: 0 auto;
  padding:8px;
}
</style>