<script lang="ts" setup>
export type KeywordType = SetOptional<Omit<Keyword, 'channelId'> & { edit?: boolean }, 'id'>;

import { Keyword } from '@tsuwari/prisma';
import type { SetOptional } from 'type-fest';
import { Ref, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import KeywordComponent from '@/components/Keyword.vue';
import { useUpdatingData } from '@/functions/useUpdatingData';

const { data } = useUpdatingData(`/v1/channels/{dashboardId}/keywords`);
const keywords = ref<Array<KeywordType>>([]);
const keywordsBeforeEdit = ref<Array<KeywordType>>([]);
const { t } = useI18n();

watch(data, (v: any[]) => {
  keywords.value = v;
  keywordsBeforeEdit.value = [];
});

function insert() {
  keywords.value = [
    {
      text: '',
      response: '',
      cooldown: 5,
      enabled: true,
      edit: true,
      isReply: false,
      isRegular: false,
    },
    ...keywords.value,
  ];
}

async function deletekeyword(index: number) {
  keywords.value = keywords.value.filter((_, i) => i !== index);
}

function cancelEdit(keyword: Ref<KeywordType>) {
  const index = keywords.value.indexOf(keyword.value);
  if (keyword.value.id && keywords.value) {
    const editableCommand = keywordsBeforeEdit.value?.find((c) => c.id === keyword.value.id);
    if (editableCommand) {
      keywords.value[index] = {
        ...editableCommand,
        edit: false,
      };

      keywordsBeforeEdit.value = keywordsBeforeEdit.value.filter(
        (v, i) => i !== keywordsBeforeEdit.value.indexOf(editableCommand),
      );
    }
  } else {
    keywords.value = keywords.value.filter((v, i) => i !== index);
  }
}
</script>

<template>
  <div class="m-1.5 md:m-3">
    <div class="flow-root">
      <div class="btn btn-primary btn-sm float-left mb-5 md:w-auto rounded w-full">
        <button
          class="bg-purple-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-purple-700 inline-block leading-tight px-6 py-2.5 rounded shadow text-white text-xs transition uppercase"
          @click="insert">
          {{ t('buttons.addNew') }}
        </button>
      </div>
    </div>

    <masonry-wall :items="keywords" :gap="8">
      <template #default="{ item, index }">
        <div :key="index" class="block card rounded shadow text-white">
          <KeywordComponent
            :keyword="item"
            :keywords="keywords"
            :keywords-before-edit="keywordsBeforeEdit"
            @delete="deletekeyword"
            @cancel-edit="cancelEdit" />
        </div>
      </template>
    </masonry-wall>
  </div>
</template>
