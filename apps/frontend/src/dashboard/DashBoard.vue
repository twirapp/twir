<!-- eslint-disable @typescript-eslint/ban-ts-comment -->
<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { get, useTitle, useIntervalFn } from '@vueuse/core';
// @ts-ignore
import FilePondPluginFileValidateType from 'filepond-plugin-file-validate-type/dist/filepond-plugin-file-validate-type.esm.js';
// @ts-ignore
import FilePondPluginImagePreview from 'filepond-plugin-image-preview/dist/filepond-plugin-image-preview.esm.js';
import { Form, Field } from 'vee-validate';
import { ref } from 'vue';
import vueFilePond from 'vue-filepond';
import { useI18n } from 'vue-i18n';

import 'filepond/dist/filepond.min.css';
import 'filepond-plugin-image-preview/dist/filepond-plugin-image-preview.min.css';

import MyBtn from '@/components/elements/MyBtn.vue';
import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

const title = useTitle();
title.value = 'Dashboard';

const { t } = useI18n({
  useScope: 'global',
  inheritLocale: true,
});

const isBotMod = ref(false);

const dashboard = useStore(selectedDashboardStore);

useIntervalFn(
  async () => {
    const dash = get(dashboard);
    if (!dash) return;

    const { data } = await api(`v1/channels/${dashboard.value.channelId}/bot/checkmod`);

    isBotMod.value = data;
  },
  1000,
  { immediate: true },
);

selectedDashboardStore.subscribe(() => {
  isBotMod.value = false;
});

async function patchBotConnection(action: 'join' | 'leave') {
  await api.patch(`v1/channels/${dashboard.value.channelId}/bot/connection`, {
    action,
  });
}

const feedBack = ref('');
const myFiles = ref([]);
const FilePond = vueFilePond(FilePondPluginFileValidateType, FilePondPluginImagePreview);
function updateFiles(files: any) {
  myFiles.value = files.map((f: any) => f.file);
}

async function sendForm() {
  const uploadedFiles = [];

  if (!feedBack.value) return;
  if (myFiles.value.length) {
    const formData = new FormData();

    for (const file of myFiles.value) {
      formData.append('files', file);
    }

    const filesRequest = await api.post('/v1/files', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    uploadedFiles.push(...filesRequest.data);
  }

  await api.post('/v1/feedback', {
    text: feedBack.value,
    files: uploadedFiles,
  });

  myFiles.value = [];
  feedBack.value = '';
}
</script>

<template>
  <div class="m-1.5 md:m-3">
    <div class="masonry md:masonry-md sm:masonry-sm space-y-2">
      <div class="block break-inside card rounded shadow text-white">
        <h2
          class="border-b border-gray-700 card-title flex font-bold justify-center outline-none p-2">
          <p>{{ t('pages.dashboard.widgets.status.title') }}</p>
        </h2>
        <div class="p-4 w-full">
          <div
            class="mb-4 px-6 py-5 rounded text-base"
            :class="{ 'bg-[#ED4245]': !isBotMod, 'bg-green-600': isBotMod }">
            <div v-if="!isBotMod">
              <div v-html="t('pages.dashboard.widgets.status.notMod')" />
            </div>
            <div v-else>
              {{ t('pages.dashboard.widgets.status.mod') }}
            </div>
          </div>

          <div
            class="flex flex-col md:flex-row md:justify-end md:space-x-1 md:space-y-0 md:text-right space-y-1">
            <button
              type="button"
              class="bg-red-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-red-700 inline-block leading-tight px-6 py-2.5 rounded shadow text-white text-xs transition uppercase"
              @click="() => patchBotConnection('leave')">
              {{ t('pages.dashboard.widgets.status.buttons.leave') }}
            </button>
            <button
              type="button"
              class="bg-green-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-green-700 inline-block leading-tight px-6 py-2.5 rounded shadow text-white text-xs transition uppercase"
              @click="() => patchBotConnection('join')">
              {{ t('pages.dashboard.widgets.status.buttons.join') }}
            </button>
          </div>
        </div>
      </div>

      <div class="block break-inside card rounded shadow text-white">
        <h2
          class="border-b border-gray-700 card-title flex font-bold justify-center outline-none p-2">
          <p>{{ t('pages.dashboard.widgets.feedback.title') }}</p>
        </h2>
        <Form @submit="sendForm">
          <div class="inline-block p-4 w-full">
            <div class="flex justify-center">
              <div class="mb-3 w-full">
                <Field
                  v-model.trim="feedBack"
                  name="feedback"
                  as="textarea"
                  :placeholder="t('pages.dashboard.widgets.feedback.placeholder')"
                  class="bg-clip-padding bg-white block border border-gray-300 border-solid ease-in-out focus:bg-white focus:border-blue-600 focus:outline-none focus:text-gray-700 font-normal form-control m-0 px-3 py-1.5 rounded text-base text-gray-700 transition w-full"
                  rows="3" />
                <file-pond
                  ref="pond"
                  name="test"
                  class-name="mt-3"
                  label-idle="Drop files here..."
                  :allow-multiple="true"
                  accepted-file-types="image/jpeg, image/png"
                  :files="myFiles"
                  :credits="[]"
                  :max-files="5"
                  @updatefiles="updateFiles" />
              </div>
            </div>

            <div class="text-right">
              <MyBtn color="green" type="submit">
                {{ t('pages.dashboard.widgets.feedback.buttons.send') }}
              </MyBtn>
            </div>
          </div>
        </Form>
      </div>
    </div>
  </div>
</template>

<style>
.filepond--root {
  max-height: 350px;
}

@media (min-width: 30em) {
  .filepond--item {
    width: calc(50% - 0.5em);
  }
}

@media (min-width: 50em) {
  .filepond--item {
    width: calc(33.33% - 0.5em);
  }
}
</style>
