<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { Form, Field } from 'vee-validate';
import { toRef } from 'vue';
import { useI18n } from 'vue-i18n';

import { KeywordType } from '@/dashboard/Keywords.vue';
import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

const props = defineProps<{
  keyword: KeywordType,
  keywords: KeywordType[],
  keywordsBeforeEdit: KeywordType[]
}>();

const keyword = toRef(props, 'keyword');
const keywords = toRef(props, 'keywords');
const keywordsBeforeEdit = toRef(props, 'keywordsBeforeEdit');
const selectedDashboard = useStore(selectedDashboardStore);
const { t } = useI18n();
const emit = defineEmits<{
  (e: 'delete', index: number): void
}>();

async function saveKeyword() {
  const index = keywords.value.indexOf(keyword.value);

  let data;

  if (keyword.value.id) {
    const request = await api.patch(`/v1/channels/${selectedDashboard.value.channelId}/keywords/${keyword.value.id}`, keyword.value);  
    data = request.data;
  } else {
    const request = await api.post(`/v1/channels/${selectedDashboard.value.channelId}/keywords`, keyword.value);
    data = request.data;
  }

  keywords.value[index] = data;
}

async function deletekeyword() {
  const index = keywords.value.indexOf(keyword.value);
  if (keyword.value.id) {
    await api.delete(`/v1/channels/${selectedDashboard.value.channelId}/keywords/${keyword.value.id}`);
  }

  emit('delete', index);
}

function cancelEdit() {
  const index = keywords.value.indexOf(keyword.value);
  if (keyword.value.id && keywords.value) {
    const editableCommand = keywordsBeforeEdit.value?.find(c => c.id === keyword.value.id);
    if (editableCommand) {
      keywords.value[index] = {
        ...editableCommand,
        edit: false,
      };
      keywordsBeforeEdit.value?.splice(keywordsBeforeEdit.value.indexOf(editableCommand), 1);
    }
  } else {
    keywords.value?.splice(index, 1);
  }
}
</script>

<template>
  <div class="p-4">
    <Form
      v-slot="{ errors }"
      @submit="saveKeyword"
    >
      <div
        v-for="error of errors"
        :key="error"
        class="bg-red-600 mb-4 px-6 py-5 rounded text-base text-red-700"
        role="alert"
      >
        {{ error }}
      </div>
      <div
        class="gap-1 grid grid-cols-1"
      >
        <div>
          <div class="label mb-1">
            <span class="label-text">{{ t('pages.keywords.card.text.title') }}</span>
          </div>
          <Field
            v-model.lazy.trim="keyword.text"
            name="text"
            as="input" 
            type="text"
            :placeholder="t('pages.keywords.card.text.placeholder')"
            :disabled="!keyword.edit"
            class="form-control input input-bordered input-sm px-3 py-1.5 rounded text-gray-700 w-full"
          />
        </div>

        <div class="mt-5">
          <div class="label mb-1">
            <span class="label-text">{{ t('pages.keywords.card.response.title') }}</span>
          </div>
          <Field
            v-model.lazy.trim="keyword.response"
            name="response"
            as="input" 
            type="text"
            :placeholder="t('pages.keywords.card.response.placeholder')"
            :disabled="!keyword.edit"
            class="form-control input input-bordered input-sm px-3 py-1.5 rounded text-gray-700 w-full"
          />
        </div>

        <div class="mt-5">
          <div class="label mb-1">
            <span class="label-text">{{ t('pages.keywords.card.cooldown.title') }}</span>
          </div>
          <Field
            v-model.number="keyword.cooldown"
            name="cooldown"
            as="input" 
            type="number"
            :placeholder="t('pages.keywords.card.cooldown.placeholder')"
            :disabled="!keyword.edit"
            class="form-control input input-bordered input-sm px-3 py-1.5 rounded text-gray-700 w-full"
          />
        </div>
      </div>

      <div class="flex justify-between mt-5">
        <div>
          <button
            v-if="!keyword.edit"
            type="button"
            class="bg-purple-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-putple-700 inline-block leading-tight px-6 py-2.5 rounded shadow text-xs transition uppercase"
            @click="() => {
              keyword.edit = true;
              if (keyword.id) keywordsBeforeEdit?.push(JSON.parse(JSON.stringify(keyword)))
            }"
          >
            {{ t('buttons.edit') }}
          </button>
          <button
            v-else
            class="bg-purple-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-purple-700 inline-block leading-tight px-6 py-2.5 rounded shadow text-white text-xs transition uppercase"
            @click="cancelEdit"
          >
            {{ t('buttons.cancel') }}
          </button>
        </div>
        <div v-if="keyword.edit">
          <button
            v-if="keyword.id"
            type="button"
            class="bg-red-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-red-700 inline-block leading-tight px-6 py-2.5 rounded shadow text-white text-xs transition uppercase"
            @click="deletekeyword"
          >
            {{ t('buttons.delete') }}
          </button>
          <button
            type="submit"
            class="bg-green-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-green-700 inline-block leading-tight ml-2 px-6 py-2.5 rounded shadow text-white text-xs transition uppercase"
          >
            {{ t('buttons.save') }}
          </button>
        </div>
      </div>
    </Form>
  </div>
</template>

<style scoped>
input, select {
  @apply border-inherit
}
input:disabled, select:disabled {
  @apply bg-zinc-400 opacity-100 border-transparent
}
</style>