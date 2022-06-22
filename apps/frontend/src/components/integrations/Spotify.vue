<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { ChannelIntegration } from '@tsuwari/prisma';
import { computed, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import { api } from '@/plugins/api';
import { setSpotifyProfile, spotifyProfileStore } from '@/stores/spotifyProfile';
import { selectedDashboardStore } from '@/stores/userStore';

const router = useRouter();
const spotifyIntegration = ref<Partial<ChannelIntegration>>({
  enabled: true,
});
const selectedDashboard = useStore(selectedDashboardStore);
const spotifyProfile = useStore(spotifyProfileStore);
const authurl = computed(() => {
  return `${window.location.origin}/api/v1/channels/${selectedDashboard.value.channelId}/integrations/spotify/auth`;
});
const { t } = useI18n({
  useScope: 'global',
});

selectedDashboardStore.subscribe(d => {
  api(`/v1/channels/${d.channelId}/integrations/spotify`).then(async (r) => {
    spotifyIntegration.value = r.data;
    fetchSpotifyProfile();
  });
});


async function fetchSpotifyProfile() {
  const { data } = await api(`v1/channels/${selectedDashboard.value.channelId}/integrations/spotify/profile`);
  setSpotifyProfile(data);
}

async function patch() {
  const { data } = await api.patch(`v1/channels/${selectedDashboard.value.channelId}/integrations/spotify`, {
    enabled: spotifyIntegration.value.enabled,
  });

  spotifyIntegration.value = data;
}

onMounted(async () => {
  const route = router.currentRoute.value;
  const params = new URLSearchParams(window.location.search);
  const code = params.get('code');

  if (route.params.integration === 'spotify' && code) {
    if (code) {
      await api.post(`v1/channels/${selectedDashboard.value.channelId}/integrations/spotify/token`, {
        code,
      });
      await fetchSpotifyProfile();

      return router.push('/dashboard/integrations');
    }
  }

  if (!spotifyProfile.value) {
    fetchSpotifyProfile();
  }
});

</script>

<template>
  <div class="flex flex-col card rounded card-compact bg-base-200 drop-shadow p-4">
    <div class="flex justify-between mb-5">
      <div>
        <h2 class="card-title flex font-bold">
          Spotify
        </h2>
      </div>
      <div class="form-check form-switch">
        <input
          id="flexSwitchCheckDefault"
          v-model="spotifyIntegration.enabled"
          class="form-check-input appearance-none w-9 -ml-10 rounded-full float-left h-5 align-top bg-no-repeat bg-contain bg-gray-300 focus:outline-none cursor-pointer shadow"
          type="checkbox"
          role="switch"
          @change="patch"
        >
      </div>
    </div>

    <div class="mb-5">
      <div
        v-if="spotifyProfile"
      >
        <div class="flex justify-center mb-3">
          <img
            v-if="spotifyProfile.images"
            :src="spotifyProfile.images[0].url"
            class="rounded-full w-32 ring-2 ring-white"
            alt="Avatar"
          >
        </div>
        <p class="text-center break-words">
          {{ spotifyProfile.display_name }}#{{ spotifyProfile.id }}
        </p>
      </div>
      <div v-else>
        Not logged in
      </div>
    </div>

    <div class="mt-auto text-right">
      <a
        class="px-6 py-2.5 inline-block bg-purple-600 text-white font-medium text-xs leading-tight uppercase rounded shadow hover:bg-purple-700 hover:shadow focus:bg-purple-700 focus:shadow focus:outline-none focus:ring-0 active:bg-purple-800 active:shadow transition duration-150 ease-in-out"
        :href="authurl"
      >
        {{ t('buttons.login') }}
      </a>
    </div>
  </div>
</template>

<style scoped>
.tooltip { 
width: 140px;
background: #59c7f9;
color: #ffffff;
text-align: center;
padding: 10px 20px 10px 20px;
border-radius: 10px;
top: calc(100% + 11px);
left: 50%;
transform: translate-x(-50%)
 }
.tooltip-box { 
position: relative
 }
.triangle { 
border-width: 0 6px 6px;
border-color: transparent;
border-bottom-color: #59c7f9;
position: absolute;
top: -6px;
left: 50%;
transform: translate-x(-50%)
 }
</style>