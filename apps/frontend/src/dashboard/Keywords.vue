<script lang="ts" setup>
export type KeywordType = SetOptional<Omit<Keyword, 'channelId'> & { edit?: boolean }, 'id'>

import { useStore } from '@nanostores/vue';
import { Keyword } from '@tsuwari/prisma';
import { useAxios } from '@vueuse/integrations/useAxios';
import type { SetOptional } from 'type-fest';
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import KeywordComponent from '@/components/Keyword.vue';
import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

const selectedDashboard = useStore(selectedDashboardStore);
const { execute, data: axiosData } = useAxios(`/v1/channels/${selectedDashboard.value.channelId}/keywords`, api, { immediate: false });
const keywords = ref<Array<KeywordType>>([]);
const keywordsBeforeEdit = ref<Array<KeywordType>>([]);
const { t } = useI18n();

selectedDashboardStore.subscribe((v) => {
  execute(`/v1/channels/${v.channelId}/keywords`);
});

watch(axiosData, (v: any[]) => {
  keywords.value = v;
  keywordsBeforeEdit.value = [];
});

function insert() {
  keywords.value.unshift({
    text: '',
    response: '',
    cooldown: 5,
    enabled: true,
    edit: true,
  });
}

async function deletekeyword(index: number) {
  keywords.value = keywords.value.filter((_, i) => i !== index);
}
</script>

<template>
  <div class="m-3">
    <div class="p-1">
      <div class="flow-root">
        <div class="float-left rounded btn btn-primary btn-sm w-full mb-1 md:w-auto">
          <button
            class="px-6 py-2.5 inline-block bg-purple-600 text-white font-medium text-xs leading-tight uppercase rounded shadow hover:bg-purple-700    focus:outline-none focus:ring-0  transition duration-150 ease-in-out"
            @click="insert"
          >
            {{ t('buttons.addNew') }}
          </button>
        </div>
      </div>
    </div>

    <div class="grid lg:grid-cols-3 md:grid-cols-2 grid-cols-1 gap-2">
      <div
        v-for="keyword, index of keywords"
        :key="index"
        class="block rounded card text-white shadow"
      >
        <KeywordComponent 
          :keyword="keyword"
          :keywords="keywords"
          :keywords-before-edit="keywordsBeforeEdit"
          @delete="deletekeyword"
        />
      </div>
    </div>
  </div>
</template>
