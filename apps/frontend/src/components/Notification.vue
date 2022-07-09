<script lang="ts" setup>
 import { Popover, PopoverButton, PopoverPanel } from '@headlessui/vue';
import { useStore } from '@nanostores/vue';
import { Notification, ViewedNotification, NotificationMessage } from '@tsuwari/prisma';
import { useAxios } from '@vueuse/integrations/useAxios';
import { computed, ComputedRef, Ref, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';


import Bell from '@/assets/icons/bell.svg?component';
import Mark from '@/assets/icons/check.svg?component';
import MyBtn from '@/components/elements/MyBtn.vue';
import { api } from '@/plugins/api';
import { localeStore } from '@/stores/locale';
import { selectedDashboardStore, userStore } from '@/stores/userStore';

const selectedDashboard = useStore(selectedDashboardStore);
const user = useStore(userStore);

const { execute: executeNew, data: newNotificationsData } = useAxios(`/v1/channels/${selectedDashboard.value.channelId}/notifications/new`, api, { immediate: false });
const { execute: executeViewed, data: viewedNotificationsData } = useAxios(`/v1/channels/${selectedDashboard.value.channelId}/notifications/viewed`, api, { immediate: false });

type Viewed = ViewedNotification & { notification: Notification & { messages: NotificationMessage[] } }
type NotViewed = Notification & { messages: NotificationMessage[] }

const newNotifications = ref<NotViewed[]>([]);
const viewedNotifications = ref<Viewed[]>([]);
const showNew = ref(true);
const { t } = useI18n();
const selectedLang = useStore(localeStore);

const notifications: ComputedRef<Ref<NotViewed[]> | Ref<Viewed[]>> = computed(() => {
  return showNew.value ? newNotifications : viewedNotifications;
});

async function markNotificationAsReaded(notification: Viewed | NotViewed) {
  await api.post(`v1/channels/${selectedDashboard.value.channelId}/notifications/viewed`, {
    notificationId: notification.id,
  });

  newNotifications.value = newNotifications.value.filter(v => v.id !== notification.id);
  executeViewed(`/v1/channels/${selectedDashboard.value.channelId}/notifications/viewed`);
}

watch(newNotificationsData, (v: any[]) => {
  newNotifications.value = v;
});

watch(viewedNotificationsData, (v) => {
  viewedNotifications.value = v;
});

selectedDashboardStore.subscribe(async (v) => {
  executeNew(`/v1/channels/${v.channelId}/notifications/new`);
  executeViewed(`/v1/channels/${v.channelId}/notifications/viewed`);
});
</script>

<template>
  <div class="block inline-flex items-center relative">
    <Popover>
      <div
        v-if="newNotifications.length"
        class="-translate-y-1/2 0 absolute align-baseline bg-[#772CE8] bottom-auto font-bold inline-block leading-none left-auto px-1.5 py-0.5 rotate-0 rounded scale-x-100 scale-y-100 select-none skew-x-0 skew-y-0 text-center text-white text-xs translate-x-2/4 whitespace-nowrap z-10"
      >
        {{ newNotifications.length }}
      </div>
      <PopoverButton
        class="hover:text-slate-300"
      >
        <Bell class="inline" />
      </PopoverButton>
      
      <div
        class="-translate-x-3/4 absolute float-left mt-3 w-96 z-10"
      >
        <PopoverPanel
          class="bg-[#202020] max-h-[55vh] overflow-auto scrollbar scrollbar-thin scrollbar-thumb-gray-600 scrollbar-track-gray-500 z-10"
        >
          <div
            class="flex items-center justify-between mt-1 px-2 space-y-1"
          >
            <div class="font-bold">
              {{ t("navbar.notifications.title") }}
            </div>
            <div>
              <MyBtn
                color="purple"
                size="small"
                @click="showNew = !showNew"
              >
                {{ t(showNew ? `navbar.notifications.showOld` : `navbar.notifications.showNew`) }}
              </MyBtn>
              <!-- <button
                type="button"
                class="bg-[#522f87] duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-[#772CE8] inline-block leading-tight px-1 py-1 rounded shadow text-white text-xs transition uppercase"
                @click="showNew = !showNew"
              >
                {{ t(showNew ? `navbar.notifications.showOld` : `navbar.notifications.showNew`) }}
              </button> -->
            </div>
          </div>
          <div
            v-if="!notifications.value?.length"
            class="mx-3 my-3 text-center"
          >
            {{ t('navbar.notifications.noUnreadNotifications') }}
          </div>
          <!-- v-element-visibility="[(state) => onNotificationVisibility(notification.text, state)]" -->
          <div
            v-for="notification of notifications.value"
            v-else
            :key="notification.id"
            class="block font-normal hover:bg-[#393636] mt-1 px-2 py-2 text-sm w-full"
          >
            {{ ('messages' in notification ? notification?.messages : notification.notification.messages).find(m => m.langCode === selectedLang.toUpperCase())?.text }}
            <div
              v-if="('messages' in notification) && selectedDashboard.channelId === user?.id"
              class="flex flex-col md:flex-row md:justify-end md:space-x-1 md:space-y-0 md:text-right mt-1 pr-2 space-y-1"
            >
              <button
                type="button"
                class="duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-[#772CE8] inline-block leading-tight px-1 py-1 rounded shadow text-white text-xs transition uppercase"
                @click="markNotificationAsReaded(notification)"
              >
                <Mark style="width: 16px; height: 16px;" />
              </button>
            </div>
          </div>
        </PopoverPanel>
      </div>
    </Popover>
  </div>
</template>