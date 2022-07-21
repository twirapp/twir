<script lang="ts" setup>
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { api } from '@/plugins/api';

const { availableLocales } = useI18n();

const locales = (availableLocales as string[]).map((l) => l.toUpperCase());
const form = ref({
  imageSrc: null,
  userName: null,
  messages: Object.fromEntries(locales.map((l) => [l, {}])) as Record<string, any>,
});

async function sendForm() {
  await api.post('/admin/notifications', {
    ...form.value,
    messages: Object.entries(form.value.messages).map((m) => ({ langCode: m[0], ...m[1] })),
  });
}
</script>

<template>
  <div class="m-1.5 md:m-3">
    <div class="p-1">
      <div class="gap-2 grid grid-cols-1 lg:grid-cols-2">
        <div class="block card rounded shadow text-white">
          <h2 class="border-b border-gray-700 card-title flex font-bold justify-center outline-none p-2">
            <p>Notifications</p>
          </h2>
          <div class="p-4 space-y-5 w-full">
            <div class="flex flex-col mb-4 rounded space-y-3">
              <input
                v-model="form.imageSrc"
                class="border border-grey-light flex-1 form-control h-10 input leading-normal px-3 relative rounded text-gray-700 w-full"
                placeholder="Image src"
              >
              <input
                v-model="form.userName"
                class="border border-grey-light flex-1 form-control h-10 input leading-normal px-3 relative rounded text-gray-700 w-full"
                placeholder="Username"
              >
            </div>

            <div
              v-for="lang in locales"
              :key="lang"
            >
              {{ lang }}
              <input
                v-model="form.messages[lang].title"
                class="border border-grey-light flex-1 form-control h-5 input leading-normal px-3 relative rounded text-gray-700 w-full"
                placeholder="Alert title"
              >
              <input
                v-model="form.messages[lang].text"
                class="border border-grey-light flex-1 form-control h-5 input leading-normal px-3 relative rounded text-gray-700 w-full"
                placeholder="Alert text"
              >
            </div>

            <div class="flex flex-col md:flex-row md:justify-end md:space-x-1 md:space-y-0 md:text-right mt-5 space-y-1">
              <button
                type="button"
                class="bg-green-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-green-700 inline-block leading-tight px-6 py-2.5 rounded shadow text-white text-xs transition uppercase"
                @click="sendForm"
              >
                Send
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
