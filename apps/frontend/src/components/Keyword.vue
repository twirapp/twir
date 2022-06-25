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
        class="bg-red-600 rounded py-5 px-6 mb-4 text-base text-red-700 mb-3"
        role="alert"
      >
        {{ error }}
      </div>
      <div
        class="grid grid-cols-1 gap-1"
      >
        <div>
          <div class="label mb-1">
            <span class="label-text">{{ t('pages.keywords.card.text.title') }}</span>
          </div>
          <Field
            v-model.lazy="keyword.text"
            name="text"
            as="input" 
            type="text"
            :placeholder="t('pages.keywords.card.text.placeholder')"
            :disabled="!keyword.edit"
            class="form-control px-3 py-1.5 text-gray-700 rounded input input-bordered w-full input-sm"
          />
        </div>

        <div class="mt-5">
          <div class="label mb-1">
            <span class="label-text">{{ t('pages.keywords.card.response.title') }}</span>
          </div>
          <Field
            v-model.lazy="keyword.response"
            name="response"
            as="input" 
            type="text"
            :placeholder="t('pages.keywords.card.response.placeholder')"
            :disabled="!keyword.edit"
            class="form-control px-3 py-1.5 text-gray-700 rounded input input-bordered w-full input-sm"
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
            class="form-control px-3 py-1.5 text-gray-700 rounded input input-bordered w-full input-sm"
          />
        </div>
      </div>

      <div class="flex justify-between mt-5">
        <div>
          <button
            v-if="!keyword.edit"
            type="button"
            class="inline-block px-6 py-2.5 bg-purple-600 font-medium text-xs leading-tight uppercase rounded shadow hover:bg-putple-700  focus:outline-none focus:ring-0 transition duration-150 ease-in-out"
            @click="() => {
              keyword.edit = true;
              if (keyword.id) keywordsBeforeEdit?.push(JSON.parse(JSON.stringify(keyword)))
            }"
          >
            {{ t('buttons.edit') }}
          </button>
          <button
            v-else
            class="px-6 py-2.5 inline-block bg-purple-600 text-white font-medium text-xs leading-tight uppercase rounded shadow hover:bg-purple-700  focus:outline-none focus:ring-0  transition duration-150 ease-in-out"
            @click="cancelEdit"
          >
            {{ t('buttons.cancel') }}
          </button>
        </div>
        <div v-if="keyword.edit">
          <button
            v-if="keyword.id"
            type="button"
            class="inline-block px-6 py-2.5 bg-red-600 text-white font-medium text-xs leading-tight uppercase rounded shadow hover:bg-red-700  focus:outline-none focus:ring-0   transition duration-150 ease-in-out"
            @click="deletekeyword"
          >
            {{ t('buttons.delete') }}
          </button>
          <button
            type="submit"
            class="inline-block ml-2 px-6 py-2.5 bg-green-600 text-white font-medium text-xs leading-tight uppercase rounded shadow hover:bg-green-700  focus:outline-none focus:ring-0   transition duration-150 ease-in-out"
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