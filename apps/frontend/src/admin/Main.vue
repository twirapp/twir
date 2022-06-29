<script lang="ts" setup>
import { reactive, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { api } from '@/plugins/api';

const { t } = useI18n({
  useScope: 'global',
  inheritLocale: true,
});

const form = ref({
  title: null,
  text: '',
  imageSrc: null,
  userName: null,
});

async function sendForm() {
  await api.post('/admin/notifications', form.value);
}

</script>

<template>
  <div class="m-1.5 md:m-3">
    <div class="p-1">
      <div class="gap-2 grid grid-cols-1 lg:grid-cols-3">
        <div
          class="block card rounded shadow text-white"
        >
          <h2 class="border-b border-gray-700 card-title flex font-bold justify-center outline-none p-2">
            <p>Notifications</p>

          <!-- <a
              href="/"
              class="btn btn-outline btn-error btn-sm rounded"
            >Remove</a> -->
          </h2>
          <div class="p-4 w-full">
            <div
              class="flex flex-col mb-4 rounded space-y-3"
            >
              <input 
                v-model="form.title"
                class="border border-grey-light flex-1 form-control h-10 input leading-normal px-3 relative rounded text-gray-700 w-full"
                placeholder="title"
              >
              <input 
                v-model="form.text"
                class="border border-grey-light flex-1 form-control h-10 input leading-normal px-3 relative rounded text-gray-700 w-full"
                placeholder="text"
              >
              <input 
                v-model="form.imageSrc"
                class="border border-grey-light flex-1 form-control h-10 input leading-normal px-3 relative rounded text-gray-700 w-full"
                placeholder="Image src"
              >
              <input 
                v-model="form.userName"
                class="border border-grey-light flex-1 form-control h-10 input leading-normal px-3 relative rounded text-gray-700 w-full"
                placeholder="username"
              >
            </div>
         
            <div class="flex flex-col md:flex-row md:justify-end md:space-x-1 md:space-y-0 md:text-right space-y-1">
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