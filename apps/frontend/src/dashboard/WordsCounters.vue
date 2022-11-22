<script lang="ts" setup>
export type WordsCounter = {
  id?: string;
  phrase: string;
  counter: number;
  enabled: boolean;
};
export type CounterType = WordsCounter & { edit?: boolean };

import { Ref, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import WordsCounterComponent from '@/components/WordsCounter.vue';
import { useUpdatingData } from '@/functions/useUpdatingData';

const { data } = useUpdatingData(`/v1/channels/{dashboardId}/words_counters`);
const counters = ref<Array<CounterType>>([]);
const countersBeforeEdit = ref<Array<CounterType>>([]);
const { t } = useI18n();

watch(data, (v: any[]) => {
  counters.value = v;
});

function insert() {
  counters.value = [
    {
      phrase: '',
      counter: 0,
      edit: true,
      enabled: true,
    },
    ...counters.value,
  ];
}

async function deletecounter(index: number) {
  counters.value = counters.value.filter((_, i) => i !== index);
}

function cancelEdit(counter: Ref<CounterType>) {
  const index = counters.value.indexOf(counter.value);
  if (counter.value.id && counters.value) {
    const editableCounters = countersBeforeEdit.value?.find((c) => c.id === counter.value.id);
    if (editableCounters) {
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      counter.value[index] = {
        ...editableCounters,
        edit: false,
      };

      countersBeforeEdit.value = countersBeforeEdit.value.filter(
        (v, i) => i !== countersBeforeEdit.value.indexOf(editableCounters),
      );
    }
  } else {
    counters.value = counters.value.filter((v, i) => i !== index);
  }
}
</script>

<template>
  <div class="m-1.5 md:m-3">
    <div class="flow-root">
      <div class="btn btn-primary btn-sm float-left mb-5 md:w-auto rounded w-full">
        <button
          class="
            bg-purple-600
            duration-150
            ease-in-out
            focus:outline-none
            focus:ring-0
            font-medium
            hover:bg-purple-700
            inline-block
            leading-tight
            px-6
            py-2.5
            rounded
            shadow
            text-white text-xs
            transition
            uppercase
          "
          @click="insert"
        >
          {{ t('buttons.addNew') }}
        </button>
      </div>
    </div>

    <masonry-wall :items="counters" :gap="8">
      <template #default="{ item, index }">
        <div :key="index" class="block card rounded shadow text-white">
          <WordsCounterComponent
            :counter="item"
            :counters="counters"
            :counters-before-edit="countersBeforeEdit"
            @delete="deletecounter"
            @cancel-edit="cancelEdit"
          />
        </div>
      </template>
    </masonry-wall>
  </div>
</template>
