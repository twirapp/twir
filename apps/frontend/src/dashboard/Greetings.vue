<script lang="ts" setup>
export type GreeTingType = SetOptional<
  Omit<Greeting, 'channelId'> & { userName: string; edit?: boolean },
  'id'
>;

import { Greeting } from '@tsuwari/prisma';
import type { SetOptional } from 'type-fest';
import { Ref, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import GreetingComponent from '@/components/Greeting.vue';
import { useUpdatingData } from '@/functions/useUpdatingData';

const { t } = useI18n({
  useScope: 'global',
});

const { data } = useUpdatingData(`/v1/channels/{dashboardId}/greetings`);
const greetings = ref<Array<GreeTingType>>([]);
const greetingsBeforeEdit = ref<Array<GreeTingType>>([]);

watch(data, (v: any[]) => {
  greetings.value = v;
  greetingsBeforeEdit.value = [];
});

function insert() {
  greetings.value = [
    {
      userName: '',
      userId: '',
      text: '',
      edit: true,
      enabled: true,
      isReply: true,
    },
    ...greetings.value,
  ];
}

async function deleteGreeting(index: number) {
  greetings.value = greetings.value.filter((_, i) => i !== index);
}

function cancelEdit(greeting: Ref<GreeTingType>) {
  const index = greetings.value.indexOf(greeting.value);
  if (greeting.value.id && greetings.value) {
    const editableCommand = greetingsBeforeEdit.value?.find((c) => c.id === greeting.value.id);
    if (editableCommand) {
      greetings.value[index] = {
        ...editableCommand,
        edit: false,
      };

      greetingsBeforeEdit.value = greetingsBeforeEdit.value.filter(
        (v, i) => i !== greetingsBeforeEdit.value.indexOf(editableCommand),
      );
    }
  } else {
    greetings.value = greetings.value.filter((v, i) => i !== index);
  }
}
</script>

<template>
  <div class="m-1.5 md:m-3">
    <div class="flow-root">
      <div class="btn btn-primary btn-sm float-left mb-5 md:w-auto rounded w-full">
        <button
          class="bg-purple-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-purple-700 inline-block leading-tight px-6 py-2.5 rounded shadow text-white text-xs transition uppercase"
          @click="insert">
          {{ t('pages.greetings.buttons.add') }}
        </button>
      </div>
    </div>

    <masonry-wall :items="greetings" :gap="8">
      <template #default="{ item, index }">
        <div :key="index" class="block card rounded shadow text-white">
          <GreetingComponent
            :greeting="item"
            :greetings="greetings"
            :greetings-before-edit="greetingsBeforeEdit"
            @delete="deleteGreeting"
            @cancel-edit="cancelEdit" />
        </div>
      </template>
    </masonry-wall>
  </div>
</template>
