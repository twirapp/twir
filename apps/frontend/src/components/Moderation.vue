<script lang="ts" setup>
import { ModerationUpdateDto } from '@tsuwari/shared';
import { toRef } from 'vue';

type Settings = ModerationUpdateDto['items'][0] & { 
  checkClips: boolean,
  blackListSentences: string[],
  triggerLength: number,
}
const props = defineProps<{
  settings: Settings
}>();


const settings = toRef(props, 'settings', {
  ...props.settings,
  checkClips: props.settings.checkClips ?? false as boolean,
  blackListSentences: props.settings.blackListSentences ?? [],
  triggerLength: props.settings.triggerLength ?? 300,
});
</script>

<template>
  <h2 class="card-title p-2 flex justify-between border-b border-gray-700 outline-none">
    <p>{{ settings.type.charAt(0).toUpperCase() + settings.type.substring(1, settings.type.length) }}</p>
  </h2>
  <div class="p-0">
    <div
      class="rounded-lg py-5 px-6 mb-4 text-base mb-4"
    >
      <div class="flex items-center justify-center">
        <div
          class="inline-flex shadow-md hover:shadow-lg focus:shadow-lg"
          role="group"
        >
          <button
            type="button"
            class="rounded-l inline-block px-6 py-2.5 text-white font-medium text-xs leading-tight uppercase focus:outline-none focus:ring-0 transition duration-150 ease-in-out"
            :class="[settings.enabled ? 'bg-green-500 hover:bg-green-600' : 'bg-red-500 hover:bg-red-600' ]"
            @click="() => settings.enabled = !settings.enabled"
          >
            {{ settings.enabled ? 'Enabled' : 'Disabled' }}
          </button>
          <button
            type="button"
            class="inline-block px-6 py-2.5 text-white font-medium text-xs leading-tight uppercase focus:outline-none focus:ring-0 transition duration-150 ease-in-out"
            :class="[settings.subscribers ? 'bg-green-500 hover:bg-green-600' : 'bg-red-500 hover:bg-red-600' ]"
            @click="() => settings.subscribers = !settings.subscribers"
          >
            Moderate Subscribers
          </button>
          <button
            type="button"
            class="inline-block px-6 py-2.5 text-white font-medium text-xs leading-tight uppercase focus:outline-none focus:ring-0 transition duration-150 ease-in-out rounded-r"
            :class="[settings.vips ? 'bg-green-500 hover:bg-green-600' : 'bg-red-500 hover:bg-red-600' ]"
            @click="() => settings.vips = !settings.vips"
          >
            Moderate vips
          </button>
        </div>
      </div>
      <div class="mt-3">
        <label
          :for="'timeoutMessage' + settings.type"
          class="form-label inline-block mb-2"
        >Timeout message</label>
        <input
          :id="'timeoutMessage' + settings.type"
          v-model="settings.banMessage"
          type="text"
          class="
            form-control
            block
            w-full
            px-3
            py-1.5
            text-base
            font-normal
            text-gray-700
            bg-white bg-clip-padding
            border border-solid border-gray-300
            rounded
            transition
            ease-in-out
            m-0
            focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none
          "
          placeholder="Message for timeout"
        >
      </div>
      <div class="mt-3">
        <label
          :for="'timeoutTime' + settings.type"
          class="form-label inline-block mb-2"
        >Timeout time</label>
        <input
          :id="'timeoutTime' + settings.type"
          v-model="settings.banTime"
          type="number"
          class="
            form-control
            block
            w-full
            px-3
            py-1.5
            text-base
            font-normal
            text-gray-700
            bg-white bg-clip-padding
            border border-solid border-gray-300
            rounded
            transition
            ease-in-out
            m-0
            focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none
          "
        >
      </div>
      <div class="mt-3">
        <label
          :for="'warningMessage' + settings.type"
          class="form-label inline-block mb-2"
        >Warning message</label>
        <input
          :id="'warningMessage' + settings.type"
          v-model="settings.warningMessage"
          type="text"
          class="
            form-control
            block
            w-full
            px-3
            py-1.5
            text-base
            font-normal
            text-gray-700
            bg-white bg-clip-padding
            border border-solid border-gray-300
            rounded
            transition
            ease-in-out
            m-0
            focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none
          "
          placeholder="Message for warning"
        >
      </div>
      <div
        v-if="settings.type === 'links'"
        class="mt-3"
      >
        <div class="form-check form-switch">
          <label
            class="form-check-label inline-block"
            for="flexSwitchCheckChecked"
          >Moderate clips</label>
          <input
            id="flexSwitchCheckChecked"
            v-model="settings.checkClips"
            class="form-check-input appearance-none w-9 -ml-10 rounded-full float-left h-5 align-top bg-white bg-no-repeat bg-contain bg-gray-300 focus:outline-none cursor-pointer shadow-sm"
            type="checkbox"
            role="switch"
          >
        </div>
      </div>
      <div
        v-if="settings.type === 'blacklists'"
        class="mt-3"
      >
        <span class="label">  
          <span>Blacklisted words
            <a
              class="px-2 py-0.5 inline-block bg-green-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md  hover:shadow-lg focus:shadow-lg focus:outline-none focus:ring-0 active:bg-purple-800 active:shadow-lg transition duration-150 ease-in-out cursor-pointer"
              @click="settings.blackListSentences?.push('')"
            >
              +
            </a>
          </span>
        </span>
        <div class="input-group pt-1 pr-2 grid lg:grid-cols-2 md:grid-cols-2 sm:grid-cols-2 grid-cols-2 xl:grid-cols-3 gap-1 max-h-40 scrollbar-thin overflow-auto scrollbar scrollbar-thumb-gray-900 scrollbar-track-gray-600">
          <div
            v-for="word, wordIndex in settings.blackListSentences"
            :key="wordIndex"
            class="flex flex-wrap items-stretch mb-4 relative"
          >
            <input
              v-model.lazy="settings.blackListSentences[wordIndex]"
              type="text"
              class="flex-shrink flex-grow flex-auto leading-normal w-px flex-1 border border-grey-light text-gray-700 rounded px-3 py-1.5 relative rounded-r-none"
            >
            <div
              class="flex -mr-px cursor-pointer"
              @click="settings.blackListSentences.splice(wordIndex, 1)"
            >
              <span class="flex items-center leading-normal bg-red-500 rounded rounded-l-none border-0 border-l-0 border-grey-light px-5 py-1.5 whitespace-no-wrap text-grey-dark text-sm">X</span>
            </div>
          </div>
        </div>
      </div>
      <div
        v-if="settings.type === 'symbols'"
        class="mt-3"
      >
        <label
          :for="'maxSymbols' + settings.type"
          class="form-label inline-block mb-2"
        >Max symbols in message (in percent)</label>
        <input
          :id="'maxSymbols' + settings.type"
          v-model="settings.maxPercentage"
          type="text"
          class="
            form-control
            block
            w-full
            px-3
            py-1.5
            text-base
            font-normal
            text-gray-700
            bg-white bg-clip-padding
            border border-solid border-gray-300
            rounded
            transition
            ease-in-out
            m-0
            focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none
          "
          placeholder="50"
        >
      </div>
      <div
        v-if="settings.type === 'longMessage'"
        class="mt-3"
      >
        <label
          :for="'longMessage' + settings.type"
          class="form-label inline-block mb-2"
        >Max message length</label>
        <input
          :id="'longMessage' + settings.type"
          v-model="settings.triggerLength"
          type="text"
          class="
            form-control
            block
            w-full
            px-3
            py-1.5
            text-base
            font-normal
            text-gray-700
            bg-white bg-clip-padding
            border border-solid border-gray-300
            rounded
            transition
            ease-in-out
            m-0
            focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none
          "
          placeholder="50"
        >
      </div>
      <div
        v-if="settings.type === 'caps'"
        class="mt-3"
      >
        <label
          :for="'maxCaps' + settings.type"
          class="form-label inline-block mb-2"
        >Max caps in message (in percent)</label>
        <input
          :id="'maxCaps' + settings.type"
          v-model="settings.maxPercentage"
          type="text"
          class="
            form-control
            block
            w-full
            px-3
            py-1.5
            text-base
            font-normal
            text-gray-700
            bg-white bg-clip-padding
            border border-solid border-gray-300
            rounded
            transition
            ease-in-out
            m-0
            focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none
          "
          placeholder="50"
        >
      </div>
      <div
        v-if="settings.type === 'emotes'"
        class="mt-3"
      >
        <label
          :for="'maxEmotes' + settings.type"
          class="form-label inline-block mb-2"
        >Max emotes in message.</label>
        <input
          :id="'maxEmotes' + settings.type"
          v-model="settings.triggerLength"
          type="text"
          class="
            form-control
            block
            w-full
            px-3
            py-1.5
            text-base
            font-normal
            text-gray-700
            bg-white bg-clip-padding
            border border-solid border-gray-300
            rounded
            transition
            ease-in-out
            m-0
            focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none
          "
          placeholder="50"
        >
      </div>
    </div>
  </div>
</template>