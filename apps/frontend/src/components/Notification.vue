<!-- eslint-disable vue/no-v-html -->
<script lang="ts" setup>
import { Popover, PopoverButton, PopoverPanel } from '@headlessui/vue';
import { useStore } from '@nanostores/vue';
import { Notification, NotificationMessage } from '@tsuwari/prisma';
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';


import Bell from '@/assets/icons/bell.svg?component';
import Mark from '@/assets/icons/check.svg?component';
import MyBtn from '@/components/elements/MyBtn.vue';
import { useUpdatingData } from '@/functions/useUpdatingData';
import { api } from '@/plugins/api';
import { localeStore } from '@/stores/locale';
import { selectedDashboardStore, userStore } from '@/stores/userStore';

const selectedDashboard = useStore(selectedDashboardStore);
const user = useStore(userStore);

const { data: newNotifications } = useUpdatingData<MyNotification[]>(`/v1/channels/{dashboardId}/notifications/new`);
const { execute: executeViewed, data: viewedNotificationsData } = useUpdatingData(`/v1/channels/{dashboardId}/notifications/viewed`);

type MyNotification = Notification & { messages: NotificationMessage[] }

const viewedNotifications = ref<MyNotification[]>([]);
const showNew = ref(true);
const { t } = useI18n();
const selectedLang = useStore(localeStore);

const notifications = computed(() => {
  return showNew.value ? newNotifications : viewedNotifications;
});

async function markNotificationAsReaded(notification: MyNotification) {
  await api.post(`v1/channels/${selectedDashboard.value.channelId}/notifications/viewed`, {
    notificationId: notification.id,
  });

  newNotifications.value = newNotifications?.value?.filter(v => v.id !== notification.id);
  executeViewed();
}

watch(viewedNotificationsData, (v) => {
  viewedNotifications.value = v.map((v: any) => v.notification);
});
</script>

<template>
  <div class="block inline-flex items-center relative">
    <Popover>
      <div
        v-if="newNotifications?.length"
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
        class="-translate-x-56 absolute float-left md:-translate-x-3/4 md:w-96 mt-3 w-80 z-10"
      >
        <PopoverPanel
          class="bg-[#202020] break-words max-h-[55vh] overflow-auto pr-2 scrollbar z-10"
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
            </div>
          </div>
          <div
            v-if="!notifications.value?.length"
            class="mx-3 my-3 text-center"
          >
            {{ t('navbar.notifications.noUnreadNotifications') }}
          </div>
          <div
            v-for="notification of notifications.value"
            v-else
            :key="notification.id"
            class="block flex font-normal hover:bg-[#393636] jus mt-1 px-2 py-2 space-x-2 text-sm w-full"
          >
            <img
              v-if="notification.imageSrc"
              :src="notification.imageSrc"
              alt="ALET_IMG"
              class="h-12 rounded w-12"
            >

            <div class="flex flex-col w-full">
              <div class="flex justify-between">
                <p
                  class="break-all font-bold"
                >
                  {{ notification.messages.find(m => m.langCode === selectedLang.toUpperCase())?.title ?? "" }}
                </p>
                <p class="italic">
                  {{ new Date(notification?.createdAt).toLocaleDateString() }}
                </p>
              </div>
              <div
                class="notification"
                v-html="notification.messages.find(m => m.langCode === selectedLang.toUpperCase())?.text"
              />
              <div
                v-if="('messages' in notification) && selectedDashboard.channelId === user?.id && showNew"
                class="text-right"
              >
                <button
                  type="button"
                  class="duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-[#772CE8] inline-block leading-tight px-1 py-1 rounded shadow text-white transition uppercase"
                  @click="markNotificationAsReaded(notification)"
                >
                  <Mark style="width: 16px; height: 16px;" />
                </button>
              </div>
            </div>
          </div>
        </PopoverPanel>
      </div>
    </Popover>
  </div>
</template>

<style scoped>
.notification :deep(a) {
  @apply text-purple-500 font-bold
}
</style>