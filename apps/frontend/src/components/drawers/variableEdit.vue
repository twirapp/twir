<script lang="ts" setup>
import Editor, { useMonaco } from '@guolao/vue-monaco-editor';
import { mdiPlus, mdiClose } from '@mdi/js';
import { type ChannelCustomvar } from '@tsuwari/typeorm/entities/ChannelCustomvar';
import { onUnmounted, ref } from 'vue';

import confirmDeletion from '@/components/confirmDeletion.vue';

const props = defineProps<{
  variable: ChannelCustomvar
}>();
const emits = defineEmits<{
  (event: 'cancel'): () => void
}>();
const variable = ref(props.variable);
const { monacoRef, unload } = useMonaco();

onUnmounted(() => !monacoRef.value && unload());

function onDelete() {
  console.log(variable);
  emits('cancel');
}
</script>

<template>
  <div>
    <v-list-item>
      <div class="d-flex justify-space-between">
        <h1>Edit variable</h1>
        <div class="d-flex d-inline">
          <v-btn size="small" class="mt-1 mr-2" @click="$emit('cancel')">
            Cancel
          </v-btn>
          <confirmDeletion :cb="() => onDelete()">
            <v-btn color="red" size="small" class="mt-1 mr-2">
              Delete
            </v-btn>
          </confirmDeletion>
          <v-btn color="green" size="small" class="mt-1">
            Save
          </v-btn>
        </div>
      </div>
    </v-list-item>

    <v-divider></v-divider>

    <v-list-item>
      <v-form class="mt-2">
        <div
          class="d-flex flex-column"
        >
          <v-text-field 
            v-model="variable.name" 
            label="Name" 
            :rules="[
              v => !!v || 'Field is required'
            ]"
          />

          <v-select
            v-model="variable.type"
            label="Variable type"
            :items="['SCRIPT', 'TEXT']"
          ></v-select>

          <v-textarea
            v-if="variable.type === 'TEXT'"
            v-model="variable.response"
            auto-grow
            label="Variable response"
            rows="1"
            row-height="5"
            class="mt-2"
          />
        </div>

        <v-sheet
          v-if="variable.type === 'SCRIPT'"
          rounded
          class="mt-2 pa-4"
          color="#484749"
        >
          Do not forget about semicolons when writing scripts. It's important.
        </v-sheet>
        <Editor 
          v-if="variable.type === 'SCRIPT'"
          v-model:value="variable.evalValue"
          class="mt-4"
          height="60vh"
          theme="vs-dark"
          defaultLanguage="javascript"
          :options="{}"
        />
      </v-form>
    </v-list-item>
  </div>
</template>