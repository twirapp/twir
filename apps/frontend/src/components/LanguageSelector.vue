<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

import { localeStore } from '@/stores/locale';

const i18n = useI18n();

function setLocale(v: string) {
  i18n.locale.value = v;
  localeStore.set(v);
}
</script>

<template>
  <div class="flex items-center justify-center">
    <div>
      <div class="dropdown relative select-none">
        <div
          id="langSelector"
          class="hover:text-slate-300         flex
          items-center space-x-2 dropdown-toggle
           rounded px-1 cursor-pointer"
          data-bs-toggle="dropdown"
          aria-expanded="false"
        >
          <span
            class="fi rounded-sm"
            :class="`fi-${$i18n.locale}`"
          />
          <p>
            {{ $t("name", $i18n.locale) }}
          </p>
          <svg
            aria-hidden="true"
            focusable="false"
            data-prefix="fas"
            data-icon="caret-down"
            class="w-2 ml-2"
            role="img"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 320 512"
          >
            <path
              fill="currentColor"
              d="M31.3 192h257.3c17.8 0 26.7 21.5 14.1 34.1L174.1 354.8c-7.8 7.8-20.5 7.8-28.3 0L17.2 226.1C4.6 213.5 13.5 192 31.3 192z"
            />
          </svg>
        </div>

        <ul
          class="
          dropdown-menu
          absolute
          hidden
bg-[#202020]
          text-base
          z-50
          float-left
          py-1
          list-none
          text-left
          rounded
          mt-1
          m-0
          bg-clip-padding
          border-none
          space-y-0.5
          w-max
          space-y-0.5
          max-h-[55vh]
                   scrollbar-thin overflow-auto scrollbar scrollbar-thumb-gray-600 scrollbar-track-gray-500
        "
          aria-labelledby="langSelector"
        >
          <div
            v-for="(lang) in $i18n.availableLocales"
            :key="lang"
            class="flex px-1 hover:bg-[#393636] hover:rounded space-x-2 items-center"
            @click="setLocale(lang)"
          >
            <div class="flex items-center space-x-2">
              <span
                class="fi rounded-sm"
                :class="`fi-${lang}`"
              />
              <p>{{ $t('name', 'unknown', { locale: lang }) }}</p>
            </div>
          </div>
        </ul>
      </div>
    </div>
  </div>
</template>