<script lang="ts" setup>
import { useTitle } from '@vueuse/core';
import { useI18n } from 'vue-i18n';
import { useRoute } from 'vue-router';

import { publicRoutes, adminRoutes } from './sidebarLinks';

const currentRoute = useRoute();
const { t } = useI18n({
  useScope: 'global',
  inheritLocale: true,
});
const title = useTitle();

</script>

<template>
  <aside class="border-r border-stone-700 h-screen hidden min-w-max overflow-auto scrollbar select-none shadow sm:block w-54">
    <div class="flex items-center justify-center py-1 sidebar-header">
      <div>
        <router-link
          to="/"
          class="flex-row font-bold inline-flex items-center mt-5 text-xl"
        >
          <svg
            width="30"
            height="29"
            viewBox="0 0 20 19"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            class="mr-2"
          >
            <path
              d="M11.22 2.84008H13.8V5.18008C13.8 6.50008 13.98 6.94008 14.36 7.32008C14.72 7.68008 15.32 7.84008 15.86 7.84008H17.06C17.44 7.84008 17.96 7.76008 18.26 7.60008C18.64 7.40008 18.88 7.14008 19.02 6.74008C19.18 6.34008 19.28 5.42008 19.32 4.58008C18.76 4.38008 17.96 4.00008 17.54 3.64008C17.52 4.38008 17.5 5.02008 17.46 5.30008C17.42 5.58008 17.34 5.70008 17.26 5.74008C17.2 5.78008 17.06 5.80008 16.96 5.80008H16.54C16.42 5.80008 16.32 5.78008 16.26 5.72008C16.2 5.64008 16.2 5.46008 16.2 5.10008V0.760078H8.95999V2.82008C8.95999 4.18008 8.73999 5.92008 6.93999 7.18008C7.45999 7.44008 8.49999 8.12008 8.89999 8.52008C10.86 7.04008 11.22 4.72008 11.22 2.88008V2.84008ZM4.81999 0.0800781C3.95999 1.42008 2.15999 3.06008 0.599994 4.04008C0.979994 4.54008 1.51999 5.52008 1.77999 6.04008C3.67999 4.82008 5.73999 2.84008 7.05999 0.980077L4.81999 0.0800781ZM5.43999 4.26008C4.23999 6.24008 2.23999 8.24008 0.399994 9.50008C0.799994 10.0201 1.43999 11.2601 1.65999 11.7801C2.23999 11.3401 2.81999 10.8201 3.41999 10.2401V18.7201H5.71999V7.72008C6.39999 6.88008 7.03999 6.00008 7.55999 5.14008L5.43999 4.26008ZM14.9 10.8001C14.36 11.8201 13.64 12.7201 12.8 13.5201C11.96 12.7201 11.28 11.8201 10.78 10.8001H14.9ZM16.54 8.56008L16.1 8.66008H7.57999V10.8001H10.4L8.61999 11.3401C9.25999 12.7001 10.04 13.8801 11 14.9001C9.63999 15.7601 8.11999 16.4001 6.47999 16.8001C6.91999 17.2801 7.47999 18.1801 7.77999 18.8001C9.59999 18.2601 11.3 17.5001 12.8 16.5001C14.24 17.5401 15.98 18.3201 18.04 18.8001C18.38 18.1601 19.02 17.2001 19.56 16.6801C17.7 16.3401 16.06 15.7401 14.7 14.9401C16.22 13.4801 17.4 11.6201 18.14 9.28008L16.54 8.56008Z"
              fill="white"
            />
          </svg>
          Tsuwari
        </router-link>
      </div>
    </div>

    <ul class="mt-3 px-0.5 relative">
      <li
        v-for="(route, index) in currentRoute.fullPath.includes('/admin') ? adminRoutes : publicRoutes"
        :key="index"
      >
        <RouterLink
          :to="route.path"
          class="border-slate-300 duration-300 ease-in-out flex h-12 hover:bg-[#202122] items-center mt-1 overflow-hidden px-6 py-4 ripple-surface-primary rounded text-ellipsis text-sm text-white transition whitespace-nowrap"
          :class="{
            'bg-neutral-700': currentRoute.path === route.path,
          }"
          @click="title = `Tsuwari - ${route.name.charAt(0).toUpperCase() + route.name.substring(1)}`"
        >
          <span
            v-if="route.icon"
            class="h-3 mr-3 w-3"
          >
            <img :src="route.icon">
          </span>
             
          <span>{{ t(`pages.${route.name}.sidebarName`) }}</span>
        </RouterLink>
      </li>
    </ul>
  </aside>
</template>
