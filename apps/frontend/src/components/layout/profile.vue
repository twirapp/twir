<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { ref } from 'vue';

import { userStore } from '@/stores/user';

const menu = ref(false);
const user = useStore(userStore);
</script>

<template>
  <v-menu
    v-model="menu"
    :close-on-content-click="false"
    location="end"
  >
    <template #activator="{ props }">
      <!-- <v-btn -->
      <!-- color="indigo" -->
      <!-- v-bind="props" -->
      <!-- > -->
      <img
        v-bind="props" 
        class="rounded-circle" 
        style="height:50px;cursor: pointer;" 
        src="https://randomuser.me/api/portraits/women/85.jpg" 
      />
      <!-- </v-btn> -->
    </template>

    <v-card>
      <v-list>
        <v-list-item
          prepend-avatar="https://cdn.vuetifyjs.com/images/john.jpg"
          title="John Leider"
          subtitle="Founder of Vuetify"
        >
        </v-list-item>
      </v-list>

      <v-divider></v-divider>

      <v-list>
        <v-list-item>
          <v-switch
            color="purple"
            label="Enable hints"
            hide-details
          ></v-switch>
        </v-list-item>
      </v-list>

      <v-card-actions>
        <v-btn
          color="red"
          variant="outlined"
          block
          @click="menu = false"
        >
          Logout
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-menu>
</template>