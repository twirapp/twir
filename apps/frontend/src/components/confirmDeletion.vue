<script setup lang="ts">
import { ref } from 'vue';

const modalState = ref(false);

defineProps<{
  cb: () => any | Promise<any>;
}>();
</script>

<template>
  <div class="text-center">
    <span @click="modalState = true"><slot></slot></span>

    <v-dialog
      v-model="modalState"
      max-width="200"
      max-heigh="200"
    >
      <v-card>
        <v-card-text>
          Are you sure?
        </v-card-text>
        <v-card-actions class="d-flex justify-end">
          <v-btn color="green" variant="outlined" @click="modalState = false">
            No
          </v-btn>
          <v-btn
            class="ml-1" color="red" variant="outlined" @click="() => {
              modalState = false;
              cb()
            }"
          >
            Yes
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>