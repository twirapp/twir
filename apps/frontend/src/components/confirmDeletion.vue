<script setup lang="ts">
import { ref } from 'vue';

const isVisible = ref(false);

defineEmits<{
  (event: 'onConfirmed'): () => void
}>();
</script>

<template>
  <div class="text-center">
    <v-dialog
      v-model="isVisible"
    >
      <template #activator>
        <slot ></slot>
      </template>

      <v-card>
        <v-card-text>
          Please confirm deletion
        </v-card-text>
        <v-card-actions>
          <v-btn color="primary" block @click="isVisible = false">
            Close Dialog
          </v-btn>
          <v-btn
            color="red" block @click="() => {
              isVisible = false;
              $emit('onConfirmed')
            }"
          >
            Confirm
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>