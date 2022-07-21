<script lang="ts" setup>
import { ModerationSettingsDto } from '@tsuwari/shared';
import { toRef } from 'vue';
import { useI18n } from 'vue-i18n';

import Add from '@/assets/buttons/add.svg';
import Remove from '@/assets/buttons/remove.svg';

type Settings = ModerationSettingsDto & {
  checkClips: boolean;
  blackListSentences: string[];
  triggerLength: number;
};
const props = defineProps<{
  settings: Settings;
}>();
const { t } = useI18n({
  useScope: 'global',
});

const settings = toRef(props, 'settings', {
  ...props.settings,
  checkClips: props.settings.checkClips ?? (false as boolean),
  blackListSentences: props.settings.blackListSentences ?? [],
  triggerLength: props.settings.triggerLength ?? 300,
});
</script>

<template>
  <h2 class="border-b border-gray-700 card-title flex font-bold form-switch justify-between outline-none p-2">
    <p>{{ settings.type.charAt(0).toUpperCase() + settings.type.substring(1, settings.type.length) }}</p>
    <input
      :id="'enabledState' + settings.type"
      v-model="settings.enabled"
      class="align-top appearance-none bg-contain bg-gray-300 bg-no-repeat cursor-pointer float-left focus:outline-none form-check-input h-5 rounded-full shadow w-9"
      type="checkbox"
      role="switch"
    >
  </h2>

  <div class="mb-[0.5rem] px-6 py-5 rounded text-base">
    <div class="flex items-center justify-center">
      <div
        class="inline-flex shadow"
        role="group"
      />
    </div>

    <label
      :for="'timeoutMessage' + settings.type"
      class="form-label inline-block mb-1"
    >{{
      t('pages.moderation.timeout.title')
    }}</label>
    <input
      :id="'timeoutMessage' + settings.type"
      v-model="settings.banMessage"
      type="text"
      class="bg-clip-padding bg-white block border border-gray-300 border-solid ease-in-out focus:outline-none font-normal form-control m-0 px-3 py-1.5 rounded text-base text-gray-700 transition w-full"
      :placeholder="t('pages.moderation.timeout.placeholder')"
    >

    <div class="mt-5">
      <label
        :for="'timeoutTime' + settings.type"
        class="form-label inline-block mb-1"
      >{{
        t('pages.moderation.time.title')
      }}</label>
      <input
        :id="'timeoutTime' + settings.type"
        v-model.number="settings.banTime"
        type="number"
        class="bg-clip-padding bg-white block border border-gray-300 border-solid ease-in-out focus:outline-none font-normal form-control m-0 px-3 py-1.5 rounded text-base text-gray-700 transition w-full"
      >
    </div>

    <div class="mt-5">
      <label
        :for="'warningMessage' + settings.type"
        class="form-label inline-block mb-1"
      >Warning message</label>
      <input
        :id="'warningMessage' + settings.type"
        v-model="settings.warningMessage"
        type="text"
        class="bg-clip-padding bg-white block border border-gray-300 border-solid ease-in-out focus:outline-none font-normal form-control m-0 px-3 py-1.5 rounded text-base text-gray-700 transition w-full"
        :placeholder="t('pages.moderation.warning.placeholder')"
      >
    </div>
    <div
      v-if="settings.type === 'links'"
      class="mt-5"
    >
      <div class="flex form-check justify-between">
        <label
          class="form-check-label inline-block"
          for="flexSwitchModClips"
        >{{ t('pages.moderation.clips') }}</label>

        <div class="form-switch">
          <input
            id="flexSwitchModClips"
            v-model="settings.checkClips"
            class="align-top appearance-none bg-contain bg-gray-300 bg-no-repeat cursor-pointer float-left focus:outline-none form-check-input h-5 rounded-full shadow w-9"
            type="checkbox"
            role="switch"
          >
        </div>
      </div>
    </div>
    <div
      v-if="settings.type === 'blacklists'"
      class="mt-5"
    >
      <span class="flex items-center label">
        <span>{{ t('pages.moderation.blacklist') }}</span>
        <span
          class="bg-green-600 cursor-pointer duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-green-700 inline-block leading-tight ml-2 px-1 py-1 rounded shadow text-white text-xs transition uppercase"
          @click="settings.blackListSentences?.push('')"
        ><Add /></span>
      </span>

      <div
        class="gap-2 grid grid-cols-2 input-group lg:grid-cols-2 max-h-40 md:grid-cols-2 mt-1 overflow-auto pr-2 pt-1 scrollbar sm:grid-cols-2 xl:grid-cols-3"
      >
        <div
          v-for="(word, wordIndex) in settings.blackListSentences"
          :key="wordIndex"
          class="flex flex-wrap items-stretch relative"
        >
          <input
            v-model.lazy="settings.blackListSentences[wordIndex]"
            type="text"
            class="border border-grey-light flex-auto flex-grow flex-shrink leading-normal px-3 py-1.5 relative rounded rounded-r-none text-gray-700 w-px"
          >
          <div
            class="-mr-px cursor-pointer flex"
            @click="settings.blackListSentences.splice(wordIndex, 1)"
          >
            <span
              class="bg-red-600 border-0 border-grey-light border-l-0 flex hover:bg-red-700 items-center leading-normal px-5 py-1.5 rounded rounded-l-none text-grey-dark text-sm whitespace-no-wrap"
            ><Remove /></span>
          </div>
        </div>
      </div>
    </div>
    <div
      v-if="settings.type === 'symbols'"
      class="mt-5"
    >
      <label
        :for="'maxSymbols' + settings.type"
        class="form-label inline-block mb-1"
      >{{
        t('pages.moderation.symbols')
      }}</label>
      <input
        :id="'maxSymbols' + settings.type"
        v-model.number="settings.maxPercentage"
        type="text"
        class="bg-clip-padding bg-white block border border-gray-300 border-solid ease-in-out focus:outline-none font-normal form-control m-0 px-3 py-1.5 rounded text-base text-gray-700 transition w-full"
        placeholder="50"
      >
    </div>
    <div
      v-if="settings.type === 'longMessage'"
      class="mt-5"
    >
      <label
        :for="'longMessage' + settings.type"
        class="form-label inline-block mb-1"
      >{{
        t('pages.moderation.lnght')
      }}</label>
      <input
        :id="'longMessage' + settings.type"
        v-model.number="settings.triggerLength"
        type="text"
        class="bg-clip-padding bg-white block border border-gray-300 border-solid ease-in-out focus:outline-none font-normal form-control m-0 px-3 py-1.5 rounded text-base text-gray-700 transition w-full"
        placeholder="50"
      >
    </div>
    <div
      v-if="settings.type === 'caps'"
      class="mt-5"
    >
      <label
        :for="'maxCaps' + settings.type"
        class="form-label inline-block mb-1"
      >{{
        t('pages.moderation.caps')
      }}</label>
      <input
        :id="'maxCaps' + settings.type"
        v-model.number="settings.maxPercentage"
        type="text"
        class="bg-clip-padding bg-white block border border-gray-300 border-solid ease-in-out focus:outline-none font-normal form-control m-0 px-3 py-1.5 rounded text-base text-gray-700 transition w-full"
        placeholder="50"
      >
    </div>
    <div
      v-if="settings.type === 'emotes'"
      class="mt-5"
    >
      <label
        :for="'maxEmotes' + settings.type"
        class="form-label inline-block mb-1"
      >{{
        t('pages.moderation.emotes')
      }}</label>
      <input
        :id="'maxEmotes' + settings.type"
        v-model.number="settings.triggerLength"
        type="text"
        class="bg-clip-padding bg-white block border border-gray-300 border-solid ease-in-out focus:outline-none font-normal form-control m-0 px-3 py-1.5 rounded text-base text-gray-700 transition w-full"
        placeholder="50"
      >
    </div>

    <div class="mt-5">
      <div class="flex form-check justify-between">
        <label
          class="form-check-label inline-block"
          for="flexSwitchSubs"
        >Moderate subscribers</label>

        <div class="form-switch">
          <input
            :id="'moderateSubscribers' + settings.subscribers"
            v-model="settings.subscribers"
            class="align-top appearance-none bg-contain bg-gray-300 bg-no-repeat cursor-pointer float-left focus:outline-none form-check-input h-5 rounded-full shadow w-9"
            type="checkbox"
            role="switch"
          >
        </div>
      </div>

      <div class="flex form-check justify-between">
        <label
          class="form-check-label inline-block"
          for="flexSwitchVip"
        >Moderate vips</label>

        <div class="form-switch">
          <input
            :id="'moderateVips' + settings.vips"
            v-model="settings.vips"
            class="align-top appearance-none bg-contain bg-gray-300 bg-no-repeat cursor-pointer float-left focus:outline-none form-check-input h-5 rounded-full shadow w-9"
            type="checkbox"
            role="switch"
          >
        </div>
      </div>
    </div>
  </div>
</template>0
