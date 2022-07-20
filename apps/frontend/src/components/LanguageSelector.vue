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
          class="cursor-pointer dropdown-toggle flex hover:text-slate-300 items-center rounded space-x-2"
          data-bs-toggle="dropdown"
          aria-expanded="false"
        >
          <span
            class="fi rounded-sm"
            :class="`fi-${$i18n.locale}`"
          />
          <!-- <p class="hidden md:block">
            {{ $t("name", $i18n.locale) }}
          </p> -->
          <svg
            aria-hidden="true"
            focusable="false"
            data-prefix="fas"
            data-icon="caret-down"
            class="ml-2 w-2"
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
          class="absolute
          bg-[#202020]
          bg-clip-padding
          border-none
          dropdown-menu
          float-left
          hidden
          list-none
          m-0
          max-h-[55vh]
          mt-1
          mx-2
          overflow-auto
          py-1
          rounded
          scrollbar
          space-y-1
          text-base
          text-left
          w-max
          z-50"
          aria-labelledby="langSelector"
        >
          <div
            v-for="(lang) in $i18n.availableLocales"
            :key="lang"
            class="flex hover:bg-[#393636] hover:rounded items-center mx-1 px-2"
            @click="setLocale(lang)"
          >
            <div class="flex items-center space-x-2.5">
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