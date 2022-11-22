<template>
  <button
    class="youtube-player-btn-icon"
    @click="openModal"
  >
    <div
      v-if="isSettingsFetching"
      class="animate-spin border-2 border-[#AFAFAF] border-r-transparent border-solid h-5 inline-block rounded-full w-5"
      role="status"
    />
    <Settings
      v-else
      class="stroke-icon"
    />
  </button>
  <TransitionRoot
    appear
    :show="isModalOpen"
    as="template"
  >
    <Dialog
      as="div"
      class="relative z-10"
      @close="closeModal"
    >
      <TransitionChild
        as="template"
        enter="duration-300 ease-out"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="duration-200 ease-in"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <div class="bg-black bg-opacity-40 fixed inset-0" />
      </TransitionChild>

      <div class="fixed inset-0 overflow-y-auto">
        <div class="flex items-center justify-center min-h-full p-4 text-center">
          <TransitionChild
            as="template"
            enter="duration-300 ease-out"
            enter-from="opacity-0 scale-95"
            enter-to="opacity-100 scale-100"
            leave="duration-200 ease-in"
            leave-from="opacity-100 scale-100"
            leave-to="opacity-0 scale-95"
          >
            <DialogPanel
              class="align-middle bg-[#2C2C2C] max-w-2xl overflow-hidden rounded-lg shadow-xl text-left transform transition-all w-full"
            >
              <div class="border-[#525252] border-b flex items-center justify-between md:p-6 p-5">
                <DialogTitle
                  as="h3"
                  class="font-semibold leading-6 text-white text-xl"
                >
                  Youtube song request settings
                </DialogTitle>
                <button
                  @click="closeModal"
                >
                  <Cross class="h-6 hover:stroke-[#D0D0D0] stroke-[#AFAFAF] w-6" />
                </button>
              </div>
              
              <form
                v-if="isDataFetched"
                class=""
                @submit.prevent="submitForm"
              >
                <div class="divide-[#525252] divide-y grid md:px-6 px-5">
                  <div class="py-5">
                    <span class="inline-block leading-tight mb-4 text-white">General</span>
                    <div class="flex flex-col gap-5 justify-between sm:flex-row w-full">
                      <TswSwitch
                        id="accept-only-when-online"
                        name="acceptOnlyWhenOnline"
                        label="Accept only when online"
                        direction="col"
                      />
                      <TswTextInput
                        id="channel-points-reward-name"
                        name="channelPointsRewardName"
                        direction="col"
                        class="flex-1"
                        label="Channel points reward name"
                      />
                    </div>
                  </div>
                  
                  <div class="py-5">
                    <span class="inline-block leading-tight mb-4 text-white">User</span>
                    <div class="gap-3 grid grid-cols-1 md:grid-cols-4">
                      <TswNumberInput
                        id="user.max-requests"
                        name="user.maxRequests"
                        label="Max requests"
                        direction="col"
                      />
                      <TswNumberInput
                        id="user.min-watch-time"
                        name="user.minWatchTime"
                        label="Min watch time"
                        direction="col"
                      />
                      <TswNumberInput
                        id="user.min-follow-time"
                        name="user.minFollowTime"
                        label="Min follow time"
                        direction="col"
                      />
                      <TswNumberInput
                        id="user.min-messages"
                        name="user.minMessages"
                        label="Min messages"
                        direction="col"
                      />
                    </div>
                  </div>
                
                  <div class="py-5">
                    <span class="inline-block leading-tight mb-4 text-white">Song</span>
                    <div class="gap-3 grid grid-flow-col">
                      <TswNumberInput
                        id="song.max-length"
                        name="song.maxLength"
                        label="Max length"
                        direction="col"
                      />
                      <TswNumberInput
                        id="song.min-views"
                        name="song.minViews"
                        label="Min views"
                        direction="col"
                      />
                    </div>
                    <TswArrayInput
                      name="song.acceptedCategories"
                      label="Accepted categories"
                      class="mt-2"
                    />
                  </div>
                  <div class="py-5">
                    <span class="inline-block leading-tight mb-4 text-white">Black list</span>
                    <div class="gap-y-2 grid">
                      <TswArrayInput
                        name="blacklist.usersIds"
                        label="User ids list"
                      />
                      <TswArrayInput
                        name="blacklist.songsIds"
                        label="Song ids list"
                      />
                      <TswArrayInput
                        name="blacklist.channelsIds"
                        label="Channel ids list"
                      />
                      <TswArrayInput
                        name="blacklist.artistsNames"
                        label="Artists names"
                      />
                    </div>
                  </div>
                </div>
                
                
                
                <div class="border-[#525252] border-t gap-2 grid-flow-col inline-grid justify-end md:px-6 md:py-5 p-5 w-full">
                  <button
                    type="submit"
                    class="bg-[#644EE8] focus:outline-none font-medium inline-flex items-center justify-center px-4 py-2 rounded text-sm text-white"
                  >
                    <template v-if="isPostingData">
                      Saving
                      <div
                        class="animate-spin border-[1.5px] border-r-transparent border-solid border-white h-[14px] inline-block ml-[6px] rounded-full w-[14px]"
                        role="status"
                      />
                    </template>
                    <template v-else>
                      Save
                    </template>
                  </button>
                </div>
              </form>
            </DialogPanel>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { V1 } from '@tsuwari/types/api';
import { AxiosError } from 'axios';
import { useForm, configure } from 'vee-validate';
import { computed, ref } from 'vue';
import { object, string, boolean, number, array } from 'yup';

import Cross from '@/assets/icons/cross.svg?component';
import Settings from '@/assets/icons/settings.svg?component';
import TswArrayInput from '@/components/elements/TswArrayInput.vue';
import TswNumberInput from '@/components/elements/TswNumberInput.vue';
import TswSwitch from '@/components/elements/TswSwitch.vue';
import TswTextInput from '@/components/elements/TswTextInput.vue';
import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore.js';

configure({
  validateOnBlur: true,
});

const isModalOpen = ref(false);
const isDataFetched = ref(false);
const isSettingsFetching = ref(false);
const error = ref<string | null>(null);
const isPostingData = ref<boolean>(false);

const schema = computed(() =>
  object({
    user: object({
      maxRequests: number().typeError('Must be a nubmer').min(0, 'Must be greater or equal to zero'),
      minWatchTime: number().typeError('Must be a nubmer').min(0, 'Must be greater or equal to zero'),
      minMessages: number().typeError('Must be a nubmer').min(0, 'Must be greater or equal to zero'),
      minFollowTime: number().typeError('Must be a nubmer').min(0, 'Must be greater or equal to zero'),
    }),
    maxRequests: number().typeError('Must be a nubmer').min(0, 'Must be greater or equal to zero'),
    acceptOnlyWhenOnline: boolean().required(),
    channelPointsRewardName: string().optional(),
    song: object({
      maxLength: number().typeError('Must be a nubmer').min(0, 'Must be greater or equal to zero'),
      minViews: number().typeError('Must be a nubmer').min(0, 'Must be greater or equal to zero'),
      acceptedCategories: array().of(string()),
    }),
    blackList: object({
      usersIds: array().of(string()),
      songsIds: array().of(string()),
      channelsIds: array().of(string()),
      artistsNames: array().of(string()),
    }),
  }),
);

function closeModal() {
  isModalOpen.value = false;
}

const { values, validate, setValues, meta } = useForm<Required<V1['CHANNELS']['MODULES']['YouTube']['POST']>>({
  validationSchema: schema,
  keepValuesOnUnmount: true,
});

const submitForm = async () => {
  if (!meta.value.touched) {
    return console.log('change form to send it');
  }
  const validationResult = await validate();
  if (validationResult.valid) {
    try {
      isPostingData.value = true;
      await api.post(
        `/v1/channels/${selectedDashboardStore.get().channelId}/modules/youtube-sr`,
        values,
      );
    } catch (err) {
      error.value = (err as AxiosError).message;
    } finally {
      isPostingData.value = false;
    }
  } else {
    console.log('For not valid');
  }
  meta.value.touched = false;
};

async function openModal() {
  if (!isDataFetched.value) {
    isSettingsFetching.value = true;
    try {
      const response = await api.get<V1['CHANNELS']['MODULES']['YouTube']['GET']>(
        `/v1/channels/${selectedDashboardStore.get().channelId}/modules/youtube-sr`,
      );
      setValues(response.data);
      isDataFetched.value = true;
      isModalOpen.value = true;
      meta.value.touched = false;
    } catch (err) {
      error.value = (err as AxiosError).message;
      console.error(err);
    } finally {
      isSettingsFetching.value = false;
    }
  } else {
    isModalOpen.value = true;
  }
}
</script>
