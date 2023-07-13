<script setup lang='ts'>
import { IconSettings } from '@tabler/icons-vue';
import { NTooltip, NButton, NModal } from 'naive-ui';
import { ref } from 'vue';
import { defineSlots, FunctionalComponent } from 'vue/dist/vue.js';

const props = withDefaults(defineProps<{
	name: string,
	hasSaveButton: boolean
}>(), {
	hasSaveButton: false,
});

defineSlots<{
	icon: FunctionalComponent<any>
	settings: FunctionalComponent<any>
}>();

const showSettings = ref(false);
const modalWidth = '600px';
</script>

<template>
  <tr>
    <td>
      <n-tooltip trigger="hover" placement="left">
        <template #trigger>
          <slot name="icon" />
        </template>
        {{ name }}
      </n-tooltip>
    </td>
    <td></td>
    <td>
      <n-button strong secondary type="info" @click="showSettings = true">
        <IconSettings />
        Settings
      </n-button>
    </td>
  </tr>

  <n-modal
    v-model:show="showSettings"
    :mask-closable="false"
    :segmented="true"
    preset="card"
    :title="name"
    class="modal"
    :style="{
      width: modalWidth,
      position: 'fixed',
      left: `calc(50% - ${modalWidth}/2)`,
      top: '50px',
    }"
  >
    <template #header>
      {{ name }}
    </template>
    <slot name="settings" />

    <template #action>
      <n-button @click="showSettings = false">
        Close
      </n-button>
    </template>
  </n-modal>
</template>
